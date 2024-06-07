# KMailOps 📧

## Description
### A Kubernetes Email Operator with Multiple Providers

This project builds a Kubernetes operator written in Go that manages custom resources (CRDs) for configuring and sending emails through transactional email providers like MailerSend 📨 and/or Mailgun 🔫.

### Project Goals

- Simplify email sending in Kubernetes: Automate email configuration and delivery for applications running in Kubernetes clusters 🤖
- Support multiple providers: Integrate with both MailerSend and Mailgun APIs for sending emails 🤝
- CRD-based configuration: Define email sending configurations and messages using CRDs for a declarative approach


## Getting Started

### Prerequisites
- go version v1.20.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/kmailops:tag
```

**NOTE:** This image ought to be published in the personal registry you specified. 
And it is required to have access to pull the image from the working environment. 
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/kmailops:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin 
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Testing Code

1. Create Secret for your provider's api in base64

```
apiVersion: v1
kind: Secret
metadata:
  name: mailgun-api-token-secret
  namespace: default  # Replace with the namespace where you want the secret

type: opaque

stringData:
  api-key: [provider-apikey-in-base64]
```

2. Use the samples in ```crd/samples``` to create E-mail and E-mail Config CRD. Note that providers are specified in the annotations.emailprovider.

3. Alternatively build this code and edit the controller deployment in controller-deployment.yaml to suit. controller-deployment currently has this configurations:

```
name: email-operator
namespace: default
service_account: email-operator
resources: 0.8 CPU, 256Mi memory
replicas: 1
```

## Technical Specifications

_Custom Resource Definitions (CRDs)_

- `EmailSenderConfig`: This CRD defines the configuration details for an email provider, including the API token (stored in a Kubernetes secret) and sender email address 🔑
- `Email`: This CRD defines an email message to be sent, referencing the sender configuration and specifying the recipient email, subject, and body content. The status field of the Email CRD will reflect the delivery status (success/failure), message ID (if successful), and any error messages ⚠️

## Operator Implementation

The operator is written in Go and leverages the Operator SDK framework for ease of development 🛠️. It watches for changes to the Email and EmailSenderConfig CRDs using Kubernetes controllers 👀. Upon creation or update of an EmailSenderConfig, the operator validates the configuration settings ✔️. When a new Email CRD is created, the operator:

- Fetches the email sending configuration details from the referenced EmailSenderConfig 📥
- Uses the appropriate provider API (MailerSend or Mailgun) based on the configuration to send the email 📧
- Updates the status of the Email CRD with the delivery outcome (success/failure) and relevant information ℹ️

## Security

The operator securely retrieves the API token for email providers from Kubernetes secrets to avoid storing sensitive information in plain text 🛡️.

## TO-DO

1. Refactor Code to functions. Use Regex for domain name grabbing for mailgun to prevent hardcoding.
2. Refactor sendMail case statements to individual functions.
3. Log more operator actions for observability.
Use advanced statuses to check statuses of the operator
4. Test dockerfile deployment on multiple distributions
5. Test operator with a deployment and service manifest
6. Deploy and integrate kube-prometheus for advanced monitoring.
7. Add more metrics endpoints to code to monitor with kube-prometheus
8. Version this operator API


## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

