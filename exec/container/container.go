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

package container

import (
	"github.com/chaosblade-io/chaosblade-operator/channel"
	"github.com/chaosblade-io/chaosblade-operator/exec/model"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
)

type ResourceModelSpec struct {
	model.BaseResourceExpModelSpec
}

// NewResourceModelSpec returns the container model spec
func NewResourceModelSpec(client *channel.Client) model.ResourceExpModelSpec {
	modelSpec := &ResourceModelSpec{
		model.NewBaseResourceExpModelSpec("container", client),
	}
	dockerModelSpecs := NewDockerSubResourceModelSpec(client).ExpModels()

	spec.AddFlagsToModelSpec(getResourceFlags, dockerModelSpecs...)
	modelSpec.RegisterExpModels(dockerModelSpecs...)
	addActionExamples(modelSpec)
	return modelSpec
}

func addActionExamples(modelSpec *ResourceModelSpec) {
	for _, expModelSpec := range modelSpec.ExpModelSpecs {
		for _, commandSpec := range expModelSpec.Actions() {
			if expModelSpec.Name() == "network" {
				commandSpec.SetLongDesc("Kubernetes Container network scenes, same as the network scenes of the underlying resources")
			} else if expModelSpec.Name() == "process" {
				commandSpec.SetLongDesc("Kubernetes Container process scenes, same as the process scenes of the underlying resources")
			} else if expModelSpec.Name() == "" {
				commandSpec.SetLongDesc("For scenarios of container resources themselves under Kubernetes, such as deleting containers, it is important to note that to execute container scenarios, Pod must be determined first, so Pod related parameters need to be configured")
			}
		}
	}
}

func getResourceFlags() []spec.ExpFlagSpec {
	coverageFlags := model.GetResourceCoverageFlags()
	commonFlags := model.GetResourceCommonFlags()
	containerFlags := model.GetContainerFlags()
	return append(append(coverageFlags, commonFlags...), containerFlags...)
}
