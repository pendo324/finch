// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build !linux

package container

import (
	"github.com/runfinch/common-tests/tests"
)

func getCGroupMode() tests.CGMode {
	return 0
}
