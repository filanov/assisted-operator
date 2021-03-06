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

	"github.com/go-openapi/swag"

	"github.com/openshift/assisted-service/client/installer"
	"github.com/openshift/assisted-service/models"

	servicev1 "github.com/filanov/assisted-operator/api/v1"
	"github.com/go-logr/logr"
	aiclient "github.com/openshift/assisted-service/client"
	"github.com/openshift/assisted-service/pkg/auth"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=service.example.com,resources=clusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=service.example.com,resources=clusters/status,verbs=get;update;patch

func (r *ClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	var (
		ctx = context.Background()
		log = r.Log.WithValues("cluster", req.NamespacedName)
	)

	// your logic here
	cluster := &servicev1.Cluster{}
	err := r.Get(ctx, req.NamespacedName, cluster)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.Info("Cluster resource not found, no work to do")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get cluster")
		return ctrl.Result{}, err
	}

	service := &servicev1.AssistedService{}
	err = r.Get(ctx, client.ObjectKey{
		Namespace: req.Namespace,
		Name:      cluster.Spec.AssistedService,
	}, service)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.Info("AssistedService resource not found, no work to do")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get related AssistedService")
		return ctrl.Result{}, err
	}

	cfg := aiclient.Config{
		URL: &url.URL{
			Scheme: aiclient.DefaultSchemes[1],
			Host:   service.Spec.URL,
			Path:   aiclient.DefaultBasePath,
		},
		AuthInfo: auth.UserAuthHeaderWriter("Bearer " + service.Spec.Token),
	}
	c := aiclient.New(cfg)

	reply, err := c.Installer.RegisterCluster(ctx, &installer.RegisterClusterParams{
		NewClusterParams: &models.ClusterCreateParams{
			Name:       swag.String(cluster.Spec.Name),
			PullSecret: cluster.Spec.PullSecret,
		},
	})

	if err != nil {
		log.Error(err, "failed to create cluster")
		return ctrl.Result{}, err
	}

	log.Info(fmt.Sprintf("succesfully created cluster %+v", reply))
	return ctrl.Result{}, nil
}

func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&servicev1.Cluster{}).
		Complete(r)
}
