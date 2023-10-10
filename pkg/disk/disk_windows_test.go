package disk

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/runfinch/finch/pkg/flog"
	"github.com/runfinch/finch/pkg/mocks"
	fpath "github.com/runfinch/finch/pkg/path"
)

func TestDisk_NewUserDataDiskManager(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	lcc := mocks.NewLimaCmdCreator(ctrl)
	ecc := mocks.NewCommandCreator(ctrl)
	dfs := mocks.NewMockdiskFS(ctrl)
	finch := fpath.Finch("mock_finch")
	homeDir := "mock_home"
	log := mocks.NewLogger(ctrl)

	NewUserDataDiskManager(lcc, ecc, dfs, finch, homeDir, nil, log, nil, nil)
}

func TestUserDataDiskManager_EnsureUserDataDisk(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		finch   fpath.Finch
		rootDir string
		wantErr error
		mockSvc func(
			dfs *mocks.MockdiskFS,
			mFs afero.Fs,
			cmd *mocks.Command,
			ecc *mocks.CommandCreator,
			log *mocks.Logger,
			ec *mocks.ElevatedCommand,
			sd *mocks.UserDataDiskManagerSystemDeps,
		)
	}{
		{
			name:    "first run",
			finch:   fpath.Finch("mock_finch_path"),
			rootDir: "finch_root_dir",
			wantErr: nil,
			mockSvc: func(
				dfs *mocks.MockdiskFS,
				mFs afero.Fs,
				cmd *mocks.Command,
				ecc *mocks.CommandCreator,
				log *mocks.Logger,
				ec *mocks.ElevatedCommand,
				sd *mocks.UserDataDiskManagerSystemDeps,
			) {
				log.EXPECT().Debugf("diskPath: %s", `finch_root_dir\.finch\.disks\a39e889c581af5b4.vhdx`)
				log.EXPECT().Debugf("disksDir: %s", `finch_root_dir\.finch\.disks`)

				dfs.EXPECT().Stat(`finch_root_dir\.finch\.disks`).Return(nil, fs.ErrNotExist)
				dfs.EXPECT().MkdirAll(`finch_root_dir\.finch\.disks`, fs.FileMode(0o700)).Return(nil)

				dfs.EXPECT().Stat(`finch_root_dir\.finch\.disks\a39e889c581af5b4.vhdx`).Return(nil, fs.ErrNotExist)

				log.EXPECT().Infof("creating disk at path: %s", `finch_root_dir\.finch\.disks\a39e889c581af5b4.vhdx`)

				tempFile, err := mFs.Create("tempfile")
				require.NoError(t, err)

				tempFileMatcher := newTempFileNameMatcher(`.*finchCreateDiskOutput.*`)

				dfs.EXPECT().OpenFile(
					tempFileMatcher,
					os.O_RDWR|os.O_CREATE|os.O_EXCL,
					fs.FileMode(0o600),
				).Return(tempFile, nil)

				dpGoOutStr := `{"level": "info", "time": "", "msg": "DiskPart successfully created the virtual disk file.
DiskPart successfully selected the virtual disk file.
DiskPart successfully attached the virtual disk file.
DiskPart succeeded in creating the specified partition.
DiskPart successfully formatted the volume.
DiskPart successfully detached the virtual disk file."}
`
				_, err = tempFile.Write([]byte(dpGoOutStr))
				require.NoError(t, err)

				sd.EXPECT().Executable().Return(`mock_finch_path\finch.exe`, nil)

				ec.EXPECT().Run(
					`"mock_finch_path\bin\dpgo.exe"`,
					`"mock_finch_path"`,
					[]string{
						"--json",
						"--debug",
						"disk",
						"create",
						"--path",
						`"finch_root_dir\.finch\.disks\a39e889c581af5b4.vhdx"`,
						"--size",
						"51200",
						">",
						`"tempfile"`,
						"2>&1",
					},
				).Return(nil)

				dfs.EXPECT().Open("tempfile").Return(tempFile, nil)
				log.EXPECT().Debugf("create disk cmd stdout: %s", dpGoOutStr)
				dfs.EXPECT().Remove("tempfile").Return(nil)

				logs := []flog.Log{{
					Level: "info",
					Time:  "",
					Message: `DiskPart successfully created the virtual disk file.
DiskPart successfully selected the virtual disk file.
DiskPart successfully attached the virtual disk file.
DiskPart succeeded in creating the specified partition.
DiskPart successfully formatted the volume.
DiskPart successfully detached the virtual disk file.`,
				}}
				log.EXPECT().Debugf("create disk cmd stdout parsed: %v", logs)
			},
		},
		// {
		// 	name:    "disk already exists",
		// 	wantErr: nil,
		// 	mockSvc: func(dfs *mocks.MockdiskFS, cmd *mocks.Command, ecc *mocks.CommandCreator, log *mocks.Logger) {
		// 		cmd.EXPECT().Output().Return(listSuccessOutput, nil)
		// 	},
		// },
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ecc := mocks.NewCommandCreator(ctrl)
			mFs := afero.NewMemMapFs()
			dfs := mocks.NewMockdiskFS(ctrl)
			cmd := mocks.NewCommand(ctrl)
			log := mocks.NewLogger(ctrl)
			ec := mocks.NewElevatedCommand(ctrl)
			sd := mocks.NewUserDataDiskManagerSystemDeps(ctrl)
			tc.mockSvc(dfs, mFs, cmd, ecc, log, ec, sd)
			dm := NewUserDataDiskManager(nil, ecc, dfs, tc.finch, tc.rootDir, nil, log, ec, sd)
			err := dm.EnsureUserDataDisk()
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

type tempFileNameMatcher struct {
	pattern *regexp.Regexp
}

func newTempFileNameMatcher(pattern string) *tempFileNameMatcher {
	return &tempFileNameMatcher{
		pattern: regexp.MustCompile(pattern),
	}
}

func (m *tempFileNameMatcher) Matches(x interface{}) bool {
	s, ok := x.(string)
	if !ok {
		return false
	}

	return m.pattern.MatchString(s)
}

func (m *tempFileNameMatcher) String() string {
	return fmt.Sprintf("matches pattern /%v/", m.pattern)
}
