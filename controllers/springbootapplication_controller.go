/*
Copyright 2020 the original author or authors.

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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1alpha1 "github.com/spring-cloud-incubator/mononoke/api/v1alpha1"
)

// SpringBootApplicationReconciler reconciles a SpringBootApplication object
type SpringBootApplicationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.mononoke.local,resources=springbootapplications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.mononoke.local,resources=springbootapplications/status,verbs=get;update;patch

func (r *SpringBootApplicationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("springbootapplication", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *SpringBootApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.SpringBootApplication{}).
		Complete(r)
}
