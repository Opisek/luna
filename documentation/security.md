# Security & Privacy
Luna has been designed with a strong focus on security.

This document is divided into four parts:
1. The [Security Measures](#1-security-measures) chapter describes all the security mechanisms built into Luna.
2. The [System Compromise Effect Analysis](#2-system-compromise-effect-analysis) chapter describes to what extent a compromise of different parts of the system would compromise the user's security.
3. The [User Considerations](#3-user-considerations) chapter describes different limitations of Luna, how they might affect users (including the administrator), and suggestions based on one's individual security requirements.
4. The [Possible or Planned Additional Security Measures](#4-possible-or-planned-additional-security-measures) chapter goes over additional security measures that could be implemented in the future.

## 1 Security Measures
### 1.1 Internal IDs
### 1.2 Error Verbosity
### 1.3 Avoidance of Timing Side-Channels
### 1.4 Password Hashing
### 1.5 Authorization Token Design
### 1.6 Database Encryption
### 1.7 Web Attacks Protections
#### 1.7.1 CSRF Protections
#### 1.7.2 CORS Protections
#### 1.7.3 XSS Protections
#### 1.7.4 SQL Injection Protections
#### 1.7.5 Token Leak Protections

## 2 System Compromise Effect Analysis
### 2.1 Database Leak
### 2.2 Cryptographic Key Store Leak
### 2.3 Combined Database and Cryptographic Key Store Leak
### 2.3 Database Compromise
### 2.4 Frontend Compromise
### 2.5 Backend Compromise
### 2.5 Proxy Compromise

## 3 User Considerations
### 3.1 Zero-Day Vulnerabilities

### 3.2 Privacy Concerns
#### 3.2.1 Gravatar
#### 3.2.2 Remote Resource IP Leak
#### 3.2.3 Geolocation Privacy Concerns
#### 3.2.4 IP Address and User Agent Logging Privacy Concerns

### 3.3 Hosting of Illicit Content

## 4 Possible or Planned Additional Security Measures
### 4.1 Option to Disable iCal Uploading
### 4.2 Option to Disable iCal Caching
### 4.3 Decentralized User Database Encryption Key Storage
- Salted password hash used as seed for KDF.
- Resulting key is used to encrypt that user's sensitive database entries.
- The key is encrypted using the server's master key together with the user's ID.
- The key is never stored on the server, instead it is placed in an HTTP-only cookie.

Still need to think about calendar shares.
They will probably have to be re-encrypted using a different key.
### 4.4 Further Database Encryption
### 4.5 Multi-Factor Authentication and Passkeys
### 4.6 External Authentication