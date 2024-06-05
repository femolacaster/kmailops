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

// Define a sample controller using the Operator SDK (pseudocode)
type EmailReconciler struct {
	// ... other fields
	
  }


func (r *EmailReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	// Fetch the Email CRD instance
	email := &k8s.apiextensions.v1.CustomResource{}
	err := r.client.Get(ctx, request.NamespacedName, email)
	if err !=
}
