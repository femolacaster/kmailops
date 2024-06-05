package crd

import (
	v1 "k8s.io/api/core/v1" // Import the core v1 API group
)

// Define CRD for EmailSenderConfig
type EmailSenderConfigSpec struct {
	APITokenSecretRef v1.LocalObjectReference `json:"apiTokenSecretRef"`
	SenderEmail       string                  `json:"senderEmail"`
}

// Define CRD for Email
type EmailSpec struct {
	SenderConfigRef string `json:"senderConfigRef"`
	RecipientEmail  string `json:"recipientEmail"`
	Subject         string `json:"subject"`
	Body            string `json:"body"`
}
