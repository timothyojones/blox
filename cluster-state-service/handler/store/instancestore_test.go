// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package store

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/blox/blox/cluster-state-service/handler/mocks"
	storetypes "github.com/blox/blox/cluster-state-service/handler/store/types"
	"github.com/blox/blox/cluster-state-service/handler/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	containerInstanceARN1 = "arn:aws:ecs:us-east-1:123456789123:container-instance/4b6d45ea-a4b4-4269-9d04-3af6ddfdc597"
	containerInstanceARN2 = "arn:aws:ecs:us-east-1:123456789123:container-instance/3af93452-d6b7-6759-0923-4f5123cfd025"
	status1               = "active"
	status2               = "inactive"
	version               = int64(1)
)

type instanceStoreMockContext struct {
	mockCtrl      *gomock.Controller
	datastore     *mocks.MockDataStore
	etcdTxStore   *mocks.MockEtcdTXStore
	instance1     types.ContainerInstance
	instance2     types.ContainerInstance
	instanceJSON1 string
	instanceJSON2 string
	instanceKey1  string
	instanceKey2  string
}

func NewContainerInstanceStoreMockContext(t *testing.T) *instanceStoreMockContext {
	context := instanceStoreMockContext{}
	context.mockCtrl = gomock.NewController(t)
	context.datastore = mocks.NewMockDataStore(context.mockCtrl)
	context.etcdTxStore = mocks.NewMockEtcdTXStore(context.mockCtrl)

	context.instance1 = types.ContainerInstance{
		Detail: &types.InstanceDetail{
			ContainerInstanceARN: &containerInstanceARN1,
			ClusterARN:           &clusterARN1,
			Status:               &status1,
			Version:              &version,
		},
	}
	context.instanceJSON1 = marshalInstance(t, context.instance1)
	context.instanceKey1 = instanceKeyPrefix + clusterName1 + "/" + containerInstanceARN1

	context.instance2 = types.ContainerInstance{
		Detail: &types.InstanceDetail{
			ContainerInstanceARN: &containerInstanceARN2,
			ClusterARN:           &clusterARN2,
			Status:               &status2,
			Version:              &version,
		},
	}
	context.instanceJSON2 = marshalInstance(t, context.instance2)
	context.instanceKey2 = instanceKeyPrefix + clusterName2 + "/" + containerInstanceARN2

	return &context
}

func TestInstanceStoreNilDatastore(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()
	_, err := NewContainerInstanceStore(nil, context.etcdTxStore)

	if err == nil {
		t.Error("Expected an error when datastore is nil")
	}
}

func TestInstanceStoreNilEtcdTxStore(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()
	_, err := NewContainerInstanceStore(context.datastore, nil)

	if err == nil {
		t.Error("Expected an error when etcd transactional store is nil")
	}
}

func TestInstanceStore(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	if instanceStore == nil {
		t.Error("Instancestore should not be nil")
	}
}

func TestAddContainerInstanceEmptyInstanceJSON(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	err := instanceStore.AddContainerInstance("")

	if err == nil {
		t.Error("Expected an error when instance JSON is empty in AddContainerInstance")
	}
}

func TestAddContainerInstanceJSONUnmarshalError(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	err := instanceStore.AddContainerInstance("invalidJSON")

	if err == nil {
		t.Error("Expected an error when instance JSON is invalid in AddContainerInstance")
	}
}

func TestAddContainerInstanceContainerInstanceDetailNotSet(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	instance := types.ContainerInstance{}

	instanceJSON, err := json.Marshal(instance)
	err = instanceStore.AddContainerInstance(string(instanceJSON))

	if err == nil {
		t.Error("Expected an error when container instance detail is not set")
	}
}

func TestAddContainerInstanceContainerInstanceARNNotSet(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	instance := types.ContainerInstance{
		Detail: &types.InstanceDetail{
			ClusterARN: &clusterARN1,
		},
	}

	instanceJSON, err := json.Marshal(instance)
	err = instanceStore.AddContainerInstance(string(instanceJSON))

	if err == nil {
		t.Error("Expected an error when container instance ARN is not set")
	}
}

func TestAddContainerInstanceClusterARNNotSet(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	instance := types.ContainerInstance{
		Detail: &types.InstanceDetail{
			ContainerInstanceARN: &containerInstanceARN1,
		},
	}

	instanceJSON, err := json.Marshal(instance)
	err = instanceStore.AddContainerInstance(string(instanceJSON))

	if err == nil {
		t.Error("Expected an error when cluster ARN is not set")
	}
}

func TestAddContainerInstanceEmptyContainerInstanceARN(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	instanceARN := ""
	instance := types.ContainerInstance{
		Detail: &types.InstanceDetail{
			ContainerInstanceARN: &instanceARN,
			ClusterARN:           &clusterARN1,
		},
	}

	instanceJSON, err := json.Marshal(instance)
	err = instanceStore.AddContainerInstance(string(instanceJSON))

	if err == nil {
		t.Error("Expected an error when container instance ARN is an empty string")
	}
}

func TestAddContainerInstanceEmptyClusterARN(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	clusterARN := ""
	instance := types.ContainerInstance{
		Detail: &types.InstanceDetail{
			ContainerInstanceARN: &containerInstanceARN1,
			ClusterARN:           &clusterARN,
		},
	}

	instanceJSON, err := json.Marshal(instance)
	err = instanceStore.AddContainerInstance(string(instanceJSON))

	if err == nil {
		t.Error("Expected an error when cluster ARN is an empty string")
	}
}

func TestAddContainerInstanceSTMRepeatableFails(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	context.etcdTxStore.EXPECT().GetV3Client().Return(nil)
	context.etcdTxStore.EXPECT().NewSTMRepeatable(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("Error when getting key"))

	err := instanceStore.AddContainerInstance(context.instanceJSON1)
	assert.Error(t, err, "Expected error when STM repeatable fails to execute with an error")
}

func TestAddUnversionedContainerInstanceEmptyVersion(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	context.datastore.EXPECT().Get(gomock.Any()).Times(0)
	context.datastore.EXPECT().Add(gomock.Any(), gomock.Any()).Times(0)

	instance := context.instance1
	instance.Detail.Version = nil
	instanceJSON := marshalInstance(t, instance)

	err := instanceStore.AddUnversionedContainerInstance(instanceJSON)

	if err == nil {
		t.Error("Expected an error when adding unversioned container instance with empty version")
	}
}

func TestAddUnversionedContainerInstanceInvalidVersion(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	context.datastore.EXPECT().Get(gomock.Any()).Times(0)
	context.datastore.EXPECT().Add(gomock.Any(), gomock.Any()).Times(0)

	err := instanceStore.AddUnversionedContainerInstance(context.instanceJSON1)

	if err == nil {
		t.Error("Expected an error when adding unversioned container instance with invalid version")
	}
}

func TestGetContainerInstanceEmptyClusterName(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	_, err := instanceStore.GetContainerInstance("", containerInstanceARN1)
	if err == nil {
		t.Error("Expected an error when instance ARN is empty in GetContainerInstance")
	}
}

func TestGetContainerInstanceEmptyInstanceARN(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	_, err := instanceStore.GetContainerInstance(clusterName1, "")
	if err == nil {
		t.Error("Expected an error when instance ARN is empty in GetContainerInstance")
	}
}

func TestGetContainerInstanceGetFails(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	context.datastore.EXPECT().Get(context.instanceKey1).Return(nil, errors.New("Error when getting key"))
	_, err := instanceStore.GetContainerInstance(clusterName1, containerInstanceARN1)
	if err == nil {
		t.Error("Expected an error when datastore get fails")
	}
}

func TestGetContainerInstanceGetNoResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	context.datastore.EXPECT().Get(context.instanceKey1).Return(make(map[string]string), nil)
	instance, err := instanceStore.GetContainerInstance(clusterName1, containerInstanceARN1)
	if err != nil {
		t.Error("Unexpected error when datastore get returns empty results")
	}
	if instance != nil {
		t.Error("Expected GetContainerInstance to return nil when get from datastore is empty")
	}
}

func TestGetContainerInstanceGetMultipleResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	resp := map[string]string{
		containerInstanceARN1: context.instanceJSON1,
		containerInstanceARN2: context.instanceJSON2,
	}
	context.datastore.EXPECT().Get(context.instanceKey1).Return(resp, nil)
	_, err := instanceStore.GetContainerInstance(clusterName1, containerInstanceARN1)
	if err == nil {
		t.Error("Expected an error when datastore get returns multiple results")
	}
}

func TestGetContainerInstanceGetInvalidJSONResult(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	resp := map[string]string{
		containerInstanceARN1: "invalidJSON",
	}
	context.datastore.EXPECT().Get(context.instanceKey1).Return(resp, nil)
	_, err := instanceStore.GetContainerInstance(clusterName1, containerInstanceARN1)
	if err == nil {
		t.Error("Expected an error when datastore get returns invalid JSON results")
	}
}

func TestGetContainerInstanceWithClusterNameAndInstanceARN(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	resp := map[string]string{
		containerInstanceARN1: context.instanceJSON1,
	}
	context.datastore.EXPECT().Get(context.instanceKey1).Return(resp, nil)
	instance, err := instanceStore.GetContainerInstance(clusterName1, containerInstanceARN1)
	assert.NoError(t, err, "Unexpected error when getting an instance")
	if !reflect.DeepEqual(*instance, context.instance1) {
		t.Error("Expected the returned instance to match the one returned from the datastore")
	}
}

func TestGetContainerInstanceWithClusterARNAndInstanceARN(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	resp := map[string]string{
		containerInstanceARN1: context.instanceJSON1,
	}
	context.datastore.EXPECT().Get(context.instanceKey1).Return(resp, nil)
	instance, err := instanceStore.GetContainerInstance(clusterARN1, containerInstanceARN1)
	if err != nil {
		t.Error("Unexpected error when getting an instance")
	}
	if !reflect.DeepEqual(*instance, context.instance1) {
		t.Error("Expected the returned instance to match the one returned from the datastore")
	}
}

func TestListContainerInstancesGetWithPrefixInvalidJson(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	resp := map[string]string{
		containerInstanceARN1: "invalidJSON",
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)
	_, err := instanceStore.ListContainerInstances()
	if err == nil {
		t.Error("Expected an error when datastore GetWithPrefix fails")
	}
}

func TestListContainerInstancesGetWithPrefixFails(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(nil, errors.New("GetWithPrefix failed"))
	_, err := instanceStore.ListContainerInstances()
	if err == nil {
		t.Error("Expected an error when datastore GetWithPrefix fails")
	}
}

func TestListContainerInstancesGetWithPrefixNoResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(make(map[string]string), nil)
	instances, err := instanceStore.ListContainerInstances()
	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns empty")
	}

	if len(instances) > 0 {
		t.Error("Expected ListContainerInstances result to be empty when GetWithPrefix result is empty")
	}
}

func TestListContainerInstancesGetWithPrefixMultipleResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	resp := map[string]string{
		containerInstanceARN1: context.instanceJSON1,
		containerInstanceARN2: context.instanceJSON2,
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)
	instances, err := instanceStore.ListContainerInstances()

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns empty")
	}
	if len(instances) != len(resp) {
		t.Error("Expected ListContainerInstances result to be the same length as the GetWithPrefix result")
	}
	for _, v := range instances {
		value, ok := resp[*v.Detail.ContainerInstanceARN]
		if !ok {
			t.Errorf("Expected GetWithPrefix result to contain the same elements as ListContainerInstances result. Missing %v", v)
		} else {
			instance := unmarshalString(t, value)
			if !reflect.DeepEqual(v, instance) {
				t.Errorf("Expected GetWithPrefix result to contain the same elements as ListContainerInstances result. %v does not match %v", v, instance)
			}
		}
	}
}

func TestFilterContainerInstancesNoFilters(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	var filters map[string]string
	_, err := instanceStore.FilterContainerInstances(filters)
	if err == nil {
		t.Error("Expected an error when filter map is empty FilterContainerInstances")
	}
}

func TestFilterContainerInstancesEmptyValue(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	_, err := instanceStore.FilterContainerInstances(map[string]string{instanceStatusFilter: ""})
	if err == nil {
		t.Error("Expected an error when filterValue is empty in FilterContainerInstances")
	}
}

func TestFilterContainerInstancesUnsupportedFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	filters := map[string]string{"invalidFilter": "value"}
	_, err := instanceStore.FilterContainerInstances(filters)
	if err == nil {
		t.Error("Expected an error when unsupported filter key is provided in FilterContainerInstances")
	}
}

func TestFilterContainerInstancesDatastoreGetWithPrefixFails(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(nil, errors.New("GetWithPrefix failed"))

	filters := map[string]string{instanceStatusFilter: status1}
	instanceStore := instanceStore(t, context)
	_, err := instanceStore.FilterContainerInstances(filters)
	if err == nil {
		t.Error("Expected an error when datastore GetWithPrefix fails in FilterContainerInstances")
	}
}

func TestFilterContainerInstancesDatastoreGetWithPrefixReturnsNoResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(make(map[string]string), nil)

	instanceStore := instanceStore(t, context)
	filters := map[string]string{instanceStatusFilter: status1}
	instances, err := instanceStore.FilterContainerInstances(filters)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns empty map in FilterContainerInstances")
	}

	if instances == nil || len(instances) != 0 {
		t.Error("Result should be empty when datastore GetWithPrefix returns empty map in FilterContainerInstances")
	}
}

func TestFilterContainerInstancesNoResultsMatchStatusFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	resp := map[string]string{
		containerInstanceARN1: context.instanceJSON1,
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)

	instanceStore := instanceStore(t, context)
	filters := map[string]string{instanceStatusFilter: status2}
	instances, err := instanceStore.FilterContainerInstances(filters)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns results in FilterContainerInstances")
	}

	if instances == nil || len(instances) != 0 {
		t.Error("Result should be empty when status filter does not match any results in FilterContainerInstances")
	}
}

func TestFilterContainerInstancesMultipleResultsOneMatchesStatusFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	resp := map[string]string{
		containerInstanceARN1: context.instanceJSON1,
		containerInstanceARN2: context.instanceJSON2,
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)

	instanceStore := instanceStore(t, context)
	filters := map[string]string{instanceStatusFilter: status1}
	instances, err := instanceStore.FilterContainerInstances(filters)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns results in FilterContainerInstances")
	}

	if instances == nil || len(instances) != 1 {
		t.Error("Result should have 1 instance when 1 instance matches results in FilterContainerInstances")
	}

	if !reflect.DeepEqual(instances[0], context.instance1) {
		t.Error("Expected the returned instance to match the instance with status" + status1)
	}
}

func TestFilterContainerInstancesMultipleResultsMatchStatusFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instance := types.ContainerInstance{
		Detail: &types.InstanceDetail{
			ContainerInstanceARN: &containerInstanceARN2,
			ClusterARN:           &clusterARN2,
			Status:               &status1,
			Version:              &version,
		},
	}
	instanceJSON := marshalInstance(t, instance)

	resp := map[string]string{
		containerInstanceARN1: context.instanceJSON1,
		containerInstanceARN2: instanceJSON,
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)

	instanceStore := instanceStore(t, context)
	filters := map[string]string{instanceStatusFilter: status1}
	instances, err := instanceStore.FilterContainerInstances(filters)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns results in FilterContainerInstances")
	}
	validateFilterContainerInstancesResultsMatchDatastoreResponse(t, instances, resp)
}

func TestFilterContainerInstancesClusterNameFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	resp := map[string]string{
		containerInstanceARN1: context.instanceJSON1,
	}

	instancesForClusterPrefix := instanceKeyPrefix + clusterName1 + "/"
	context.datastore.EXPECT().GetWithPrefix(instancesForClusterPrefix).Return(resp, nil)

	instanceStore := instanceStore(t, context)
	filters := map[string]string{instanceClusterFilter: clusterName1}
	instances, err := instanceStore.FilterContainerInstances(filters)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns results in FilterContainerInstances")
	}
	validateFilterContainerInstancesResultsMatchDatastoreResponse(t, instances, resp)
}

func TestFilterContainerInstancesClusterARNFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	resp := map[string]string{
		containerInstanceARN1: context.instanceJSON1,
	}
	instancesForClusterPrefix := instanceKeyPrefix + clusterName1 + "/"
	context.datastore.EXPECT().GetWithPrefix(instancesForClusterPrefix).Return(resp, nil)

	instanceStore := instanceStore(t, context)
	filters := map[string]string{instanceClusterFilter: clusterARN1}
	instances, err := instanceStore.FilterContainerInstances(filters)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns results in FilterContainerInstances")
	}
	validateFilterContainerInstancesResultsMatchDatastoreResponse(t, instances, resp)
}

func TestFilterContainerInstancesStatusAndClusterARNFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instance := types.ContainerInstance{
		Detail: &types.InstanceDetail{
			ContainerInstanceARN: &containerInstanceARN2,
			ClusterARN:           &clusterARN1,
			Status:               &status2,
			Version:              &version,
		},
	}
	instanceJSON := marshalInstance(t, instance)

	resp := map[string]string{
		containerInstanceARN1: context.instanceJSON1, // clusterARN1, status1
		containerInstanceARN2: instanceJSON,          // clusterARN1, status2
	}
	instancesForClusterPrefix := instanceKeyPrefix + clusterName1 + "/"
	context.datastore.EXPECT().GetWithPrefix(instancesForClusterPrefix).Return(resp, nil)

	instanceStore := instanceStore(t, context)
	filters := map[string]string{instanceStatusFilter: status1, instanceClusterFilter: clusterARN1}
	instances, err := instanceStore.FilterContainerInstances(filters)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns results in FilterContainerInstances")
	}

	if instances == nil || len(instances) != 1 {
		t.Error("Result should have 1 instance when 1 instance matches results in FilterContainerInstances")
	}

	if !reflect.DeepEqual(instances[0], context.instance1) {
		t.Error("Expected the returned instance to match the instance with status" + status1)
	}
}

func validateFilterContainerInstancesResultsMatchDatastoreResponse(t *testing.T, instances []types.ContainerInstance, datastoreResp map[string]string) {
	if instances == nil || len(instances) != len(datastoreResp) {
		t.Error("Number or instances in result should match response from datastore")
	}

	for _, v := range instances {
		value, ok := datastoreResp[*v.Detail.ContainerInstanceARN]
		if !ok {
			t.Errorf("Expected FilterContainerInstances result to contain the same elements as datastore GetWithPrefix result. Missing %v", v)
		} else {
			instance := unmarshalString(t, value)
			if !reflect.DeepEqual(v, instance) {
				t.Errorf("Expected FilterContainerInstances result to contain the same elements as GetWithPrefix result. %v does not match %v", v, instance)
			}
		}
	}
}

func TestStreamContainerInstancesDataStoreStreamReturnsError(t *testing.T) {
	ctx := NewContainerInstanceStoreMockContext(t)
	defer ctx.mockCtrl.Finish()

	tstCtx := context.Background()
	ctx.datastore.EXPECT().StreamWithPrefix(gomock.Any(), instanceKeyPrefix).Return(nil, errors.New("StreamWithPrefix failed"))

	instanceStore := instanceStore(t, ctx)
	instaceRespChan, err := instanceStore.StreamContainerInstances(tstCtx)
	if err == nil {
		t.Error("Expected an error when datastore StreamWithPrefix returns an error")
	}
	if instaceRespChan != nil {
		t.Error("Unexpected instance response channel when there is a datastore channel setup error")
	}
}

func TestStreamContainerInstancesValidJSONInDSChannel(t *testing.T) {
	ctx := NewContainerInstanceStoreMockContext(t)
	defer ctx.mockCtrl.Finish()

	tstCtx := context.Background()
	dsChan := make(chan map[string]string)
	defer close(dsChan)
	ctx.datastore.EXPECT().StreamWithPrefix(gomock.Any(), instanceKeyPrefix).Return(dsChan, nil)

	instanceStore := instanceStore(t, ctx)
	instanceRespChan, err := instanceStore.StreamContainerInstances(tstCtx)
	if err != nil {
		t.Error("Unexpected error when calling stream instances")
	}
	if instanceRespChan == nil {
		t.Error("Expected valid non-nil instanceRespChannel")
	}

	instanceResp := addContainerInstanceToDSChanAndReadFromInstanceRespChan(ctx.instanceJSON1, dsChan, instanceRespChan)
	if instanceResp.Err != nil {
		t.Error("Unexpected error when reading instance from channel")
	}
	if !reflect.DeepEqual(ctx.instance1, instanceResp.ContainerInstance) {
		t.Error("Expected instance in instance response to match that in the stream")
	}
}

func TestStreamContainerInstancesInvalidJSONInDSChannel(t *testing.T) {
	ctx := NewContainerInstanceStoreMockContext(t)
	defer ctx.mockCtrl.Finish()

	tstCtx := context.Background()
	dsChan := make(chan map[string]string)
	defer close(dsChan)
	ctx.datastore.EXPECT().StreamWithPrefix(gomock.Any(), instanceKeyPrefix).Return(dsChan, nil)

	instanceStore := instanceStore(t, ctx)
	instanceRespChan, err := instanceStore.StreamContainerInstances(tstCtx)
	if err != nil {
		t.Error("Unexpected error when calling stream instances")
	}
	if instanceRespChan == nil {
		t.Error("Expected valid non-nil instanceRespChannel")
	}

	invalidJSON := "invalidJSON"
	instanceResp := addContainerInstanceToDSChanAndReadFromInstanceRespChan(invalidJSON, dsChan, instanceRespChan)

	if instanceResp.Err == nil {
		t.Error("Expected an error when dsChannel returns an invalid instance json")
	}
	if !reflect.DeepEqual(types.ContainerInstance{}, instanceResp.ContainerInstance) {
		t.Error("Expected empty instance in response when there is a decode error")
	}

	_, ok := <-instanceRespChan
	if ok {
		t.Error("Expected instance response channel to be closed")
	}
}

func TestStreamContainerInstancesCancelUpstreamContext(t *testing.T) {
	ctx := NewContainerInstanceStoreMockContext(t)
	defer ctx.mockCtrl.Finish()

	tstCtx, cancel := context.WithCancel(context.Background())
	dsChan := make(chan map[string]string)
	defer close(dsChan)
	ctx.datastore.EXPECT().StreamWithPrefix(gomock.Any(), instanceKeyPrefix).Return(dsChan, nil)

	instanceStore := instanceStore(t, ctx)
	instanceRespChan, err := instanceStore.StreamContainerInstances(tstCtx)
	if err != nil {
		t.Error("Unexpected error when calling stream instances")
	}
	if instanceRespChan == nil {
		t.Error("Expected valid non-nil instanceRespChannel")
	}

	cancel()

	_, ok := <-instanceRespChan
	if ok {
		t.Error("Expected instance response channel to be closed")
	}
}

func TestStreamContainerInstancesCloseDownstreamChannel(t *testing.T) {
	ctx := NewContainerInstanceStoreMockContext(t)
	defer ctx.mockCtrl.Finish()

	tstCtx := context.Background()
	dsChan := make(chan map[string]string)
	ctx.datastore.EXPECT().StreamWithPrefix(gomock.Any(), instanceKeyPrefix).Return(dsChan, nil)

	instanceStore := instanceStore(t, ctx)
	instanceRespChan, err := instanceStore.StreamContainerInstances(tstCtx)
	if err != nil {
		t.Error("Unexpected error when calling stream instances")
	}
	if instanceRespChan == nil {
		t.Error("Expected valid non-nil instanceRespChannel")
	}

	close(dsChan)

	_, ok := <-instanceRespChan
	if ok {
		t.Error("Expected instance response channel to be closed")
	}
}

func TestDeleteContainerInstanceEmptyInstanceARN(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	err := instanceStore.DeleteContainerInstance(clusterName1, "")
	if err == nil {
		t.Error("Expected an error when instance ARN is empty in DeleteContainerInstance")
	}
}

func TestDeleteContainerInstanceEmptyClusterName(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	err := instanceStore.DeleteContainerInstance("", containerInstanceARN1)
	if err == nil {
		t.Error("Expected an error when instance ARN is empty in DeleteContainerInstance")
	}
}

func TestDeleteContainerInstanceDeleteFails(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	context.datastore.EXPECT().Delete(context.instanceKey1).Return(int64(0), errors.New("Error when deleting key"))
	err := instanceStore.DeleteContainerInstance(clusterName1, containerInstanceARN1)
	if err == nil {
		t.Error("Expected an error when datastore delete fails")
	}
}

func TestDeleteContainerInstanceDeleteNoError(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	context.datastore.EXPECT().Delete(context.instanceKey1).Return(int64(1), nil)
	err := instanceStore.DeleteContainerInstance(clusterName1, containerInstanceARN1)
	if err != nil {
		t.Errorf("Error deleting container instance from data store: %v", err)
	}
}

func TestDeleteContainerInstanceWithClusterARNAndInstanceARN(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	context.datastore.EXPECT().Delete(context.instanceKey1).Return(int64(1), nil)
	err := instanceStore.DeleteContainerInstance(clusterARN1, containerInstanceARN1)
	if err != nil {
		t.Errorf("Error deleting container instance from data store: %v", err)
	}
}

func instanceStore(t *testing.T, context *instanceStoreMockContext) ContainerInstanceStore {
	instanceStore, err := NewContainerInstanceStore(context.datastore, context.etcdTxStore)
	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}
	return instanceStore
}

func marshalInstance(t *testing.T, instance types.ContainerInstance) string {
	instanceJSON, err := json.Marshal(instance)
	if err != nil {
		t.Error("Failed to marshal instance: ", err)
	}
	return string(instanceJSON)
}

func unmarshalString(t *testing.T, str string) types.ContainerInstance {
	var instance types.ContainerInstance
	err := json.Unmarshal([]byte(str), &instance)
	if err != nil {
		t.Error("Failed to unmarshal string: ", err)
	}
	return instance
}

func addContainerInstanceToDSChanAndReadFromInstanceRespChan(instanceToAdd string, dsChan chan map[string]string, instanceRespChan chan storetypes.ContainerInstanceErrorWrapper) storetypes.ContainerInstanceErrorWrapper {
	var instanceResp storetypes.ContainerInstanceErrorWrapper

	doneChan := make(chan bool)
	defer close(doneChan)
	go func() {
		instanceResp = <-instanceRespChan
		doneChan <- true
	}()

	dsVal := map[string]string{containerInstanceARN1: instanceToAdd}
	dsChan <- dsVal
	<-doneChan

	return instanceResp
}
