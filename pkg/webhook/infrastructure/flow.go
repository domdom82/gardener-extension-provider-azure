// Copyright (c) 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package infrastructure

import (
	"context"
	"fmt"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	extensionscontextwebhook "github.com/gardener/gardener/extensions/pkg/webhook/context"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/gardener/gardener-extension-provider-azure/pkg/azure"
)

type flowMutator struct {
	logger logr.Logger
	client client.Client
}

// NewFlowMutator returns a new Infrastructure flowMutator that uses mutateFunc to perform the mutation.
func NewFlowMutator(mgr manager.Manager, logger logr.Logger) extensionswebhook.Mutator {
	return &flowMutator{
		client: mgr.GetClient(),
		logger: logger,
	}
}

// Mutate mutates the given object on creation and adds the annotation `aws.provider.extensions.gardener.cloud/use-flow=true`
// if the seed has the label `aws.provider.extensions.gardener.cloud/use-flow` == `new`.
func (m *flowMutator) Mutate(ctx context.Context, new, old client.Object) error {
	if old != nil || new.GetDeletionTimestamp() != nil {
		return nil
	}

	newInfra, ok := new.(*extensionsv1alpha1.Infrastructure)
	if !ok {
		return fmt.Errorf("could not mutate: object is not of type Infrastructure")
	}

	gctx := extensionscontextwebhook.NewGardenContext(m.client, new)
	cluster, err := gctx.GetCluster(ctx)
	if err != nil {
		return err
	}

	// force terraformer case
	if metav1.HasAnnotation(cluster.Shoot.ObjectMeta, azure.AnnotationKeyUseTF) {
		return nil
	}

	if cluster.Seed.Labels[azure.SeedLabelKeyUseFlow] == azure.SeedLabelUseFlowValueNew {
		if newInfra.Annotations == nil {
			newInfra.Annotations = map[string]string{}
		}
		newInfra.Annotations[azure.AnnotationKeyUseFlow] = "true"
	}

	return nil
}