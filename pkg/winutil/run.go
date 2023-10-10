// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package winutil

//go:generate mockgen -copyright_file=../../copyright_header -destination=../mocks/winutil_run_windows.go -package=mocks -mock_names ElevatedCommand=ElevatedCommand . ElevatedCommand
type ElevatedCommand interface {
	Run(exePath, wd string, args []string) error
}
