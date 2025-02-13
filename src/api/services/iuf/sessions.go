/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a
 *  copy of this software and associated documentation files (the "Software"),
 *  to deal in the Software without restriction, including without limitation
 *  the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *  and/or sell copies of the Software, and to permit persons to whom the
 *  Software is furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included
 *  in all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 *  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 *  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 *  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *  OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package services_iuf

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"time"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func (s iufService) GetSession(sessionName string) (iuf.Session, error) {
	rawConfigMapData, err := s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Get(
			context.TODO(),
			sessionName,
			v1.GetOptions{},
		)
	if err != nil {
		s.logger.Error(err)
		return iuf.Session{}, err
	}

	res, err := s.ConfigMapDataToSession(rawConfigMapData.Data[LABEL_SESSION])
	if err != nil {
		s.logger.Error(err)
		return res, err
	}
	return res, err
}

func (s iufService) ListSessions(activityName string) ([]iuf.Session, error) {
	rawConfigMapList, err := s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		List(
			context.TODO(),
			v1.ListOptions{
				LabelSelector: fmt.Sprintf("type=%s,%s=%s", LABEL_SESSION, LABEL_ACTIVITY_REF, activityName),
			},
		)
	if err != nil {
		s.logger.Error(err)
		return []iuf.Session{}, err
	}
	var res []iuf.Session
	for _, rawConfigMap := range rawConfigMapList.Items {
		tmp, err := s.ConfigMapDataToSession(rawConfigMap.Data[LABEL_SESSION])
		if err != nil {
			s.logger.Error(err)
			return []iuf.Session{}, err
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (s iufService) ConfigMapDataToSession(data string) (iuf.Session, error) {
	var res iuf.Session
	err := json.Unmarshal([]byte(data), &res)
	if err != nil {
		s.logger.Error(err)
		return res, err
	}
	return res, err
}

func (s iufService) CreateSession(session iuf.Session, name string, activity iuf.Activity) (iuf.Session, error) {
	configmap, err := s.iufObjectToConfigMapData(session, name, LABEL_SESSION)
	if err != nil {
		s.logger.Error(err)
		return iuf.Session{}, err
	}
	configmap.Labels[LABEL_ACTIVITY_REF] = activity.Name
	_, err = s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Create(
			context.TODO(),
			&configmap,
			v1.CreateOptions{},
		)
	return session, err
}

func (s iufService) UpdateSessionAndActivity(session iuf.Session) error {
	err := s.UpdateSession(session)
	if err != nil {
		return err
	}

	// if the session update was successful, we also want to update the activity
	s.logger.Infof("Update activity state, session state: %s", session.CurrentState)
	err = s.UpdateActivityStateFromSessionState(session)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	return nil
}

func (s iufService) UpdateSession(session iuf.Session) error {
	configmap, err := s.iufObjectToConfigMapData(session, session.Name, LABEL_SESSION)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	configmap.Labels[LABEL_ACTIVITY_REF] = session.ActivityRef
	// set completed label so metacontroller won't sync it again
	if session.CurrentState == iuf.SessionStateCompleted {
		configmap.Labels["completed"] = "true"
	}
	_, err = s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Update(
			context.TODO(),
			&configmap,
			v1.UpdateOptions{},
		)
	if err != nil {
		// does it even exist? If it doesn't, let's create it instead
		_, err := s.k8sRestClientSet.
			CoreV1().
			ConfigMaps(DEFAULT_NAMESPACE).
			Get(context.TODO(), configmap.Name, v1.GetOptions{})
		if err != nil {
			_, err := s.k8sRestClientSet.
				CoreV1().
				ConfigMaps(DEFAULT_NAMESPACE).
				Create(context.TODO(), &configmap, v1.CreateOptions{})
			if err != nil {
				s.logger.Error(err)
				return err
			}
		} else {
			s.logger.Error(err)
			return err
		}
	}

	return nil
}

func (s iufService) UpdateActivityStateFromSessionState(session iuf.Session) error {
	var activityState iuf.ActivityState
	if session.CurrentState == iuf.SessionStateCompleted {
		activityState = iuf.ActivityStateWaitForAdmin
	} else {
		activityState = iuf.ActivityState(session.CurrentState)
	}
	activity, err := s.GetActivity(session.ActivityRef)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	activity.ActivityState = activityState
	configmap, err := s.iufObjectToConfigMapData(activity, activity.Name, LABEL_ACTIVITY)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	_, err = s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Update(
			context.TODO(),
			&configmap,
			v1.UpdateOptions{},
		)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	// store history
	name := utils.GenerateName(activity.Name)
	iufHistory := iuf.History{
		ActivityState: activityState,
		StartTime:     int32(time.Now().UnixMilli()),
		Name:          name,
		SessionName:   session.Name,
	}
	configmap, err = s.iufObjectToConfigMapData(iufHistory, name, LABEL_HISTORY)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	configmap.Labels[LABEL_ACTIVITY_REF] = activity.Name
	_, err = s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Create(
			context.TODO(),
			&configmap,
			v1.CreateOptions{},
		)

	return err
}

func (s iufService) CreateIufWorkflow(session iuf.Session) (retWorkflow *v1alpha1.Workflow, err error, skipStage bool) {
	myWorkflow, err, skipStage := s.workflowGen(session)
	if err != nil {
		s.logger.Error(err)
		return nil, err, false
	} else if skipStage {
		return nil, nil, true
	}

	res, err := s.workflowClient.CreateWorkflow(context.TODO(), &workflow.WorkflowCreateRequest{
		Namespace: "argo",
		Workflow:  &myWorkflow,
	})
	if err != nil {
		s.logger.Errorf("Creating workflow for: %v FAILED", session)
		s.logger.Error(err)
		return nil, err, false
	}
	return res, nil, false
}

// RunNextStage Runs the next stage in the list of stages to execute.
func (s iufService) RunNextStage(session *iuf.Session) (response iuf.SyncResponse, err error, sessionCompleted bool) {
	// find the current stage in the list of stages, and use the next one
	var currentStage string
	found := false
	if session.CurrentStage != "" {
		for _, stage := range session.InputParameters.Stages {
			if !found {
				if stage == session.CurrentStage {
					found = true
				}
			} else {
				currentStage = stage
				break
			}
		}
	}

	if !found {
		if len(session.InputParameters.Stages) > 0 {
			// Someone updated the input parameters, perhaps. Restart from the beginning because we don't know where we are
			//  anymore
			currentStage = session.InputParameters.Stages[0]
		} else {
			// this session is done because we don't have anything to run
			s.logger.Infof("Session completed. No stages to run")
			return s.SetSessionToCompleted(session)
		}
	} else if currentStage == "" { // we found the last stage
		// this session is done
		return s.SetSessionToCompleted(session)
	}

	stage, err, skipStage := s.RunStage(session, currentStage)
	if skipStage {
		return s.RunNextStage(session)
	} else {
		return stage, err, false
	}
}

func (s iufService) SetSessionToCompleted(session *iuf.Session) (iuf.SyncResponse, error, bool) {
	session.CurrentState = iuf.SessionStateCompleted
	s.logger.Infof("Session completed. Last stage was %s", session.CurrentStage)

	err := s.UpdateSessionAndActivity(*session)
	if err != nil {
		s.logger.Errorf("Error while updating the session %v", err)
		return iuf.SyncResponse{}, err, false
	}

	return iuf.SyncResponse{}, nil, true
}

// RunStage Runs a specific stage for the given session. Creates a new Argo workflow behind the scenes for this stage.
func (s iufService) RunStage(session *iuf.Session, stageToRun string) (ret iuf.SyncResponse, err error, skipStage bool) {
	if stageToRun == "" {
		// this session is done
		s.logger.Infof("No stage specified to run. Last stage was %s and list of all stages are %v",
			session.CurrentStage, session.InputParameters.Stages)
		return iuf.SyncResponse{}, nil, false
	}

	session.CurrentStage = stageToRun
	session.CurrentState = iuf.SessionStateInProgress

	workflow, err, skipStage := s.CreateIufWorkflow(*session)
	if err != nil {
		s.logger.Error(err)
		return iuf.SyncResponse{}, err, skipStage
	} else if !skipStage {
		s.logger.Infof("workflow: %s has been created", workflow.Name)
		session.Workflows = append(session.Workflows, iuf.SessionWorkflow{Id: workflow.Name})
	}

	s.logger.Infof("Update session: %v", session)
	err = s.UpdateSessionAndActivity(*session)
	if err != nil {
		s.logger.Error(err)
		return iuf.SyncResponse{}, err, skipStage
	}

	response := iuf.SyncResponse{
		ResyncAfterSeconds: 5,
	}
	return response, nil, skipStage
}

func (s iufService) ProcessOutput(session *iuf.Session, workflow *v1alpha1.Workflow) error {
	// get activity
	activity, err := s.GetActivity(session.ActivityRef)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	// get tasks we care about (top level dag)
	tasks := workflow.Spec.Templates[0].DAG.Tasks
	switch workflow.Labels["stage_type"] {
	case "product":
		for _, task := range tasks {
			operationName := task.TemplateRef.Name
			nodeStatus := workflow.Status.Nodes.FindByDisplayName(task.Name)
			var productKey string
			for _, param := range nodeStatus.Inputs.Parameters {
				if param.Name == "global_params" {
					var valueJson map[string]interface{}
					json.Unmarshal([]byte(param.Value.String()), &valueJson)
					productManifest := valueJson["product_manifest"].(map[string]interface{})
					currentProduct := productManifest["current_product"].(map[string]interface{})
					manifest := currentProduct["manifest"].(map[string]interface{})
					productKey = s.getProductVersionKeyFromNameAndVersion(manifest["name"].(string), manifest["version"].(string))
					break
				}
			}
			s.logger.Infof("process output of: %s, product: %s, %v", operationName, productKey, nodeStatus.Outputs)
			s.updateActivityOperationOutputFromWorkflow(activity, *session, nodeStatus, operationName, productKey)
		}
		return nil
	case "global":
		// special handling of process media
		if workflow.Labels["stage"] == "process-media" {
			err := s.processOutputOfProcessMedia(&activity, workflow)
			if err != nil {
				s.logger.Error(err)
				return err
			}
			session.Products = activity.Products
			// update activity
			_, err = s.updateActivity(activity)
			if err != nil {
				s.logger.Error(err)
				return err
			}
			return nil
		} else {
			for _, task := range tasks {
				operationName := task.TemplateRef.Name
				nodeStatus := workflow.Status.Nodes.FindByDisplayName(task.Name)
				s.logger.Infof("process output of: %s, %v", operationName, nodeStatus.Outputs)
				s.updateActivityOperationOutputFromWorkflow(activity, *session, nodeStatus, operationName, "")
			}
			return nil
		}
	default:
		return fmt.Errorf("stage_type: %s is not supported", workflow.Labels["stage_type"])
	}

}

func (s iufService) processOutputOfProcessMedia(activity *iuf.Activity, workflow *v1alpha1.Workflow) error {
	nodesWithOutputs := workflow.Status.Nodes.Filter(func(nodeStatus v1alpha1.NodeStatus) bool {
		return nodeStatus.Outputs.HasOutputs() && len(nodeStatus.Outputs.Parameters) == 2
	})
	if len(nodesWithOutputs) == 0 {
		return nil
	}
	activity.OperationOutputs = map[string]interface{}{
		"stage_params": map[string]interface{}{
			"process-media": map[string]interface{}{
				"products": map[string]interface{}{},
			},
		},
	}
	activity.Products = []iuf.Product{}
	for _, nodeStatus := range nodesWithOutputs {
		var manifest map[string]interface{}
		err := yaml.Unmarshal([]byte(nodeStatus.Outputs.Parameters[0].Value.String()), &manifest)
		if err != nil {
			s.logger.Error(err)
			return err
		}
		// validate iuf product manifest
		data, _ := yaml.Marshal(manifest)
		validated := true
		err = iuf.Validate(data)
		if err != nil {
			s.logger.Error(err)
			validated = false
		}
		jsonManifest, _ := json.Marshal(manifest)
		if manifest["name"] != nil && manifest["version"] != nil {
			s.logger.Infof("manifest: %s - %s", manifest["name"], manifest["version"])
			// add product to activity object
			activity.Products = append(activity.Products, iuf.Product{
				Name:             fmt.Sprintf("%v", manifest["name"]),
				Version:          fmt.Sprintf("%v", manifest["version"]),
				Validated:        validated,
				Manifest:         string(jsonManifest),
				OriginalLocation: nodeStatus.Outputs.Parameters[1].Value.String(),
			})
			productKey := s.getProductVersionKeyFromNameAndVersion(manifest["name"].(string), manifest["version"].(string))

			activity.OperationOutputs["stage_params"].(map[string]interface{})["process-media"].(map[string]interface{})["products"].(map[string]interface{})[fmt.Sprintf("%v", productKey)] = make(map[string]interface{})

			activity.OperationOutputs["stage_params"].(map[string]interface{})["process-media"].(map[string]interface{})["products"].(map[string]interface{})[fmt.Sprintf("%v", productKey)].(map[string]interface{})["parent_directory"] = nodeStatus.Outputs.Parameters[1].Value.String()
		}
	}
	return nil
}

func (s iufService) updateActivityOperationOutputFromWorkflow(
	activity iuf.Activity,
	session iuf.Session,
	nodeStatus *v1alpha1.NodeStatus,
	operationName string,
	productKey string,
) error {
	// no-op if there is no outputs
	if nodeStatus.Outputs == nil {
		return nil
	}
	if activity.OperationOutputs == nil {
		activity.OperationOutputs = make(map[string]interface{})
	}
	if activity.OperationOutputs[session.CurrentStage] == nil {
		activity.OperationOutputs[session.CurrentStage] = make(map[string]interface{})
	}
	outputStage := activity.OperationOutputs[session.CurrentStage].(map[string]interface{})
	if outputStage[operationName] == nil {
		outputStage[operationName] = make(map[string]interface{})
	}
	outputOperation := outputStage[operationName].(map[string]interface{})
	if productKey != "" {
		if outputOperation[productKey] == nil {
			outputOperation[productKey] = make(map[string]interface{})
		}
		operationOutputOfProduct := outputOperation[productKey].(map[string]interface{})
		for _, param := range nodeStatus.Outputs.Parameters {
			operationOutputOfProduct[param.Name] = param.Value
		}

	} else {
		for _, param := range nodeStatus.Outputs.Parameters {
			outputOperation[param.Name] = param.Value
		}
	}
	activity.OperationOutputs[session.CurrentStage] = outputStage
	configmap, err := s.iufObjectToConfigMapData(activity, activity.Name, LABEL_ACTIVITY)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	_, err = s.k8sRestClientSet.
		CoreV1().
		ConfigMaps(DEFAULT_NAMESPACE).
		Update(
			context.TODO(),
			&configmap,
			v1.UpdateOptions{},
		)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}
