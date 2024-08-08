// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build linux

package command_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/runfinch/finch/pkg/command"
	"github.com/runfinch/finch/pkg/mocks"
)

const (
	mockNerdctlConfigPath  = "/etc/finch/nerdctl.toml"
	mockBuildkitSocketPath = "/etc/finch/buildkit"
	mockFinchBinPath       = "/usr/lib/usrexec/finch"
	mockSystemPath         = "/usr/bin"
	finalPath              = mockFinchBinPath + command.EnvKeyPathJoiner + mockSystemPath
)

var mockArgs = []string{"shell", "finch"}

func TestLimaCmdCreator_Create(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		mockSvc func(*mocks.Logger, *mocks.CommandCreator, *mocks.Command, *mocks.NerdctlCmdCreatorSystemDeps)
		wantErr error
	}{
		{
			name:    "happy path",
			wantErr: nil,
			mockSvc: func(logger *mocks.Logger, cmdCreator *mocks.CommandCreator, cmd *mocks.Command, lcd *mocks.NerdctlCmdCreatorSystemDeps) {
				logger.EXPECT().Debugf("Creating nerdctl command: ARGUMENTS: %v", mockArgs)
				cmdCreator.EXPECT().Create("nerdctl", mockArgs).Return(cmd)
				lcd.EXPECT().Env(command.EnvKeyPath).Return(mockSystemPath)
				lcd.EXPECT().Environ().Return([]string{})
				lcd.EXPECT().Stdin().Return(nil)
				lcd.EXPECT().Stdout().Return(nil)
				lcd.EXPECT().Stderr().Return(nil)
				cmd.EXPECT().SetEnv([]string{
					fmt.Sprintf("%s=%s", command.EnvKeyPath, finalPath),
					fmt.Sprintf("%s=%s", command.EnvKeyNerdctlTOML, mockNerdctlConfigPath),
					fmt.Sprintf("%s=%s", command.EnvKeyBuildkitHost, mockBuildkitSocketPath),
				})
				cmd.EXPECT().SetStdin(nil)
				cmd.EXPECT().SetStdout(nil)
				cmd.EXPECT().SetStderr(nil)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			cmdCreator := mocks.NewCommandCreator(ctrl)
			cmd := mocks.NewCommand(ctrl)
			logger := mocks.NewLogger(ctrl)
			lcd := mocks.NewNerdctlCmdCreatorSystemDeps(ctrl)
			tc.mockSvc(logger, cmdCreator, cmd, lcd)
			command.NewNerdctlCmdCreator(
				cmdCreator,
				logger,
				mockNerdctlConfigPath,
				mockBuildkitSocketPath,
				mockFinchBinPath,
				lcd,
			).Create(mockArgs...)
		})
	}
}

// func TestLimaCmdCreator_CreateWithoutStdio(t *testing.T) {
// 	t.Parallel()

// 	testCases := []struct {
// 		name    string
// 		mockSvc func(*mocks.Logger, *mocks.CommandCreator, *mocks.Command, *mocks.NerdctlCmdCreatorSystemDeps)
// 		wantErr error
// 	}{
// 		{
// 			name:    "happy path",
// 			wantErr: nil,
// 			mockSvc: func(logger *mocks.Logger, cmdCreator *mocks.CommandCreator, cmd *mocks.Command, lcd *mocks.NerdctlCmdCreatorSystemDeps) {
// 				logger.EXPECT().Debugf("Creating limactl command: ARGUMENTS: %v, %s: %s", mockArgs, envKeyLimaHome, mockLimaHomePath)
// 				cmdCreator.EXPECT().Create(mockLimactlPath, mockArgs).Return(cmd)
// 				lcd.EXPECT().Environ().Return([]string{})
// 				lcd.EXPECT().Env(envKeyPath).Return(mockSystemPath)
// 				cmd.EXPECT().SetEnv([]string{
// 					fmt.Sprintf("%s=%s", envKeyLimaHome, mockLimaHomePath),
// 					fmt.Sprintf("%s=%s", envKeyPath, finalPath),
// 				})
// 				cmd.EXPECT().SetStdin(nil)
// 				cmd.EXPECT().SetStdout(nil)
// 				cmd.EXPECT().SetStderr(nil)
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()

// 			ctrl := gomock.NewController(t)
// 			cmdCreator := mocks.NewCommandCreator(ctrl)
// 			cmd := mocks.NewCommand(ctrl)
// 			logger := mocks.NewLogger(ctrl)
// 			lcd := mocks.NewNerdctlCmdCreatorSystemDeps(ctrl)
// 			tc.mockSvc(logger, cmdCreator, cmd, lcd)
// 			command.NewNerdctlCmdCreator(cmdCreator, logger, mockLimaHomePath, mockLimactlPath, mockQemuBinPath, lcd).
// 				CreateWithoutStdio(mockArgs...)
// 		})
// 	}
// }

// func TestLimaCmdCreator_RunWithReplacingStdout(t *testing.T) {
// 	t.Parallel()

// 	testCases := []struct {
// 		name     string
// 		mockSvc  func(*mocks.Logger, *mocks.CommandCreator, *mocks.NerdctlCmdCreatorSystemDeps, *gomock.Controller, string, *os.File)
// 		wantErr  error
// 		stdoutRs []command.Replacement
// 		inOut    string
// 		outOut   string
// 	}{
// 		{
// 			name:     "happy path",
// 			wantErr:  nil,
// 			stdoutRs: []command.Replacement{{Source: "s1", Target: "t1"}, {Source: "s3", Target: "t3"}, {Source: "s6", Target: "t6"}},
// 			inOut:    "s1 s2 ,s3 /s4 s1.s5",
// 			outOut:   "t1 s2 ,t3 /s4 t1.s5",
// 			mockSvc: func(logger *mocks.Logger, cmdCreator *mocks.CommandCreator,
// 				lcd *mocks.NerdctlCmdCreatorSystemDeps, ctrl *gomock.Controller, inOut string, f *os.File,
// 			) {
// 				logger.EXPECT().Debugf("Creating limactl command: ARGUMENTS: %v, %s: %s", mockArgs, envKeyLimaHome, mockLimaHomePath)
// 				cmd := mocks.NewCommand(ctrl)
// 				cmdCreator.EXPECT().Create(mockLimactlPath, mockArgs).Return(cmd)
// 				lcd.EXPECT().Environ().Return([]string{})
// 				lcd.EXPECT().Stdin().Return(nil)
// 				lcd.EXPECT().Stderr().Return(nil)
// 				lcd.EXPECT().Env(envKeyPath).Return(mockSystemPath)
// 				cmd.EXPECT().SetEnv([]string{
// 					fmt.Sprintf("%s=%s", envKeyLimaHome, mockLimaHomePath),
// 					fmt.Sprintf("%s=%s", envKeyPath, finalPath),
// 				})
// 				cmd.EXPECT().SetStdin(nil)
// 				var stdoutBuf *bytes.Buffer
// 				cmd.EXPECT().SetStdout(gomock.Any()).Do(func(buf *bytes.Buffer) {
// 					stdoutBuf = buf
// 				})
// 				cmd.EXPECT().SetStderr(nil)
// 				cmd.EXPECT().Run().Do(func() {
// 					stdoutBuf.Write([]byte(inOut))
// 				})
// 				lcd.EXPECT().Stdout().Return(f)
// 			},
// 		},
// 		{
// 			name:     "overlapped replacements",
// 			wantErr:  nil,
// 			stdoutRs: []command.Replacement{{Source: "s1", Target: "s2"}, {Source: "s2", Target: "s3"}},
// 			inOut:    "s1 s2 ,s3 /s4 s1.s5",
// 			outOut:   "s3 s3 ,s3 /s4 s3.s5",
// 			mockSvc: func(logger *mocks.Logger, cmdCreator *mocks.CommandCreator,
// 				lcd *mocks.NerdctlCmdCreatorSystemDeps, ctrl *gomock.Controller, inOut string, f *os.File,
// 			) {
// 				logger.EXPECT().Debugf("Creating limactl command: ARGUMENTS: %v, %s: %s",
// 					mockArgs, envKeyLimaHome, mockLimaHomePath)
// 				cmd := mocks.NewCommand(ctrl)
// 				cmdCreator.EXPECT().Create(mockLimactlPath, mockArgs).Return(cmd)
// 				lcd.EXPECT().Environ().Return([]string{})
// 				lcd.EXPECT().Stdin().Return(nil)
// 				lcd.EXPECT().Stderr().Return(nil)
// 				lcd.EXPECT().Env(envKeyPath).Return(mockSystemPath)
// 				cmd.EXPECT().SetEnv([]string{
// 					fmt.Sprintf("%s=%s", envKeyLimaHome, mockLimaHomePath),
// 					fmt.Sprintf("%s=%s", envKeyPath, finalPath),
// 				})
// 				cmd.EXPECT().SetStdin(nil)
// 				var stdoutBuf *bytes.Buffer
// 				cmd.EXPECT().SetStdout(gomock.Any()).Do(func(buf *bytes.Buffer) {
// 					stdoutBuf = buf
// 				})
// 				cmd.EXPECT().SetStderr(nil)
// 				cmd.EXPECT().Run().Do(func() {
// 					stdoutBuf.Write([]byte(inOut))
// 				})
// 				lcd.EXPECT().Stdout().Return(f)
// 			},
// 		},
// 		{
// 			name:     "empty replacements",
// 			wantErr:  nil,
// 			stdoutRs: []command.Replacement{},
// 			inOut:    "s1 s2 ,s3 /s4 .s5",
// 			outOut:   "s1 s2 ,s3 /s4 .s5",
// 			mockSvc: func(logger *mocks.Logger, cmdCreator *mocks.CommandCreator,
// 				lcd *mocks.NerdctlCmdCreatorSystemDeps, ctrl *gomock.Controller, inOut string, f *os.File,
// 			) {
// 				logger.EXPECT().Debugf("Creating limactl command: ARGUMENTS: %v, %s: %s", mockArgs, envKeyLimaHome, mockLimaHomePath)
// 				cmd := mocks.NewCommand(ctrl)
// 				cmdCreator.EXPECT().Create(mockLimactlPath, mockArgs).Return(cmd)
// 				lcd.EXPECT().Environ().Return([]string{})
// 				lcd.EXPECT().Stdin().Return(nil)
// 				lcd.EXPECT().Stderr().Return(nil)
// 				lcd.EXPECT().Env(envKeyPath).Return(mockSystemPath)
// 				cmd.EXPECT().SetEnv([]string{
// 					fmt.Sprintf("%s=%s", envKeyLimaHome, mockLimaHomePath),
// 					fmt.Sprintf("%s=%s", envKeyPath, finalPath),
// 				})
// 				cmd.EXPECT().SetStdin(nil)
// 				var stdoutBuf *bytes.Buffer
// 				cmd.EXPECT().SetStdout(gomock.Any()).Do(func(buf *bytes.Buffer) {
// 					stdoutBuf = buf
// 				})
// 				cmd.EXPECT().SetStderr(nil)
// 				cmd.EXPECT().Run().Do(func() {
// 					stdoutBuf.Write([]byte(inOut))
// 				})
// 				lcd.EXPECT().Stdout().Return(f)
// 			},
// 		},
// 		{
// 			name:     "running cmd returns error",
// 			wantErr:  errors.New("run cmd error"),
// 			stdoutRs: []command.Replacement{{Source: "source-out", Target: "target-out"}},
// 			inOut:    "source-out",
// 			outOut:   "",
// 			mockSvc: func(logger *mocks.Logger, cmdCreator *mocks.CommandCreator,
// 				lcd *mocks.NerdctlCmdCreatorSystemDeps, ctrl *gomock.Controller, _ string, _ *os.File,
// 			) {
// 				logger.EXPECT().Debugf("Creating limactl command: ARGUMENTS: %v, %s: %s", mockArgs, envKeyLimaHome, mockLimaHomePath)
// 				cmd := mocks.NewCommand(ctrl)
// 				cmdCreator.EXPECT().Create(mockLimactlPath, mockArgs).Return(cmd)
// 				lcd.EXPECT().Environ().Return([]string{})
// 				lcd.EXPECT().Stdin().Return(nil)
// 				lcd.EXPECT().Stderr().Return(nil)
// 				lcd.EXPECT().Env(envKeyPath).Return(mockSystemPath)
// 				cmd.EXPECT().SetEnv([]string{
// 					fmt.Sprintf("%s=%s", envKeyLimaHome, mockLimaHomePath),
// 					fmt.Sprintf("%s=%s", envKeyPath, finalPath),
// 				})
// 				cmd.EXPECT().SetStdin(nil)
// 				cmd.EXPECT().SetStdout(gomock.Any())
// 				cmd.EXPECT().SetStderr(nil)
// 				cmd.EXPECT().Run().Return(errors.New("run cmd error"))
// 			},
// 		},
// 		{
// 			name:     "writing to stdout returns error",
// 			wantErr:  fs.ErrInvalid,
// 			stdoutRs: []command.Replacement{{Source: "source-out", Target: "target-out"}},
// 			inOut:    "source-out",
// 			outOut:   "",
// 			mockSvc: func(logger *mocks.Logger, cmdCreator *mocks.CommandCreator,
// 				lcd *mocks.NerdctlCmdCreatorSystemDeps, ctrl *gomock.Controller, inOut string, _ *os.File,
// 			) {
// 				logger.EXPECT().Debugf("Creating limactl command: ARGUMENTS: %v, %s: %s", mockArgs, envKeyLimaHome, mockLimaHomePath)
// 				cmd := mocks.NewCommand(ctrl)
// 				cmdCreator.EXPECT().Create(mockLimactlPath, mockArgs).Return(cmd)
// 				lcd.EXPECT().Environ().Return([]string{})
// 				lcd.EXPECT().Stdin().Return(nil)
// 				lcd.EXPECT().Stderr().Return(nil)
// 				lcd.EXPECT().Env(envKeyPath).Return(mockSystemPath)
// 				cmd.EXPECT().SetEnv([]string{
// 					fmt.Sprintf("%s=%s", envKeyLimaHome, mockLimaHomePath),
// 					fmt.Sprintf("%s=%s", envKeyPath, finalPath),
// 				})
// 				cmd.EXPECT().SetStdin(nil)
// 				var stdoutBuf *bytes.Buffer
// 				cmd.EXPECT().SetStdout(gomock.Any()).Do(func(buf *bytes.Buffer) {
// 					stdoutBuf = buf
// 				})
// 				cmd.EXPECT().SetStderr(nil)
// 				cmd.EXPECT().Run().Do(func() {
// 					stdoutBuf.Write([]byte(inOut))
// 				})
// 				lcd.EXPECT().Stdout().Return(nil)
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()

// 			ctrl := gomock.NewController(t)
// 			cmdCreator := mocks.NewCommandCreator(ctrl)
// 			logger := mocks.NewLogger(ctrl)
// 			lcd := mocks.NewNerdctlCmdCreatorSystemDeps(ctrl)

// 			stdoutFilepath := filepath.Clean(filepath.Join(t.TempDir(), "test"))
// 			stdoutFile, err := os.Create(stdoutFilepath)
// 			require.NoError(t, err)

// 			tc.mockSvc(logger, cmdCreator, lcd, ctrl, tc.inOut, stdoutFile)
// 			assert.Equal(t, tc.wantErr,
// 				command.NewNerdctlCmdCreator(
// 					cmdCreator,
// 					logger,
// 					mockLimaHomePath,
// 					mockLimactlPath,
// 					mockQemuBinPath,
// 					lcd,
// 				).RunWithReplacingStdout(tc.stdoutRs, mockArgs...))

// 			stdout, err := os.ReadFile(stdoutFilepath)
// 			require.NoError(t, err)
// 			assert.Equal(t, tc.outOut, string(stdout))
// 			assert.NoError(t, stdoutFile.Close())
// 			assert.NoError(t, os.Remove(stdoutFilepath))
// 		})
// 	}
// }
