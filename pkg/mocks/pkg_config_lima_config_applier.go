// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/runfinch/finch/pkg/config (interfaces: LimaConfigApplier)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// LimaConfigApplier is a mock of LimaConfigApplier interface.
type LimaConfigApplier struct {
	ctrl     *gomock.Controller
	recorder *LimaConfigApplierMockRecorder
}

// LimaConfigApplierMockRecorder is the mock recorder for LimaConfigApplier.
type LimaConfigApplierMockRecorder struct {
	mock *LimaConfigApplier
}

// NewLimaConfigApplier creates a new mock instance.
func NewLimaConfigApplier(ctrl *gomock.Controller) *LimaConfigApplier {
	mock := &LimaConfigApplier{ctrl: ctrl}
	mock.recorder = &LimaConfigApplierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *LimaConfigApplier) EXPECT() *LimaConfigApplierMockRecorder {
	return m.recorder
}

// Apply mocks base method.
func (m *LimaConfigApplier) Apply(arg0 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Apply", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Apply indicates an expected call of Apply.
func (mr *LimaConfigApplierMockRecorder) Apply(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Apply", reflect.TypeOf((*LimaConfigApplier)(nil).Apply), arg0)
}
