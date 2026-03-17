# Security & Privacy
Luna has been designed with a strong focus on security.

This document is divided into five parts:
1. The [Attacker Model](#1-attacker-model) goes over the assumptions made about adversaries that the system's controls defend against.
2. The [Security Measures](#2-security-measures) chapter describes the security mechanisms built into Luna.
3. The [System Compromise Effect Analysis](#3-system-compromise-effect-analysis) chapter describes to what extent a compromise of different parts of the system would compromise the user's security.
4. The [User Considerations](#4-user-considerations) chapter describes different limitations of Luna, how they might affect users (including the administrator), and suggestions based on one's individual security requirements.
5. The [Possible or Planned Additional Security Measures](#5-possible-or-planned-additional-security-measures) chapter goes over additional security measures that could be implemented in the future.

## 1 Attacker Model
The attacker is assumed to:
- Have perfect knowledge of the protocols used
- Be able to initiate communication with the server and send their own requests
- Forge `GET` and `HEAD` requests by means of Cross-Site-Request-Forgery

The attacker is assumed to not:
- Be able to break cryptographic primitives
- Be able to decrypt TLS traffic
- Sniff or manipulate traffic between the reverse proxy, the frontend, the backend, and the database

## 2 Security Measures
This chapter describes all security mechanisms built into Luna.

### 2.1 Internal IDs
Luna uses its own IDs for every resource accessed through it, like calendars or events. These IDs are based on UUIDv4 and UUIDv5.

This has two effects on privacy of users:
- To an outsider, the IDs are unpredictable by nature. In particular, an attacker can neither guess IDs of resources, nor can they deduct any information from sniffed IDs. For example, the amount of registered users cannot be deduced from IDs of users, since the IDs are not consecutive.
- The nature of the upstream sources are hidden. This means that an attacker is not able to tell which kind of calendar backend a user is using based solely on a sniffed calendar or event ID.

### 2.2 Error Verbosity
The error messages used by Luna can be serialized at four different verbosity levels, which can be selected in the settings by an administrator. These are:

1. **Broad** (minimal) level
    - Very generic and minimal level of detail
    - For users with very limited technical knowledge
    - MAY be used in high security scenarios
    - Error SHOULD only reflect the just attempted action
    - e.g. User tries to get calendars => Could not get calendars

2. **Plain** (normal) level
    - Balanced details for standard users
    - MAY be used with infrastructure terms (Frontend, Backend, Database, File, ...)
    - MUST NOT be used with technical terms
    - e.g. User tries to get calendars => Could not get calendars from source <name>

3. **Wordy** (verbose) level
    - Detailed logs for advanced users
    - Does not pose a security threat in production
    - MUST use when referring to technical terms (CalDAV, iCal, HTTP Request, ...)
    - MUST NOT contain internal IDs or errors
    - e.g. User tries to get calendars => Could not get calendars from iCal file

4. **Debug** Level
   - Full stack trace of all errors
   - MIGHT reveal sensitive information
   - MUST NOT be used in production
   - MUST be used when returning any internal IDs or errors
   - e.g. User tries to get calendars => Could not get calendars from file <ID> belonging to source <ID>

By default, the **plain** logging level is used. This way, no sensitive internal information is released to potential malicious actors. Furthermore, all endpoints related to user authentication return very generic error messages on all levels but **debug** as to not create any oracles.

Note that errors outputted to `stderr` of the backend use the **debug** level.

### 2.3 Avoidance of Timing Side-Channels
All endpoints related to user authentication (i.e., endpoints requiring the user's password) perform the same [password hashing](#24-password-hashing) procedure, even if the authentication fails at an earlier stage. This is to try and avoid timing side-channels related to user information and passwords.

For example, trying a username-password combination should take similar time, whether a user with the supplied username actually exists in the system or not.

Note that exactly identical timings are not guaranteed. This is already unfeasible due to the time that the database queries might take to run. For potential improvements in the future, see [Improvements to Avoidance of Timing Side Channels](#57-improvements-to-avoidance-of-timing-side-channels).

### 2.4 Password Hashing
User passwords are stored in a state-of-the-art hashed form with the following properties:
1. The hashing algorithm used is Argon2 with the following parameters:
   - Time: 1
   - Memory: 65536
   - Threads: 4
2. The hash uses a 32 byte long salt value generated with a CSPRNG.
3. The hash uses a 64 byte long pepper value generated with a CSPRNG.

This has the following effects:
1. Hashing the password with a cryptographic password hash function ensures that:
    - To our knowledge, there exists no more efficient way to obtain the plaintext password from the hash value than brute-forcing (for reasoning about rainbow tables, see the following two points).
    - The calculation of the hash is artificially slowed down to significantly lower the speed at which the passwords could be brute-forced, even in an offline attack, such that it becomes probabilistically infeasible for modern computers to ever "crack" a password of enough entropy (very short passwords or passwords particularly susceptible to dictionary attacks remain vulnerable).
2. Salting the hash prevents the use of rainbow tables to reduce the amount of computation needed to brute-force a password through the use pre-computed values.
3. Peppering the hash ensures that an attacker is unable to brute-force passwords given just a database leak without the pepper stored in the backend's data.

### 2.5 Rate Limiting
To further prevent online brute-force attacks, rate limiting is implemented, which dynamically adjusts the amount of requests that an attacker is able to issue in a given time-span.

This is done by recording the amount of HTTP request in the backend per client IP address that were answered with a failure status code. The counter is reset if no request fails for at least 5 minutes.

The throttle values function scale as follows:
- More than 5 failed requests: 100ms of artificial delay for every subsequent request
- More than 10 failed requests: 1s of artificial delay for every subsequent request
- More than 15 failed requests: 5s of artificial delay for every subsequent request
- More than 50 failed requests: All subsequent requests are refused to be processed

As such, the brute-force capability of an attacker is effectively limited to 50 attempts per 5 minutes.

### 2.6 Authorization Token Design
When a user logs in successfully, the following data are stored in the database:
- Session ID
- User ID
- Timestamp
- IP address
- User agent
- Secret hash

Furthermore, a JSON Web Token (JWT) is issued with the following information:
- Session ID
- User ID
- Secret

The integrity and authenticity of the JWT is guaranteed by a SHA-512 HMAC generated using a key stored in the backend.

The function of each of the saved values is as follows:
- **Session ID**: This allows the user to revoke (sign out of) specific sessions. This means that compromised devices or JWTs cannot be used to access the user's account indefinitely.
- **User ID**: This is required to correctly perform authorization based on the determined user identity.
- **Timestamp / IP address**: This allows the user to easier spot unknown sessions in case of an account compromise.
- **User agent**: This is a hurdle preventing attackers from using tokens stolen through social engineering without additional data. If the user agent of a request does not match the user agent saved in the database (OS and browser updates are possible), then the session is immediately revoked.
- **Secret / Secret hash**: Upon login, a 256 byte long secret is generated and stored in the token. A corresponding peppered SHA-256 hash is stored in the database. For a token to be accepted, the secret passed in the token must hash to the value stored in the database. This means that an attacker who somehow obtained access to the keys used for generating the JWT's HMAC is unable to forge JWTs that would be accepted, even if they gained a read access of the database, since they would need to know the secret value stored inside an existing token, or they would have to insert a new row into the database.

### 2.7 Database Encryption
In order to access users' calendars, Luna must save password and tokens for the upstreams inside the database. To help protect these secrets in case of a database leak, this information is stored in an encrypted state using keys stored inside the backend's data directory. The encryption algorithm used is PostegreSQL's built-in OpenPGP implementation. The same goes for OAuth 2.0 client secrets.

Note that, if an attacker gets a hold of the database encryption keys (this could also be a malicious administrator), then they are able to decrypt these entries. An updated version that addresses this is in the works. See [Decentralized User Database Encryption Key Storage](#53-decentralized-user-database-encryption-key-storage)

### 2.8 Web Attacks Protections
In the following, protections against common web attacks are described.

#### 2.8.1 CSRF Protections
Two mitigations against Cross-Site-Resource-Forgery are in place:
1. The authorization token issued to the client is saved with the `SameSite` attribute set to `Lax`. This means that the cookie is never sent on `POST`, `PATCH`, and `DELETE` requests originating from a different domain (cross-site). Special care is taken to limit the amount of side effects that `GET` and `HEAD` requests have in the backend. Things like the most recently used IP address are still updated in the database. The OAuth 2.0 authorization flow also relies on a `GET` request to receive an authorization code from the authorization server.
2. On all `POST`, `PATCH`, and `DELETE` requests, the equality of the `Origin` header and the `PUBLIC_URL` environmental variable is checked. If this test fails, the request is rejected.

#### 2.8.2 CORS Protections
The frontend implements a Cross-Origin Resource Sharing policy which rejects all API responses that do not carry the frontend's `PUBLIC_URL` value inside the `Access-Control-Allow-Origin` header.

#### 2.8.3 XSS Protections
Preventing Cross-Site-Scripting relies on how Svelte (the frontend library) sanitizes input before embedding it inside the DOM. Special care is taken to never inject user input directly into the rendered site.

#### 2.8.4 SQL Injection Protections
Prepared statements are used every time that values chosen by the user need to be placed inside an SQL query. User input is never simply concatenated into a raw query string.

#### 2.8.5 Token Leak Protections
An adversary may try to trick a user into giving them their authorization token via means of social engineering. To prevent this from happening by pasting unknown scripts into the browser's console, the token's cookie is saved with the `HttpOnly` attribute set.

If the attacker succeeds in tricking the user into extracting the token directly from the browser's cookie store and sending it over, the user agent protections described in [Authorization Token Design](#26-authorization-token-design) pose an additional hurdle, as the token cannot be used without extracting additional data from the user.

## 3 System Compromise Effect Analysis
In the case that parts of the system are compromised, the attacker model described in [Attacker Model](#1-attacker-model) breaks. In this chapter the extent that a compromise of the different parts of the system might have, i.e., what an attacker with control over these parts is able or is not able to do.

### 3.1 Database Leak
Giving a leak of the state of the database at a singular moment, the attacker gains infromation about:
- Usernames
- E-Mail Addresses
- Names picked for user's sources
- Calendar and event overrides
- Password hashes
- Encrypted upstream credentials

The password hashing described in [Password Hashing](#24-password-hashing) ensures that the attacker could not feasible extract the user's plaintext password from the leaked hash. This is because:
- The salt used in the password prevents the use of rainbow tables and correlation attacks.
- The password hash function used in irreversible and significantly slows down brute-force attacks.
- The pepper that is saved outside the database adds additional 64 bytes of entropy that the attacker would have to guess with no other way of verifying the correctness but to perform an online brute-force attack, running into [Rate Limiting](#25-rate-limiting).

The credentials encryption described in [Database Encryption](#27-database-encryption) means that the attacker does not get access to the user's upstream calendar sources.

If this compromise scenario occurs, it is recommended to reset the backend's cryptographic keys. Note that this will require you to purge the database, too.

### 3.2 Cryptographic Key Store Leak
If the cryptographic keys used by the backend are leaked, the attacker does not gain any additional information or capabilities without a simultaneous database leak.

If this compromise scenario occurs, it is recommended to reset the backend's cryptographic keys. Note that this will require you to purge the database, too.

### 3.3 Combined Database and Cryptographic Key Store Leak
If both the keys used by the backend and a database snapshot are leaked, then the attacker can extract the user's credentials for their upstream calendar sources.

Furthermore, the attacker obtains the pepper used for security passwords, meaning that the artificially added entropy is defeated. In such a scenario, the attacker could brute-force very short passwords or otherwise vulnerable passwords. Make sure that your users use secure and unique passwords.

If this compromise scenario occurs, it is recommended to reset the backend's cryptographic keys and purge the database. Advise all your users to change their passwords for their upstream calendar sources.

### 3.3 Database Compromise
If the attacker gains read-write SQL access to the database, then they get access to the same information as in the case of a [Database Leak](#31-database-leak). They are not able to log in as a legitimate user, because they lack the server secret used for creating token hashes to be saved in the database.

If the attacker compromises the machine that hosts the database, then they could extract the encryption keys used to secure database entries directly from the system's memory. This would have an effect comparable to [Combined Database and Cryptographic Key Store Leak](#33-combined-database-and-cryptographic-key-store-leak).

If you suspect that your database has been compromised, kill the backend immediately. This must be done to prevent the possibility of the attacker extracting cryptographic keys from the RAM upon a database query. Additionally, the cryptographic keys and the database should both be purged. Advise all your users to change their passwords for their upstream calendar sources.

### 3.4 Frontend Compromise
If the frontend is compromised, the attacker taken on a passive Man in the Middle position. This means that they can see everything that the user sends to and receives from the backend. This includes user passwords when the user logs in. Additionally, the attacker can elevate their position to an active adversary if a user sends their authorization token to the frontend. In such a case, the adversary can act as the user and access all their resources.

If you suspect a frontend compromise, it is imperative that no user accesses Luna, so that the attacker does not gain access to user tokens or passwords. Kill the frontend or the reverse proxy and kill the backend. After a safe environment is restored, make sure to drop all entries from the `sessions` database table and advise all your users to change their passwords—both for Luna and their upstream calendar sources.

### 3.5 Backend Compromise
If the backend is compromised, the attacker gains the same capabilities as given both a [Frontend Compromise](#34-frontend-compromise) and a [Combined Database and Cryptographic Key Store Leak](#33-combined-database-and-cryptographic-key-store-leak).

If you suspect a compromise, immediately kill all the services. Purge the database and the cryptographic keys, and advise all your users to change their passwords for their upstream calendar sources.

### 3.5 Proxy Compromise
The compromise of a reverse proxy poses the same danger as a [Frontend Compromise](#34-frontend-compromise).

If you suspect a proxy compromise, it is imperative that no user accesses Luna, so that the attacker does not gain access to user tokens or passwords. Kill the frontend or the reverse proxy and kill the backend. After a safe environment is restored, make sure to drop all entries from the `sessions` database table and advise all your users to change their passwords—both for Luna and their upstream calendar sources.

## 4 User Considerations
The following chapter describes different limitations of Luna, how they might affect users (including the administrator), and suggestions based on one's individual security requirements.

### 4.1 Zero-Day Vulnerabilities
Zero-day vulnerabilities in the used libraries and cryptographic ciphers are always a possibility. The attack surface is attempted to be minimized by using as few external libraries as feasible. The best way to protect against these eventualities is to keep all your software up-to-date and stay informed about the latest cybersecurity news.

### 4.2 Privacy Concerns
In the following we go over a few features that may hurt the privacy of your users when enabled, as well as what protections exist to limit that.

#### 4.2.1 Gravatar
Luna implements Gravatar for user avatars. Since this is an external service, your users reveal that they use your Luna instance to Gravatar. If this goes against your security model, you can disable the use of Gravatar in the administrative settings.

#### 4.2.2 Remote Resource IP Address Leak
Users are able to set their profile pictures to remote pictures, just as they are able to import remote `.ics` files. To make sure that your user's IP addresses are not revealed to the operators hosting these remote resources, all such content is first cached by Luna and then passed onto the user. This way, only the IP address of your Luna instance is revealed.

To also prevent the upstream from finding out your usage patterns, profile pictures are recached periodically, as opposed to on an explicit user request. The same is not possible for those `.ics` files that require user authorization, due to the credentials being encrypted as described in [Database Encryption](#27-database-encryption).

#### 4.2.3 Geolocation Privacy Concerns
To help your users notice unknown log-ins, the IP addresses used by the user are mapped to their estimated geographical location when the user views their active sessions in Luna's settings. This is done with the external service ipapi. Note that this geographical location is computed on the user's machine and never stored anywhere.

Using this feature means that all the IP addresses of a user are sent to the external third party for purposes of geolocation. If this goes against your security model, you can disable this feature in the administrative settings.

#### 4.2.4 IP Address and User Agent Logging Privacy Concerns
To help your users notice unknown log-ins, the IP addresses and user agents of the users are stored in the database. Only the first IP address that the user uses to log into a particular session and the most recently used IP address associated with each session are stored.

### 4.3 Hosting of Illicit Content
There is always a concern that a user uploads illegal content to your Luna instance, either as a calendar event or as a profile picture. Make sure you only invite people you trust to use your instance. You can also disable the ability to upload your own profile pictures in the administrative settings.

## 5 Possible or Planned Additional Security Measures
This chapter goes over additional security measures that could be implemented in the future.

### 5.1 Option to Disable iCal Uploading
To help reduce the risk of [Hosting of Illicit Content](#43-hosting-of-illicit-content), an option to disable uploading of your own `.ics` files might be added. 

### 5.2 Option to Disable Remote Resource Caching
To help reduce the risk of [Hosting of Illicit Content](#43-hosting-of-illicit-content), an option to disable caching of remote resources might be added. Note that opting for this would break the protection described in [Remote Resource IP Address Leak](#422-remote-resource-ip-address-leak).

### 5.3 Decentralized User Database Encryption Key Storage
Currently, access to the backend's cryptographic keys allows an entity to potentially decrypt users' credentials for their upstream calendar sources.

An improvement to [Database Encryption](#27-database-encryption) is planned to eliminate such a risk. In the new system, only the owner of the upstream source would be able to order their credentials to be decrypted upon their request. This would be done by deriving the encryption keys for these credentials in the following way:
- When the user logs in, their password is hashed independently a second time with a different salt.
- The resulting hash is used as a seed for a KDF.
- The KDF is used to derive a user encryption key.
- The resulting key is used to encrypt that user's sensitive database entries.
- To allow for calendars to be shared in the future, the database rows are encrypted transitively using that key, i.e., a fresh key is generated for every calendar and that key is stored in the database encrypted with the user's encryption key.
- The user encryption key is encrypted using the server's master key together with the user's ID.
- The encrypted key is never stored on the server, instead it is placed in an HTTP-only cookie. This way, the server itself is unable to decrypt the user's database entries. Only when the user makes an explicit request and sends their cookies, is it possible to access these credentials.

### 5.4 Further Database Encryption
In the future, more database columns might be encrypted.

### 5.5 Multi-Factor Authentication and Passkeys
Multi-factor authentication and passkeys will be added in the future to improve account security.

### 5.6 External Authentication
Support for external authentication services, for example based on LDAP, might be added in the future.

### 5.7 Improvements to Avoidance of Timing Side-Channels
In order to further reduce information deductible from the time taken to complete authentication-related requests, an artificial fixed delay of a few seconds might be added to all such requests. Note that such a measure would inadvertently slow down the processing time of these request for legitimate users, as well.