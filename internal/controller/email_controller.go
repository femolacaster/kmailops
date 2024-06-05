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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	emailsv1alpha1 "github.com/femolacaster/kmailops/api/v1alpha1"
)

// EmailReconciler reconciles a Email object
type EmailReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=emails.example.com,resources=emails,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=emails.example.com,resources=emails/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=emails.example.com,resources=emails/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Email object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *EmailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Get the Email resource
	email := &emailsv1alpha1.Email{}
	err := r.Get(ctx, req.NamespacedName, email)
	if err != nil {
		log.Error(err, "Failed to get Email resource")
		return ctrl.Result{}, err
	}

	// Get the referenced EmailSenderConfig
	senderConfig := &emailsv1alpha1.EmailSenderConfig{}
	err = r.Get(ctx, types.NamespacedName{Name: email.Spec.SenderConfigRef}, senderConfig)
	if err != nil {
		log.Error(err, "Failed to get EmailSenderConfig resource")
		return ctrl.Result{}, err
	}

	// Retrieve API token from Secret (replace with actual logic)
	var apiToken string
	// ... logic to get apiToken from secret referenced by senderConfig.Spec.ApiTokenSecretRef

	// Create MailerSend client with retrieved token
	mailerSendClient := mailsend.NewClient(apiToken)

	// ... continue with email sending logic
}

// SetupWithManager sets up the controller with the Manager.
func (r *EmailReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&emailsv1alpha1.Email{}).
		Complete(r)
}
