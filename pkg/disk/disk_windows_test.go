package disk

import (
	"io/fs"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

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

	NewUserDataDiskManager(lcc, ecc, dfs, finch, homeDir, nil, log, nil)
}

func TestUserDataDiskManager_EnsureUserDataDisk(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		finch   fpath.Finch
		rootDir string
		wantErr error
		mockSvc func(dfs *mocks.MockdiskFS, cmd *mocks.Command, ecc *mocks.CommandCreator, log *mocks.Logger, ec *mocks.ElevatedCommand)
	}{
		{
			name:    "first run",
			finch:   fpath.Finch("mock_finch_path"),
			rootDir: "finch_root_dir",
			wantErr: nil,
			mockSvc: func(dfs *mocks.MockdiskFS, cmd *mocks.Command, ecc *mocks.CommandCreator, log *mocks.Logger, ec *mocks.ElevatedCommand) {
				log.EXPECT().Debugf("diskPath: %s", `finch_root_dir\.finch\.disks\a39e889c581af5b4.vhdx`)
				log.EXPECT().Debugf("disksDir: %s", `finch_root_dir\.finch\.disks`)

				dfs.EXPECT().Stat(`finch_root_dir\.finch\.disks`).Return(nil, fs.ErrNotExist)
				dfs.EXPECT().MkdirAll(`finch_root_dir\.finch\.disks`, fs.FileMode(0o700)).Return(nil)

				dfs.EXPECT().Stat(`finch_root_dir\.finch\.disks\a39e889c581af5b4.vhdx`).Return(nil, fs.ErrNotExist)

				log.EXPECT().Debugf("creating disk at path: %s", `finch_root_dir\.finch\.disks\a39e889c581af5b4.vhdx`)

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
			dfs := mocks.NewMockdiskFS(ctrl)
			cmd := mocks.NewCommand(ctrl)
			log := mocks.NewLogger(ctrl)
			ec := mocks.NewElevatedCommand(ctrl)
			tc.mockSvc(dfs, cmd, ecc, log, ec)
			dm := NewUserDataDiskManager(nil, ecc, dfs, tc.finch, tc.rootDir, nil, log, ec)
			err := dm.EnsureUserDataDisk()
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
