# Identity as a Service

Identity as a Service

example: https://github.com/keycloak/keycloak

## IAM (Identity and Access Management)

### **What is IAM?**  
IAM (Identity and Access Management) is a security framework that ensures the right individuals (users, services, or systems) have the appropriate access to resources in a cloud or on-premises environment. It manages:  
- **Identities** (users, groups, roles, service accounts)  
- **Permissions** (policies defining what actions they can perform)  
- **Authentication** (verifying identity) and **Authorization** (granting access)  

### **Why Do We Need IAM?**  
1. **Security** – Prevents unauthorized access to sensitive data and systems.  
2. **Least Privilege Principle** – Grants only necessary permissions to users/applications.  
3. **Compliance** – Helps meet regulatory requirements (e.g., GDPR, HIPAA).  
4. **Centralized Control** – Manages access across multiple services from one place.  
5. **Auditability** – Tracks who did what (via logs and monitoring).  

### **How Does IAM Work?**  
1. **Authentication** – Verifies identity (e.g., via passwords, MFA, or federated login).  
2. **Authorization** – Checks policies to determine allowed actions (e.g., "Can User X delete S3 buckets?").  
3. **Policies & Roles** – Defines permissions (e.g., AWS IAM policies, Azure RBAC roles).  
4. **Access Management** – Assigns roles to users/groups (e.g., "Admin," "Read-Only").  
5. **Audit & Monitoring** – Logs access events for security analysis.  

### **Example (AWS IAM)**  
- A **policy** allows an **IAM user** to "read S3 buckets."  
- A **role** grants an **EC2 instance** permission to access DynamoDB.  
- **MFA** adds an extra layer of security for admin logins.  

## Protocols

| **Protocol** | **Role in IAM**                       | **Commonly Paired With**      |  
|--------------|---------------------------------------|-------------------------------|  
| OAuth 2.0    | API authorization                     | OIDC, JWT                     |  
| OIDC         | Authentication (SSO)                  | OAuth 2.0, SAML               |  
| SAML         | Enterprise SSO                        | LDAP, AD                      |  
| JWT          | Token format for claims               | OIDC, OAuth                   |  
| LDAP         | User directory storage                | Kerberos, SAML                |  
| SCIM         | User lifecycle management             | OAuth (for API calls)         |  
| FIDO2        | Passwordless auth                     | OIDC (as second factor)       |  

### **OAuth 2.0**  
- **Purpose**: Delegated authorization (allows apps to access resources on behalf of users).  
- **Key Use Case**: "Login with Google/Facebook," API access tokens.  
- **How It Works**:  
  - Issues short-lived **access tokens** (e.g., for APIs).  
  - Uses **grant types** (Authorization Code, Client Credentials, etc.).  
- **Not for Authentication** (but often misused as such).  

### **OpenID Connect (OIDC)**  
- **Purpose**: Authentication layer on top of OAuth 2.0.  
- **Key Use Case**: Single Sign-On (SSO), user identity verification.  
- **How It Works**:  
  - Issues **ID tokens** (JWT format) containing user claims (e.g., email, name).  
  - Uses OAuth flows + standardizes identity data (via `/userinfo` endpoint).  

### **SAML 2.0 (Security Assertion Markup Language)**  
- **Purpose**: Enterprise SSO (XML-based).  
- **Key Use Case**: Corporate logins (e.g., Active Directory integrations).  
- **How It Works**:  
  - Exchanges **SAML assertions** between Identity Provider (IdP) and Service Provider (SP).  
  - Browser redirects with base64-encoded XML messages.  

### **JWT (JSON Web Token)**  
- **Purpose**: Compact, self-contained token format for claims.  
- **Key Use Case**: Stateless authentication (used in OIDC, API auth).  
- **Structure**:  
  ```text
  Header.Claims.Signature
  ```  
  - Contains **claims** (e.g., `exp`, `sub`) and is signed/encrypted (JWS/JWE).  

### **LDAP (Lightweight Directory Access Protocol)**  
- **Purpose**: Directory services (user/group storage).  
- **Key Use Case**: On-premises IAM (e.g., Microsoft AD, OpenLDAP).  
- **How It Works**:  
  - Queries hierarchical databases (users, groups, attributes).  
  - Often paired with **Kerberos** for authentication.  

### **SCIM (System for Cross-domain Identity Management)**  
- **Purpose**: Automate user provisioning/deprovisioning.  
- **Key Use Case**: Sync users between IdP and apps (e.g., Okta → Slack).  
- **How It Works**: REST API standard for CRUD operations on user identities.  

### **FIDO2/WebAuthn**  
- **Purpose**: Passwordless authentication.  
- **Key Use Case**: Biometrics/security keys (e.g., YubiKey, Windows Hello).  
- **How It Works**: Public-key cryptography + browser APIs.  

## FAQ

### What's Workload, Workload Identity and Workload Attestation?
* **Workload**: A discrete computing process, service, or application that performs specific tasks within a computing environment. Examples include containers, microservices, virtual machines, or serverless functions.
* **Workload Identity**: An identity assigned to a workload rather than a human user, allowing the workload to authenticate and interact with other services securely without embedded credentials.
* **Workload Attestation**: The process of verifying a workload's identity and integrity by evaluating various attributes such as where it's running, how it was deployed, and its runtime characteristics before granting access to resources.

### What's SPIFFE?
* **SPIFFE (Secure Production Identity Framework For Everyone)**: An open standard that defines how services can identify and authenticate to each other securely.
* It provides a framework for establishing trust between workloads across heterogeneous environments without hardcoded credentials.
* SPIFFE defines SVID (SPIFFE Verifiable Identity Document) as the document format for identifying workloads.
* The SPIRE (SPIFFE Runtime Environment) project implements the SPIFFE specifications for practical use.

### What's Federation?
* **Identity Federation**: The process of linking and trusting identity information across different identity management systems.
* In the context of workload identity, federation allows workloads to authenticate across different trust domains, cloud providers, or organizational boundaries.
* Federation enables consistent identity verification without requiring separate credentials for each environment.

### What's JWT SVID?
* **JWT SVID (JSON Web Token SPIFFE Verifiable Identity Document)**: A type of SPIFFE identity document formatted as a signed JWT.
* It contains claims about a workload's identity and is signed by a trusted authority.
* JWT SVIDs are portable across different systems and can be verified without contacting the issuing authority.
* They typically include information about the workload's identity, permissions, and expiration time.

### What's AWS Roles Anywhere?
* A service that extends AWS IAM roles to workloads running outside of AWS (on-premises, other clouds, edge locations).
* It allows non-AWS workloads to obtain temporary AWS credentials by authenticating with X.509 certificates.
* Workloads can use these credentials to access AWS services with the same permissions model used within AWS.
* It eliminates the need for long-term credentials while maintaining the principle of least privilege.

### What's AWS OIDC Federation?
* **AWS OIDC Federation**: Integration between AWS IAM and OpenID Connect (OIDC) identity providers.
* It allows workloads that can authenticate with an OIDC provider to assume AWS IAM roles.
* Commonly used for Kubernetes workloads, CI/CD pipelines, and other cloud services to access AWS resources.
* The federation process uses JWT tokens from the OIDC provider to establish trust with AWS.


## Demo
* backend: go
* frontend: vue
```
idaas/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   ├── middleware/
│   │   │   └── routes.go
│   │   ├── config/
│   │   ├── models/
│   │   ├── repository/
│   │   └── service/
│   ├── pkg/
│   │   ├── auth/
│   │   ├── logger/
│   │   └── validator/
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── assets/
│   │   ├── components/
│   │   ├── composables/
│   │   ├── layouts/
│   │   ├── router/
│   │   ├── store/
│   │   ├── views/
│   │   ├── App.vue
│   │   └── main.js
│   ├── .env
│   ├── package.json
│   └── tailwind.config.js
└── README.md
```

## Create DB and environment file
* create datbase

```sql

CREATE DATABASE IF NOT EXISTS idaas;
CREATE DATABASE IF NOT EXISTS keycloak;

CREATE USER IF NOT EXISTS 'your_username'@'localhost' IDENTIFIED BY 'your_password';

GRANT ALL PRIVILEGES ON keycloak.* TO 'your_username'@'localhost';
GRANT ALL PRIVILEGES ON idaas.* TO 'your_username'@'localhost';

FLUSH PRIVILEGES;
```

* create a file - .env

```
KC_DB=mysql
KC_DB_URL=jdbc:mysql://keycloak-db:3306/keycloak
KC_DB_USERNAME=keycloak
KC_DB_PASSWORD=password
KEYCLOAK_ADMIN=admin
KEYCLOAK_ADMIN_PASSWORD=password

```

* create self-signed certificate
```
mkdir -p nginx/certs
openssl req -x509 -nodes -days 1095 -newkey rsa:2048 \
  -keyout nginx/certs/server.key \
  -out nginx/certs/server.crt \
  -subj "/CN=localhost"
```