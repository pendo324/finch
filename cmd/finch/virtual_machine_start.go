// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build darwin || windows

package main

import (
	"fmt"

	"github.com/runfinch/finch/pkg/disk"

	"github.com/runfinch/finch/pkg/command"
	"github.com/runfinch/finch/pkg/config"
	"github.com/runfinch/finch/pkg/dependency"
	"github.com/runfinch/finch/pkg/flog"
	"github.com/runfinch/finch/pkg/lima"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newStartVMCommand(
	ncc command.NerdctlCmdCreator,
	logger flog.Logger,
	optionalDepGroups []*dependency.Group,
	lca config.LimaConfigApplier,
	nca config.NerdctlConfigApplier,
	fs afero.Fs,
	privateKeyPath string,
	dm disk.UserDataDiskManager,
) *cobra.Command {
	return &cobra.Command{
		Use:      "start",
		Short:    "Start the virtual machine",
		RunE:     newStartVMAction(ncc, logger, optionalDepGroups, lca, dm).runAdapter,
		PostRunE: newPostVMStartInitAction(logger, ncc, fs, privateKeyPath, nca).runAdapter,
	}
}

type startVMAction struct {
	creator             command.NerdctlCmdCreator
	logger              flog.Logger
	optionalDepGroups   []*dependency.Group
	limaConfigApplier   config.LimaConfigApplier
	userDataDiskManager disk.UserDataDiskManager
}

func newStartVMAction(
	creator command.NerdctlCmdCreator,
	logger flog.Logger,
	optionalDepGroups []*dependency.Group,
	lca config.LimaConfigApplier,
	dm disk.UserDataDiskManager,
) *startVMAction {
	return &startVMAction{
		creator:             creator,
		logger:              logger,
		optionalDepGroups:   optionalDepGroups,
		limaConfigApplier:   lca,
		userDataDiskManager: dm,
	}
}

func (sva *startVMAction) runAdapter(_ *cobra.Command, _ []string) error {
	return sva.run()
}

func (sva *startVMAction) run() error {
	err := sva.assertVMIsStopped(sva.creator, sva.logger)
	if err != nil {
		return err
	}
	err = dependency.InstallOptionalDeps(sva.optionalDepGroups, sva.logger)
	if err != nil {
		sva.logger.Errorf("Dependency error: %v", err)
	}

	err = sva.limaConfigApplier.ConfigureOverrideLimaYaml()
	if err != nil {
		return err
	}

	// TODO: don't run this on Windows
	err = sva.userDataDiskManager.EnsureUserDataDisk()
	if err != nil {
		return err
	}

	limaCmd := sva.creator.CreateWithoutStdio("start", limaInstanceName)
	sva.logger.Info("Starting existing Finch virtual machine...")
	logs, err := limaCmd.CombinedOutput()
	if err != nil {
		sva.logger.SetFormatter(flog.TextWithoutTruncation)
		sva.logger.Errorf("Finch virtual machine failed to start, debug logs:\n%s", logs)
		sva.logger.SetFormatter(flog.Text)
		return err
	}
	sva.logger.Info("Finch virtual machine started successfully")
	return nil
}

func (sva *startVMAction) assertVMIsStopped(creator command.NerdctlCmdCreator, logger flog.Logger) error {
	status, err := lima.GetVMStatus(creator, logger, limaInstanceName)
	if err != nil {
		return err
	}
	switch status {
	case lima.Nonexistent:
		return fmt.Errorf("the instance %q does not exist, run `finch %s init` to create a new instance",
			limaInstanceName, virtualMachineRootCmd)
	case lima.Running:
		return fmt.Errorf("the instance %q is already running", limaInstanceName)
	default:
		return nil
	}
}
