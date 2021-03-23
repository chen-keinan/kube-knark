package utils

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/utils/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestCreateKubeKnarkFolders(t *testing.T) {
	err := CreateKubeKnarkFolders()
	assert.NoError(t, err)
	hf, err := GetHomeFolder()
	assert.NoError(t, err)
	_, err = os.Stat(hf)
	assert.NoError(t, err)
	compiledFolder, err := GetEbpfCompiledFolder()
	assert.NoError(t, err)
	_, err = os.Stat(compiledFolder)
	assert.NoError(t, err)
	sourceFolder, err := GetEbpfSourceFolder()
	assert.NoError(t, err)
	_, err = os.Stat(sourceFolder)
	assert.NoError(t, err)
	folder, err := GetSpecAPIFolder()
	assert.NoError(t, err)
	_, err = os.Stat(folder)
	assert.NoError(t, err)
}

func TestGetEnv(t *testing.T) {
	p := GetEnv("a", "p")
	assert.Equal(t, p, "p")
	os.Setenv("a", "k")
	r := GetEnv("a", "p")
	assert.Equal(t, r, "k")
}

func TestCreateEbpfCompiledFolderIfNotExist(t *testing.T) {
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	folder, err := GetEbpfCompiledFolder()
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(folder).Return(fmt.Errorf("failed to create folder")).Times(1)
	err = CreateEbpfCompiledFolderIfNotExist(fm)
	assert.Error(t, err)
}
func TestCreateSpecAPIFolderIfNotExistError(t *testing.T) {
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	folder, err := GetSpecAPIFolder()
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(folder).Return(fmt.Errorf("failed to create folder")).Times(1)
	err = CreateSpecAPIFolderIfNotExist(fm)
	assert.Error(t, err)
}
func TestCreateEbpfSourceFolderIfNotExistError(t *testing.T) {
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	folder, err := GetEbpfSourceFolder()
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(folder).Return(fmt.Errorf("failed to create folder")).Times(1)
	err = CreateEbpfSourceFolderIfNotExist(fm)
	assert.Error(t, err)
}
func TestCreateEbpfSourceFolderIfNotExist(t *testing.T) {
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	folder, err := GetEbpfSourceFolder()
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(folder).Return(nil).Times(1)
	err = CreateEbpfSourceFolderIfNotExist(fm)
	assert.NoError(t, err)
}
func TestCreateSpecAPIFolderIfNotExist(t *testing.T) {
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	folder, err := GetSpecAPIFolder()
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(folder).Return(nil).Times(1)
	err = CreateSpecAPIFolderIfNotExist(fm)
	assert.NoError(t, err)
}
func TestKFolder_CreateFolder(t *testing.T) {
	folder, err := GetHomeFolder()
	if err != nil {
		assert.NoError(t, err)
	}
	err = NewKFolder().CreateFolder(path.Join(folder, "a"))
	assert.NoError(t, err)
}
func TestGetEbpfCompiledFolder(t *testing.T) {
	folder, err := GetEbpfCompiledFolder()
	assert.NoError(t, err)
	homeFolder, err := GetHomeFolder()
	assert.NoError(t, err)
	assert.Equal(t, folder, path.Join(homeFolder, CompileSubFolder))
}
func TestGetSpecAPIFolder(t *testing.T) {
	folder, err := GetSpecAPIFolder()
	assert.NoError(t, err)
	homeFolder, err := GetHomeFolder()
	assert.NoError(t, err)
	assert.Equal(t, folder, path.Join(homeFolder, SpecSubFolder))
}
func TestGetEbpfSourceFolder(t *testing.T) {
	folder, err := GetEbpfSourceFolder()
	assert.NoError(t, err)
	homeFolder, err := GetHomeFolder()
	assert.NoError(t, err)
	assert.Equal(t, folder, path.Join(homeFolder, SourceSubFolder))
}
