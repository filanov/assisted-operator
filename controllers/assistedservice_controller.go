/*
Copyright 2020.

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
	"fmt"
	"net/url"

	servicev1 "github.com/filanov/assisted-operator/api/v1"
	"github.com/go-logr/logr"
	aiclient "github.com/openshift/assisted-service/client"
	"github.com/openshift/assisted-service/client/installer"
	"github.com/openshift/assisted-service/pkg/auth"
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AssistedServiceReconciler reconciles a AssistedService object
type AssistedServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=service.example.com,resources=assistedservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=service.example.com,resources=assistedservices/status,verbs=get;update;patch

func (r *AssistedServiceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	var (
		ctx = context.Background()
		log = r.Log.WithValues("assistedservice", req.NamespacedName)
	)

	// your logic here
	service := &servicev1.AssistedService{}
	err := r.Get(ctx, req.NamespacedName, service)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.Info("Assisted-service resource not found, no work to do")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get AssistedService")
		return ctrl.Result{}, err
	}

	log.Info(fmt.Sprintf("found service URL: %s", service.Spec.URL))
	cfg := aiclient.Config{
		URL: &url.URL{
			Scheme: aiclient.DefaultSchemes[1],
			Host:   service.Spec.URL,
			Path:   aiclient.DefaultBasePath,
		},
		AuthInfo: auth.UserAuthHeaderWriter("Bearer " + service.Spec.Token),
	}
	c := aiclient.New(cfg)
	reply, err := c.Installer.ListClusters(ctx, &installer.ListClustersParams{})
	if err != nil {
		log.Error(err, "failed to connect to service")
		return ctrl.Result{}, errors.Wrapf(err, "failed to list clusters from %s", service.Spec.URL)
	}
	clusters := reply.GetPayload()
	for _, cluster := range clusters {
		log.Info(fmt.Sprintf("found cluster ID: %s name: %s nodes: %d",
			cluster.ID.String(), cluster.Name, len(cluster.Hosts)))
	}

	return ctrl.Result{}, nil
}

func (r *AssistedServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&servicev1.AssistedService{}).
		Complete(r)
}
