# Security & Privacy
Luna has been designed with a strong focus on security.

This document is divided into four parts:
1. The [Attacker Model](#1-attacker-model) goes over the assumptions made about adversaries that the system's controls defend against.
2. The [Security Measures](#2-security-measures) chapter describes all the security mechanisms built into Luna.
3. The [System Compromise Effect Analysis](#3-system-compromise-effect-analysis) chapter describes to what extent a compromise of different parts of the system would compromise the user's security.
4. The [User Considerations](#4-user-considerations) chapter describes different limitations of Luna, how they might affect users (including the administrator), and suggestions based on one's individual security requirements.
5. The [Possible or Planned Additional Security Measures](#5-possible-or-planned-additional-security-measures) chapter goes over additional security measures that could be implemented in the future.

## 1 Attacker Model

## 2 Security Measures
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
### 2.8 Web Attacks Protections
#### 2.8.1 CSRF Protections
#### 2.8.2 CORS Protections
#### 2.8.3 XSS Protections
#### 2.8.4 SQL Injection Protections
#### 2.8.5 Token Leak Protections

## 3 System Compromise Effect Analysis
### 3.1 Database Leak
### 3.2 Cryptographic Key Store Leak
### 3.3 Combined Database and Cryptographic Key Store Leak
### 3.3 Database Compromise
### 3.4 Frontend Compromise
### 3.5 Backend Compromise
### 3.5 Proxy Compromise

## 4 User Considerations
### 4.1 Zero-Day Vulnerabilities

### 4.2 Privacy Concerns
#### 4.2.1 Gravatar
#### 4.2.2 Remote Resource IP Leak
#### 4.2.3 Geolocation Privacy Concerns
#### 4.2.4 IP Address and User Agent Logging Privacy Concerns

### 4.3 Hosting of Illicit Content

## 5 Possible or Planned Additional Security Measures
### 5.1 Option to Disable iCal Uploading
### 5.2 Option to Disable iCal Caching
### 5.3 Decentralized User Database Encryption Key Storage
- Salted password hash used as seed for KDF.
- Resulting key is used to encrypt that user's sensitive database entries.
- The key is encrypted using the server's master key together with the user's ID.
- The key is never stored on the server, instead it is placed in an HTTP-only cookie.

Still need to think about calendar shares.
They will probably have to be re-encrypted using a different key.
### 5.4 Further Database Encryption
### 5.5 Multi-Factor Authentication and Passkeys
### 5.6 External Authentication
### 5.7 Improvements to Avoidance of Timing Side-Channels
In order to further reduce information deductible from the time taken to complete authentication-related requests, an artificial fixed delay of a few seconds might be added to all such requests. Note that such a measure would inadvertently slow down the processing time of these request for legitimate users, as well.