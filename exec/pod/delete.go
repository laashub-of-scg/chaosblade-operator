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

package pod

import (
	"context"

	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"

	"github.com/chaosblade-io/chaosblade-operator/channel"
	"github.com/chaosblade-io/chaosblade-operator/exec/model"
	"github.com/chaosblade-io/chaosblade-operator/pkg/apis/chaosblade/v1alpha1"
)

type DeletePodActionSpec struct {
	spec.BaseExpActionCommandSpec
}

func NewDeletePodActionSpec(client *channel.Client) spec.ExpActionCommandSpec {
	return &DeletePodActionSpec{
		spec.BaseExpActionCommandSpec{
			ActionMatchers: []spec.ExpFlagSpec{},
			ActionFlags: []spec.ExpFlagSpec{
				&spec.ExpFlag{
					Name:   "random",
					Desc:   "Randomly select pod",
					NoArgs: true,
				},
			},
			ActionExecutor: &DeletePodActionExecutor{client: client},
			ActionExample: spec.Example{
				ExampleCommands: [] spec.ExampleCommand{
					{
						Annotation: "Deletes the POD under the specified default namespace that is app=guestbook",
						Command: "blade create k8s pod-pod delete --labels app=guestbook --namespace default --evict-count 2 --kubeconfig config",
					},

					{
						Annotation: "Query experimental status",
						Command: "blade query k8s create 4d3caa0a99c3b2dd --kubeconfig config",
					},
				},
			},
		},
	}
}

func (*DeletePodActionSpec) Name() string {
	return "delete"
}

func (*DeletePodActionSpec) Aliases() []string {
	return []string{}
}

func (*DeletePodActionSpec) ShortDesc() string {
	return "Delete pods"
}

func (d *DeletePodActionSpec) LongDesc() string {
	if d.ActionLongDesc != "" {
		return d.ActionLongDesc
	}
	return "Kubernetes Pod resource itself scenario, such as removing Pod"
}

type DeletePodActionExecutor struct {
	client *channel.Client
}

func (*DeletePodActionExecutor) Name() string {
	return "delete"
}

func (*DeletePodActionExecutor) SetChannel(channel spec.Channel) {
}

func (d *DeletePodActionExecutor) Exec(uid string, ctx context.Context, model *spec.ExpModel) *spec.Response {
	if _, ok := spec.IsDestroy(ctx); ok {
		return d.destroy(ctx, model)
	} else {
		return d.create(ctx, model)
	}
}

func (d *DeletePodActionExecutor) create(ctx context.Context, expModel *spec.ExpModel) *spec.Response {
	podObjectMetaList, err := model.ExtractPodObjectMetasFromContext(ctx)
	if err != nil {
		return spec.ReturnFailWitResult(spec.Code[spec.IllegalParameters], err.Error(),
			v1alpha1.CreateFailExperimentStatus(err.Error(), nil))
	}
	statuses := make([]v1alpha1.ResourceStatus, 0)
	success := false
	for _, meta := range podObjectMetaList {
		status := v1alpha1.ResourceStatus{
			Uid:      meta.Uid,
			Name:     meta.Name,
			Kind:     v1alpha1.PodKind,
			NodeName: meta.NodeName,
		}
		objectMeta := metav1.ObjectMeta{Name: meta.Name, Namespace: meta.Namespace}
		err := d.client.Delete(context.TODO(), &v1.Pod{ObjectMeta: objectMeta})
		if err != nil {
			logrus.Warningf("delete pod %s err, %v", meta.Name, err)
			status = status.CreateFailResourceStatus(err.Error())
		} else {
			status = status.CreateSuccessResourceStatus()
			success = true
		}
		statuses = append(statuses, status)
	}
	var experimentStatus v1alpha1.ExperimentStatus
	if success {
		experimentStatus = v1alpha1.CreateSuccessExperimentStatus(statuses)
	} else {
		experimentStatus = v1alpha1.CreateFailExperimentStatus("see resStatuses for details", statuses)
	}
	return spec.ReturnResultIgnoreCode(experimentStatus)
}

func (d *DeletePodActionExecutor) destroy(ctx context.Context, expModel *spec.ExpModel) *spec.Response {
	expObjectMetasMaps, err := model.ExtractNodeNameExpObjectMetasMapFromContext(ctx)
	if err != nil {
		spec.ReturnFailWitResult(spec.Code[spec.IllegalParameters], err.Error(),
			v1alpha1.CreateFailExperimentStatus(err.Error(), nil))
	}
	experimentStatus := v1alpha1.CreateDestroyedExperimentStatus([]v1alpha1.ResourceStatus{})
	statuses := experimentStatus.ResStatuses
	for nodeName, objectMetas := range expObjectMetasMaps {
		for _, objectMeta := range objectMetas {
			status := v1alpha1.ResourceStatus{
				Id:       objectMeta.Id,
				Uid:      objectMeta.Uid,
				Name:     objectMeta.Name,
				Kind:     v1alpha1.PodKind,
				NodeName: nodeName,
				State:    v1alpha1.DestroyedState,
				Success:  true,
			}
			statuses = append(statuses, status)
		}
	}
	experimentStatus.ResStatuses = statuses
	return spec.ReturnResultIgnoreCode(experimentStatus)
}
