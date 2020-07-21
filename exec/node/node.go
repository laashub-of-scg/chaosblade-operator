/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package node

import (
	"github.com/chaosblade-io/chaosblade-operator/channel"
	"github.com/chaosblade-io/chaosblade-operator/exec/model"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
)

type ResourceModelSpec struct {
	model.BaseResourceExpModelSpec
}

func NewResourceModelSpec(client *channel.Client) model.ResourceExpModelSpec {
	modelSpec := &ResourceModelSpec{
		model.NewBaseResourceExpModelSpec("node", client),
	}
	osModelSpecs := NewOSSubResourceModelSpec(client).ExpModels()
	selfModelSpec := NewSelfExpModelCommandSpec()

	expModelSpecs := append(osModelSpecs, selfModelSpec)
	spec.AddFlagsToModelSpec(getResourceFlags, expModelSpecs...)
	modelSpec.RegisterExpModels(osModelSpecs...)
	addActionExamples(modelSpec)
	return modelSpec
}

func addActionExamples(modelSpec *ResourceModelSpec) {
	for _, expModelSpec := range modelSpec.ExpModelSpecs {
		for _, commandSpec := range expModelSpec.Actions() {
			if expModelSpec.Name() == "network" {
				commandSpec.SetLongDesc("The kubernetes Node network scenario is the same as the network scenario of the underlying resource")
			} else if expModelSpec.Name() == "cpu" {
				commandSpec.SetLongDesc("The kubernetes Node cpu scenario is the same as the cpu scenario of the underlying resource")
			} else if expModelSpec.Name() == "process" {
				commandSpec.SetLongDesc("The kubernetes Node process scenario is the same as the process scenario of the underlying resource")
			} else if expModelSpec.Name() == "disk" {
				commandSpec.SetLongDesc("The kubernetes Node disk scenario is the same as the disk scenario of the underlying resource")
			}
		}
	}
}

func getResourceFlags() []spec.ExpFlagSpec {
	coverageFlags := model.GetResourceCoverageFlags()
	return append(coverageFlags, model.ResourceNamesFlag, model.ResourceLabelsFlag)
}

func NewSelfExpModelCommandSpec() spec.ExpModelCommandSpec {
	return &SelfExpModelCommandSpec{
		spec.BaseExpModelCommandSpec{
			ExpFlags: []spec.ExpFlagSpec{},
			ExpActions: []spec.ExpActionCommandSpec{
				// TODO
				//NewCordonActionCommandSpec(),
			},
		},
	}
}

type SelfExpModelCommandSpec struct {
	spec.BaseExpModelCommandSpec
}

func (*SelfExpModelCommandSpec) Name() string {
	return "node"
}

func (*SelfExpModelCommandSpec) ShortDesc() string {
	return "Node resource experiment for itself, for example cpu load"
}

func (*SelfExpModelCommandSpec) LongDesc() string {
	return "Node resource experiment for itself, for example cpu load"
}

func (*SelfExpModelCommandSpec) Example() string {
	return "blade c k8s node-cpu load --evict-count 1 --kubeconfig ~/.kube/config --names cn-hangzhou.192.168.0.205"
}
