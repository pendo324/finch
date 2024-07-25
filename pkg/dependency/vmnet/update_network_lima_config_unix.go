// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build darwin && !native

package vmnet

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/runfinch/finch/pkg/dependency"
	"github.com/runfinch/finch/pkg/flog"
	"github.com/runfinch/finch/pkg/path"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// defaultLimaConfig updates the lima configuration after other network dependencies are installed.
type defaultLimaConfig struct {
	fp          path.Finch
	binaries    dependency.Dependency
	sudoersFile dependency.Dependency
	fs          afero.Fs
	l           flog.Logger
}

var _ dependency.Dependency = (*defaultLimaConfig)(nil)

func newDefaultLimaConfig(
	fp path.Finch,
	binaries dependency.Dependency,
	sudoersFile dependency.Dependency,
	fs afero.Fs,
	l flog.Logger,
) *defaultLimaConfig {
	return &defaultLimaConfig{
		// TODO: consider replacing fp with only the strings that are used instead of the entire type
		fp:          fp,
		binaries:    binaries,
		sudoersFile: sudoersFile,
		fs:          fs,
		l:           l,
	}
}

// Snippet to append to a lima yaml file to setup a managed network called "finch-shared".
// This must match the value in [networks.yaml].
// TODO: Use limayaml.LimaYAML instead of appending a raw string?
const networkConfigString = `networks:
  - lima: finch-shared
`

// NetworkConfig is a struct for (partially) deserializing lima yaml.
type NetworkConfig struct {
	Networks []map[string]string `yaml:"networks"`
}

// verifyConfigHasNetworkSection deserializes a yaml file at filePath and verifies that it has the expected value.
func (overConf *defaultLimaConfig) verifyConfigHasNetworkSection(filePath string) bool {
	yamlFile, err := afero.ReadFile(overConf.fs, filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			overConf.l.Debugf("config file not found: %v", err)
		} else {
			overConf.l.Errorf("failed to read config file: %v", err)
		}
		return false
	}
	var cfg NetworkConfig
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		overConf.l.Errorf("failed to unmarshal YAML from default config file: %v", err)
		return false
	}

	networksLen := len(cfg.Networks)
	if networksLen > 1 {
		overConf.l.Errorf("default config file has incorrect number of Networks defined (%d)", networksLen)
		return false
	} else if networksLen == 0 {
		return false
	}

	return cfg.Networks[0]["lima"] == "finch-shared"
}

// appendNetworkConfiguration adds a new network config section to a file at filePath.
func (overConf *defaultLimaConfig) appendNetworkConfiguration(filePath string) error {
	f, err := overConf.fs.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("error opening file at path %s, error: %w", filePath, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			overConf.l.Errorf("error closing file at path %s, error: %v", filePath, err)
		}
	}()
	if _, err := f.WriteString(networkConfigString); err != nil {
		return fmt.Errorf("error writing to file at path %s", filePath)
	}

	return nil
}

// shouldAddNetworksConfig returns true iff binaries and sudoers are installed as
// updating the network config without those dependencies leads to a broken user experience.
func (overConf *defaultLimaConfig) shouldAddNetworksConfig() bool {
	return overConf.binaries.Installed() && overConf.sudoersFile.Installed()
}

// Installed returns true iff lima config has been updated.
func (overConf *defaultLimaConfig) Installed() bool {
	return overConf.verifyConfigHasNetworkSection(overConf.fp.LimaDefaultConfigPath())
}

// Install adds the networks config block to lima's default config yaml file.
// Only adds if the shouldAddNetworksConfig() helper function is true.
func (overConf *defaultLimaConfig) Install() error {
	if !overConf.shouldAddNetworksConfig() {
		return fmt.Errorf("skipping installation of network configuration because pre-requisites are missing")
	}
	return overConf.appendNetworkConfiguration(overConf.fp.LimaDefaultConfigPath())
}

func (overConf *defaultLimaConfig) RequiresRoot() bool {
	return false
}
