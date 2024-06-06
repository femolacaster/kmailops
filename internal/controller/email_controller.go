/*
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
*/

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	emailsv1alpha1 "github.com/femolacaster/kmailops/api/v1alpha1"
	mailersend "github.com/mailersend/mailersend-go"
	"github.com/mailgun/mailgun-go/v4"
	k8ssecret "k8s.io/api/core/v1"
)

// EmailReconciler reconciles a Email object
type EmailReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=emails.example.com,resources=emails,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=emails.example.com,resources=emails/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=emails.example.com,resources=emails/finalizers,verbs=update

// / Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *EmailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Get the object
	obj := &emailsv1alpha1.Email{}
	err := r.Get(ctx, req.NamespacedName, obj)
	if err != nil {
		log.Error(err, "Failed to get object")
		return ctrl.Result{}, err
	}

	// Handle Email object
	if obj.Kind == "Email" {
		// Get referenced EmailSenderConfig
		config := &emailsv1alpha1.EmailSenderConfig{}
		err := r.Get(ctx, types.NamespacedName{Name: obj.Spec.SenderConfigRef, Namespace: obj.Namespace}, obj)
		if err != nil {
			log.Error(err, "Failed to get referenced EmailSenderConfig")
			return ctrl.Result{}, err
		}

		// Retrieve API token from Secret
		secret := &k8ssecret.Secret{}
		err = r.Get(ctx, types.NamespacedName{Name: config.Spec.APITokenSecretRef.Name, Namespace: obj.Namespace}, secret)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Error(err, "Secret containing API token not found")
			} else {
				log.Error(err, "Failed to get secret")
			}
			return ctrl.Result{}, err
		}

		// Choose email provider based on configuration (replace with your logic)
		var provider interface{}
		if providerAnnotation, ok := obj.Annotations["email.provider"]; ok {
			if providerAnnotation == "mailersend" {
				provider = mailersend.NewMailersend(string(secret.Data["api-key"])) // Assuming key is "api-token" in secret
			} else if providerAnnotation == "mailgun" {
				// Configure Mailgun client with retrieved API key from secret (adjust based on Mailgun library)
				domain := "your-domain.com"                                           // Replace with your Mailgun domain
				provider = mailgun.NewMailgun(domain, string(secret.Data["api-key"])) // Assuming key is "api-key" in secret
			} else {
				log.Error(fmt.Errorf("unsupported provider: %s", providerAnnotation), "Invalid email provider specified")
				return ctrl.Result{}, err
			}
		} else {
			// Use default provider if not specified (replace with your choice)
			provider = mailersend.NewMailersend(string(secret.Data["api-key"])) // Assuming key is "api-token" in secret
		}

		// Send email using chosen provider
		err = sendEmail(ctx, provider, obj.Spec.RecipientEmail, obj.Spec.Subject, obj.Spec.Body, config.Spec.SenderEmail)
		if err != nil {
			log.Error(err, "Failed to send email")
			obj.Status.DeliveryStatus = "Failed"
			obj.Status.Error = err.Error()
		} else {
			// Mail sent successfully, update status
			obj.Status.DeliveryStatus = "Sent"
		}

		// Update Email object with status
		err = r.Status().Update(ctx, obj)
		if err != nil {
			log.Error(err, "Failed to update Email status")
			return ctrl.Result{}, err
		}
	}

	// Handle EmailSenderConfig object (add logic for logging creation/update)

	return ctrl.Result{}, nil
}

func sendEmail(ctx context.Context, provider interface{}, recipientEmail, subject, body, senderEmail string) error {
	switch p := provider.(type) {
	case *mailersend.Mailersend:
		// Define recipient details
		recipients := []mailersend.Recipient{
			{
				Name:  "", // Optional recipient name
				Email: recipientEmail,
			},
		}

		// Create a new message
		message := p.Email.NewMessage()

		// Set message details
		message.SetFrom(mailersend.From{
			Name:  senderEmail, // Assuming sender name is same as email
			Email: senderEmail,
		})
		message.SetRecipients(recipients)
		message.SetSubject(subject)
		message.SetHTML(body)
		message.SetText("") // Optional plain text body

		// Send the message and handle error
		_, err := p.Email.Send(ctx, message)
		return err

	case *mailgun.Mailgun:
		// Configure Mailgun message details
		message := (*p).NewMessage(senderEmail, subject, body, recipientEmail)
		// Optional: Set additional message details like ReplyTo, CC, BCC
		// message.SetReplyTo("reply@yourdomain.com")
		// message.SetCc("cc@example.com")
		// message.SetBcc("bcc@example.com")

		_, _, err := (*p).Send(ctx, message)
		return err

	default:
		return fmt.Errorf("unsupported provider: %T", provider)
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *EmailReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&emailsv1alpha1.Email{}).
		Complete(r)
}
