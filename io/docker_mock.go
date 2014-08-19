// Automatically generated by MockGen. DO NOT EDIT!
// Source: io/docker.go

package io

import (
	bytes "bytes"
	gomock "code.google.com/p/gomock/gomock"
	time "time"
)

// Mock of DockerCli interface
type MockDockerCli struct {
	ctrl     *gomock.Controller
	recorder *_MockDockerCliRecorder
}

// Recorder for MockDockerCli (not exported)
type _MockDockerCliRecorder struct {
	mock *MockDockerCli
}

func NewMockDockerCli(ctrl *gomock.Controller) *MockDockerCli {
	mock := &MockDockerCli{ctrl: ctrl}
	mock.recorder = &_MockDockerCliRecorder{mock}
	return mock
}

func (_m *MockDockerCli) EXPECT() *_MockDockerCliRecorder {
	return _m.recorder
}

func (_m *MockDockerCli) CmdRun(_param0 *RunConfig, _param1 ...string) (*bytes.Buffer, string, error) {
	_s := []interface{}{_param0}
	for _, _x := range _param1 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "CmdRun", _s...)
	ret0, _ := ret[0].(*bytes.Buffer)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockDockerCliRecorder) CmdRun(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0}, arg1...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CmdRun", _s...)
}

func (_m *MockDockerCli) CmdTag(_param0 string, _param1 bool, _param2 *TagInfo) error {
	ret := _m.ctrl.Call(_m, "CmdTag", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockDockerCliRecorder) CmdTag(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CmdTag", arg0, arg1, arg2)
}

func (_m *MockDockerCli) CmdCommit(_param0 string, _param1 *TagInfo) (string, error) {
	ret := _m.ctrl.Call(_m, "CmdCommit", _param0, _param1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDockerCliRecorder) CmdCommit(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CmdCommit", arg0, arg1)
}

func (_m *MockDockerCli) CmdBuild(_param0 *BuildConfig, _param1 string, _param2 string) error {
	ret := _m.ctrl.Call(_m, "CmdBuild", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockDockerCliRecorder) CmdBuild(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CmdBuild", arg0, arg1, arg2)
}

func (_m *MockDockerCli) CmdCopy(_param0 map[string]string, _param1 string, _param2 string, _param3 []*CopyArtifact, _param4 string) error {
	ret := _m.ctrl.Call(_m, "CmdCopy", _param0, _param1, _param2, _param3, _param4)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockDockerCliRecorder) CmdCopy(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CmdCopy", arg0, arg1, arg2, arg3, arg4)
}

func (_m *MockDockerCli) CmdLastModTime(_param0 map[string]string, _param1 string, _param2 []*CopyArtifact) (time.Time, error) {
	ret := _m.ctrl.Call(_m, "CmdLastModTime", _param0, _param1, _param2)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDockerCliRecorder) CmdLastModTime(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CmdLastModTime", arg0, arg1, arg2)
}

func (_m *MockDockerCli) CmdStop(_param0 string) error {
	ret := _m.ctrl.Call(_m, "CmdStop", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockDockerCliRecorder) CmdStop(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CmdStop", arg0)
}

func (_m *MockDockerCli) CmdRmContainer(_param0 string) error {
	ret := _m.ctrl.Call(_m, "CmdRmContainer", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockDockerCliRecorder) CmdRmContainer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CmdRmContainer", arg0)
}

func (_m *MockDockerCli) CmdRmImage(_param0 string) error {
	ret := _m.ctrl.Call(_m, "CmdRmImage", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockDockerCliRecorder) CmdRmImage(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CmdRmImage", arg0)
}

func (_m *MockDockerCli) InspectImage(_param0 string) (InspectedImage, error) {
	ret := _m.ctrl.Call(_m, "InspectImage", _param0)
	ret0, _ := ret[0].(InspectedImage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDockerCliRecorder) InspectImage(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "InspectImage", arg0)
}

func (_m *MockDockerCli) InspectContainer(_param0 string) (InspectedContainer, error) {
	ret := _m.ctrl.Call(_m, "InspectContainer", _param0)
	ret0, _ := ret[0].(InspectedContainer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDockerCliRecorder) InspectContainer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "InspectContainer", arg0)
}

// Mock of InspectedImage interface
type MockInspectedImage struct {
	ctrl     *gomock.Controller
	recorder *_MockInspectedImageRecorder
}

// Recorder for MockInspectedImage (not exported)
type _MockInspectedImageRecorder struct {
	mock *MockInspectedImage
}

func NewMockInspectedImage(ctrl *gomock.Controller) *MockInspectedImage {
	mock := &MockInspectedImage{ctrl: ctrl}
	mock.recorder = &_MockInspectedImageRecorder{mock}
	return mock
}

func (_m *MockInspectedImage) EXPECT() *_MockInspectedImageRecorder {
	return _m.recorder
}

func (_m *MockInspectedImage) CreatedTime() time.Time {
	ret := _m.ctrl.Call(_m, "CreatedTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

func (_mr *_MockInspectedImageRecorder) CreatedTime() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreatedTime")
}

func (_m *MockInspectedImage) ID() string {
	ret := _m.ctrl.Call(_m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockInspectedImageRecorder) ID() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ID")
}

func (_m *MockInspectedImage) ContainerID() string {
	ret := _m.ctrl.Call(_m, "ContainerID")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockInspectedImageRecorder) ContainerID() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ContainerID")
}

// Mock of InspectedContainer interface
type MockInspectedContainer struct {
	ctrl     *gomock.Controller
	recorder *_MockInspectedContainerRecorder
}

// Recorder for MockInspectedContainer (not exported)
type _MockInspectedContainerRecorder struct {
	mock *MockInspectedContainer
}

func NewMockInspectedContainer(ctrl *gomock.Controller) *MockInspectedContainer {
	mock := &MockInspectedContainer{ctrl: ctrl}
	mock.recorder = &_MockInspectedContainerRecorder{mock}
	return mock
}

func (_m *MockInspectedContainer) EXPECT() *_MockInspectedContainerRecorder {
	return _m.recorder
}

func (_m *MockInspectedContainer) Running() bool {
	ret := _m.ctrl.Call(_m, "Running")
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockInspectedContainerRecorder) Running() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Running")
}

func (_m *MockInspectedContainer) CreatedTime() time.Time {
	ret := _m.ctrl.Call(_m, "CreatedTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

func (_mr *_MockInspectedContainerRecorder) CreatedTime() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreatedTime")
}

func (_m *MockInspectedContainer) ContainerName() string {
	ret := _m.ctrl.Call(_m, "ContainerName")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockInspectedContainerRecorder) ContainerName() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ContainerName")
}

func (_m *MockInspectedContainer) ExitStatus() int {
	ret := _m.ctrl.Call(_m, "ExitStatus")
	ret0, _ := ret[0].(int)
	return ret0
}

func (_mr *_MockInspectedContainerRecorder) ExitStatus() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ExitStatus")
}
