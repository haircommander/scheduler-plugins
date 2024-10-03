/*
Copyright 2024 The Kubernetes Authors.

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

package usernssupported

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// UserNSSupported plugin
type UserNSSupported struct {
	handle framework.Handle
}

var _ framework.FilterPlugin = &UserNSSupported{}

// Name is the name of the plugin used in the Registry and configurations.
const Name = "UserNSSupported"

func (uns *UserNSSupported) Name() string {
	return Name
}

func (uns *UserNSSupported) Filter(ctx context.Context, cycleState *framework.CycleState, pod *corev1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	if nodeInfo.Node() == nil {
		return framework.NewStatus(framework.Error, "node not found")
	}

	if nodeInfo.Node().Status.RuntimeHandlers == nil {
		return framework.NewStatus(framework.Error, "node.Status.RuntimeHandlers not declared")
	}

	rhs := nodeInfo.Node().Status.RuntimeHandlers

	usernsSupported := false
	for _, rh := range rhs {
		if rh.Features == nil {
			continue
		}
		if rh.Features.UserNamespaces != nil && *rh.Features.UserNamespaces {
			usernsSupported = true
		}
	}
	if !usernsSupported {
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("Pod %v requested a user namespace, but node %v does not support them", pod.Name, nodeInfo.Node().Name))
	}
	return nil
}

// New initializes a new plugin and returns it.
func New(_ context.Context, _ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	return &UserNSSupported{handle: h}, nil
}
