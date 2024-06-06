package v1alpha1

import (

	// MailerSend library (assuming this is needed)
	k8ssecret "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EmailSenderConfigSpec defines the desired state of EmailSenderConfig
type EmailSenderConfigSpec struct {
	// APITokenSecretRef is a reference to the Kubernetes secret containing the MailerSend API token
	APITokenSecretRef k8ssecret.SecretReference `json:"apiTokenSecretRef"`

	// SenderEmail is the email address used for sending emails
	SenderEmail string `json:"senderEmail"`
}

// EmailSenderConfigStatus defines the observed state of EmailSenderConfig
type EmailSenderConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - observed state of config (optional)
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EmailSenderConfig is the Schema for the email sender configurations API
type EmailSenderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmailSenderConfigSpec   `json:"spec,omitempty"`
	Status EmailSenderConfigStatus `json:"status,omitempty"`
}

// EmailSpec defines the desired state of Email
type EmailSpec struct {
	// SenderConfigRef is a reference to the EmailSenderConfig resource
	SenderConfigRef string `json:"senderConfigRef"`

	// RecipientEmail is the email address of the recipient
	RecipientEmail string `json:"recipientEmail"`

	// Subject is the subject line of the email
	Subject string `json:"subject"`

	// Body is the email message body (HTML format)
	Body string `json:"body,omitempty"`
}

// EmailStatus defines the observed state of Email
type EmailStatus struct {
	// DeliveryStatus indicates the email delivery status (e.g., "Sent", "Failed")
	DeliveryStatus string `json:"deliveryStatus,omitempty"`

	// MessageID is the message ID assigned by MailerSend (if successful)
	MessageID string `json:"messageID,omitempty"`

	// Error contains any error message encountered during sending (if failed)
	Error string `json:"error,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Email is the Schema for the emails API
type Email struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmailSpec   `json:"spec,omitempty"`
	Status EmailStatus `json:"status,omitempty"`
}

type EmailList struct {
	metav1.TypeMeta `json:",inline"`
	// Omit the json tag for Items field as metav1.ListMeta already defines it
	metav1.ListMeta
	Items []Email `json:"items"`
}
