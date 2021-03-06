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

// Automatically generated by MockGen. DO NOT EDIT!
// Source: handler/store/instancestore.go

package mocks

import (
	context "context"
	types0 "github.com/blox/blox/cluster-state-service/handler/store/types"
	types "github.com/blox/blox/cluster-state-service/handler/types"
	gomock "github.com/golang/mock/gomock"
)

// Mock of ContainerInstanceStore interface
type MockContainerInstanceStore struct {
	ctrl     *gomock.Controller
	recorder *_MockContainerInstanceStoreRecorder
}

// Recorder for MockContainerInstanceStore (not exported)
type _MockContainerInstanceStoreRecorder struct {
	mock *MockContainerInstanceStore
}

func NewMockContainerInstanceStore(ctrl *gomock.Controller) *MockContainerInstanceStore {
	mock := &MockContainerInstanceStore{ctrl: ctrl}
	mock.recorder = &_MockContainerInstanceStoreRecorder{mock}
	return mock
}

func (_m *MockContainerInstanceStore) EXPECT() *_MockContainerInstanceStoreRecorder {
	return _m.recorder
}

func (_m *MockContainerInstanceStore) AddContainerInstance(instance string) error {
	ret := _m.ctrl.Call(_m, "AddContainerInstance", instance)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockContainerInstanceStoreRecorder) AddContainerInstance(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddContainerInstance", arg0)
}

func (_m *MockContainerInstanceStore) AddUnversionedContainerInstance(instance string) error {
	ret := _m.ctrl.Call(_m, "AddUnversionedContainerInstance", instance)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockContainerInstanceStoreRecorder) AddUnversionedContainerInstance(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddUnversionedContainerInstance", arg0)
}

func (_m *MockContainerInstanceStore) GetContainerInstance(cluster string, instanceARN string) (*types.ContainerInstance, error) {
	ret := _m.ctrl.Call(_m, "GetContainerInstance", cluster, instanceARN)
	ret0, _ := ret[0].(*types.ContainerInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockContainerInstanceStoreRecorder) GetContainerInstance(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetContainerInstance", arg0, arg1)
}

func (_m *MockContainerInstanceStore) ListContainerInstances() ([]types.ContainerInstance, error) {
	ret := _m.ctrl.Call(_m, "ListContainerInstances")
	ret0, _ := ret[0].([]types.ContainerInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockContainerInstanceStoreRecorder) ListContainerInstances() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListContainerInstances")
}

func (_m *MockContainerInstanceStore) FilterContainerInstances(filterMap map[string]string) ([]types.ContainerInstance, error) {
	ret := _m.ctrl.Call(_m, "FilterContainerInstances", filterMap)
	ret0, _ := ret[0].([]types.ContainerInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockContainerInstanceStoreRecorder) FilterContainerInstances(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FilterContainerInstances", arg0)
}

func (_m *MockContainerInstanceStore) StreamContainerInstances(ctx context.Context) (chan types0.ContainerInstanceErrorWrapper, error) {
	ret := _m.ctrl.Call(_m, "StreamContainerInstances", ctx)
	ret0, _ := ret[0].(chan types0.ContainerInstanceErrorWrapper)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockContainerInstanceStoreRecorder) StreamContainerInstances(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "StreamContainerInstances", arg0)
}

func (_m *MockContainerInstanceStore) DeleteContainerInstance(cluster string, instanceARN string) error {
	ret := _m.ctrl.Call(_m, "DeleteContainerInstance", cluster, instanceARN)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockContainerInstanceStoreRecorder) DeleteContainerInstance(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteContainerInstance", arg0, arg1)
}
