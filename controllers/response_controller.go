/*
Copyright 2022.

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

package controllers

import (
	"context"
	"errors"
	"github.com/mrogers950/sporc/controllers/crl"
	v1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	sporcv1alpha1 "github.com/mrogers950/sporc/api/v1alpha1"
)

// ResponseReconciler reconciles a Response object
type ResponseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=sporc.example.com,resources=responses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=sporc.example.com,resources=responses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=sporc.example.com,resources=responses/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Response object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *ResponseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// your logic here!!!!!!!!!!!!!!!!!!!
	config := &sporcv1alpha1.ResponseConfig{}
	if err := r.Get(ctx, req.NamespacedName, config); err != nil {
		if !apierr.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		// Create a default responseConfig
		if createErr := r.Create(ctx, newResponseConfig(req.Name, req.Namespace)); createErr != nil {
			return ctrl.Result{}, createErr
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// discover the configMap
	cm := &v1.ConfigMap{}
	if err := r.Client.Get(ctx, config.ConfigMapNSN(), cm); err != nil {
		return ctrl.Result{}, err
	}

	// does it have crl.pem?
	if _, has := cm.Data["crl.pem"]; !has {
		return ctrl.Result{}, errors.New("configMap does not have key \"crl.pem\"")
	}

	// read the CRL
	newCrl, err := crl.CollectCRLResponses([]byte(cm.Data["crl.pem"]))
	if err != nil {
		return ctrl.Result{}, err
	}

	// fetch current responses
	responses := &sporcv1alpha1.ResponseList{}
	listOpts := client.ListOptions{
		// should we label?
		LabelSelector: labels.SelectorFromSet(map[string]string{
			"sporc": config.ConfigMapNSN().Name,
		}),
		Namespace: config.ConfigMapNSN().Namespace,
	}

	listErr := r.List(ctx, responses, &listOpts)
	if listErr != nil && !apierr.IsNotFound(listErr) {
		return ctrl.Result{}, listErr
	}

	if apierr.IsNotFound(listErr) {
		// none found, create
		// return createNewResponses(r, newCrl, responses)
	}

	// compare CRL updates with current responses. (right now check for added revocations)
	if hasNewRevocations(r, newCrl, responses) {
		// return updateResponses(r, newCrl, responses)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResponseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&sporcv1alpha1.Response{}).
		Complete(r)
}

func hasNewRevocations(r *ResponseReconciler, newCrl, current *sporcv1alpha1.ResponseList) bool {
	return len(newCrl.Items) != len(current.Items)
}

//func createNewResponses(ctx context.Context, r *ResponseReconciler, newCrl, current *sporcv1alpha1.ResponseList) (ctrl.Result, error) {
//
// }

func updateResponses(ctx context.Context, r *ResponseReconciler, newCrl, current *sporcv1alpha1.ResponseList) (ctrl.Result, error) {
	if len(newCrl.Items) == 0 {
		return ctrl.Result{}, nil
	}

	for i, _ := range newCrl.Items {
		nc := newCrl.Items[i]
		resp := &sporcv1alpha1.Response{
			Status: nc.Status,
		}

		if createErr := r.Create(ctx, resp, &client.CreateOptions{}); createErr != nil {
			if apierr.IsAlreadyExists(createErr) {

			}
			return ctrl.Result{}, createErr
		}
	}
	return ctrl.Result{}, nil
}
