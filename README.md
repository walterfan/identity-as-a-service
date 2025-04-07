# Identity as a Service

Identity as a Service

## FAQ

### What's IAM

IAM is Identity and Access Management, which is a tool to manage fine-grained authorization. 
In other words, it lets you control who can do what on which resources.

Giving someone permissions in IAM involves the following three components:

- Principal: The identity of the person or system that you want to give permissions to
- Role: The collection of permissions that you want to give the principal
- Resource: The Google Cloud resource that you want to let the principal access
  
To give the principal permission to access the resource, you grant them the role on the resource. You grant these roles using an allow policy.

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
