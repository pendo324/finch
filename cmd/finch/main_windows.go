// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build windows

//go:generate go-winres make --file-version=git-tag --product-version=git-tag --arch amd64 --in ../../winres/winres.json

package main

import (
	"github.com/spf13/afero"

	"github.com/runfinch/finch/pkg/command"
	"github.com/runfinch/finch/pkg/config"
	"github.com/runfinch/finch/pkg/dependency"
	"github.com/runfinch/finch/pkg/dependency/credhelper"
	"github.com/runfinch/finch/pkg/disk"
	"github.com/runfinch/finch/pkg/flog"
	"github.com/runfinch/finch/pkg/path"
	"github.com/runfinch/finch/pkg/system"
	"github.com/runfinch/finch/pkg/winutil"
)

func dependencies(
	ecc *command.ExecCmdCreator,
	fc *config.Finch,
	fp path.Finch,
	fs afero.Fs,
	_ command.LimaCmdCreator,
	logger flog.Logger,
	finchDir string,
) []*dependency.Group {
	return []*dependency.Group{
		credhelper.NewDependencyGroup(
			ecc,
			fs,
			fp,
			logger,
			fc,
			finchDir,
			system.NewStdLib().Arch(),
		),
	}
}

func dataDiskManager(
	lcc command.LimaCmdCreator,
	ecc *command.ExecCmdCreator,
	fp path.Finch,
	finchRootPath string,
	fc *config.Finch,
	logger flog.Logger,
) disk.UserDataDiskManager {
	return disk.NewUserDataDiskManager(
		lcc,
		ecc,
		&afero.OsFs{},
		fp,
		finchRootPath,
		fc,
		logger,
		winutil.NewElevatedCommand(),
	)
}
