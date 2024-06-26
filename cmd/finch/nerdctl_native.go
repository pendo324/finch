// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build linux || native

package main

import (
	"github.com/runfinch/finch/pkg/command"
	"github.com/runfinch/finch/pkg/flog"
)

var aliasMap = map[string]string{}

var argHandlerMap = map[string]map[string]argHandler{}

var commandHandlerMap = map[string]commandHandler{}

func (nc *nerdctlCommand) GetCmdArgs() []string {
	return []string{""}
}

func (nc *nerdctlCommand) assertVMIsRunning(creator command.LimaCmdCreator, logger flog.Logger) error {
	return nil
}

func resolveIP(host string, logger flog.Logger, _ command.Creator) (string, error) {
	return "0.0.0.0", nil
}
