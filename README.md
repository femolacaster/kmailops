

## KMailOps 📧

### A Kubernetes Email Operator with Multiple Providers

This project builds a Kubernetes operator written in Go that manages custom resources (CRDs) for configuring and sending emails through transactional email providers like MailerSend 📨 and/or Mailgun 🔫.

### Project Goals

- Simplify email sending in Kubernetes: Automate email configuration and delivery for applications running in Kubernetes clusters 🤖
- Support multiple providers: Integrate with both MailerSend and Mailgun APIs for sending emails 🤝
- CRD-based configuration: Define email sending configurations and messages using CRDs for a declarative approach
- Cross-namespace operation: Allow the operator to manage email sending across different namespaces within the cluster

### Deliverables

- Source Code: The complete Go application code for the operator will be available on this public Github repository 🐈
- Deployment Manifests: YAML manifests for deploying the operator as a pod within the Kubernetes cluster ⚙️
- Documentation: Documentation will cover installation, configuration, and usage of the operator 📚

### Technical Specifications

_Custom Resource Definitions (CRDs)_

- `EmailSenderConfig`: This CRD defines the configuration details for an email provider, including the API token (stored in a Kubernetes secret) and sender email address 🔑
- `Email`: This CRD defines an email message to be sent, referencing the sender configuration and specifying the recipient email, subject, and body content. The status field of the Email CRD will reflect the delivery status (success/failure), message ID (if successful), and any error messages ⚠️

### Operator Implementation

The operator is written in Go and leverages the Operator SDK or Kubebuilder framework for ease of development 🛠️. It watches for changes to the Email and EmailSenderConfig CRDs using Kubernetes controllers 👀. Upon creation or update of an EmailSenderConfig, the operator validates the configuration settings ✔️. When a new Email CRD is created, the operator:

- Fetches the email sending configuration details from the referenced EmailSenderConfig 📥
- Uses the appropriate provider API (MailerSend or Mailgun) based on the configuration to send the email 📧
- Updates the status of the Email CRD with the delivery outcome (success/failure) and relevant information ℹ️

### Security

The operator securely retrieves the API token for email providers from Kubernetes secrets to avoid storing sensitive information in plain text 🛡️.

Stay tuned! 🌊 This repository will be populated with the source code, deployment manifests, architecture diagrams, and documentation as development progresses 🏗️.

Architecture Diagram

+--------------------+                 +-----------------------+
| User               |                 | Kubernetes API Server |
+--------------------+                 +-----------------------+
                     |
                     | creates/updates CRDs
                     v
+--------------------+                 +-----------------------+
| Custom Resource     |                 |      etcd             |
| Definitions (CRDs) | (Email,             +-----------------------+
| (YAML files)       |  EmailSenderConfig) |
+--------------------+                 |  (cluster state)       |
                     |
                     | watches for changes
                     v
+--------------------+                 +-----------------------+
| Kubernetes Operator |                 |  Controller Manager  |
| (Go application)   |                 +-----------------------+
+--------------------+                 |   (reconciliation loop) |
                     |
                     | interacts with CRDs
                     v
+--------------------+                 +-----------------------+
| MailerSend/Mailgun | (via their respective | Email Sending Providers|
| APIs               |  APIs)               | (external services)   |
+--------------------+                 +-----------------------+
                     |
                     | sends email based on config
                     v
+--------------------+                 +-----------------------+
| Kubernetes Cluster |                 |  Worker Nodes         |
| (infrastructure)   |                 +-----------------------+
+--------------------+                 |   (run containers)    |
                     |
+--------------------+
| Email Application  |  (can be any application 
| (containers)       |   that needs email sending)
+--------------------+

### Explanation of the diagram:

- User: Interacts with the Kubernetes cluster by creating and updating CRDs (Email and EmailSenderConfig) in YAML format.

- Custom Resource Definitions (CRDs): These YAML files define the desired state of the application (email sending configuration and email messages).

- Kubernetes API Server: The central point of communication for the Kubernetes cluster. It receives requests from the user (CRD creation/updates) and forwards them to relevant components.

- etcd: A distributed key-value store that stores the cluster state, including the CRD definitions.

- Kubernetes Operator (Go application): This is our custom-built program that watches for changes to the CRDs in etcd. It runs as a pod within the Kubernetes cluster.

- Controller Manager: This Kubernetes component manages all controllers, including our custom operator. The operator runs a reconciliation loop that continuously checks the desired state of the application (defined in CRDs) against the actual state in the cluster.

- MailerSend/Mailgun APIs: The operator interacts with the email provider APIs (MailerSend or Mailgun) based on the configuration in the EmailSenderConfig CRD. It sends emails using the provided API token and recipient information from the Email CRD.

- Email Sending Providers (external services): These are external services like MailerSend or Mailgun that handle sending emails.

- Kubernetes Cluster (infrastructure): This represents the underlying infrastructure of the Kubernetes cluster, including worker nodes that run the containerized application.

- Email Application (containers): This can be any application running in the Kubernetes cluster that needs to send emails. The operator interacts with this application indirectly through the Email CRD.

##### This architecture diagram depicts how the various components work together to enable automated email sending through a Kubernetes operator.