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
	fm := NewKFolder()
	err := CreateKubeKnarkFolders(fm)
	assert.NoError(t, err)
	hf, err := fm.GetHomeFolder()
	assert.NoError(t, err)
	_, err = os.Stat(hf)
	assert.NoError(t, err)
	compiledFolder, err := GetEbpfCompiledFolder(fm)
	assert.NoError(t, err)
	_, err = os.Stat(compiledFolder)
	assert.NoError(t, err)
	sourceFolder, err := GetEbpfSourceFolder(fm)
	assert.NoError(t, err)
	_, err = os.Stat(sourceFolder)
	assert.NoError(t, err)
	folder, err := GetSpecAPIFolder(fm)
	assert.NoError(t, err)
	_, err = os.Stat(folder)
	assert.NoError(t, err)
}

func TestCreateKubeKnarkFoldersErrorHomeFolder(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(1)
	fm.EXPECT().CreateFolder(homePath).Return(fmt.Errorf("error")).Times(1)
	err = CreateKubeKnarkFolders(fm)
	assert.Error(t, err)
}

func TestCreateKubeKnarkFoldersErrorEbpfSourceFolder(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(3)
	fm.EXPECT().CreateFolder(homePath).Return(nil).Times(1)
	sourceEbpfPath, err := GetEbpfSourceFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(sourceEbpfPath).Return(fmt.Errorf("error")).Times(1)
	err = CreateKubeKnarkFolders(fm)
	assert.Error(t, err)
}

func TestCreateKubeKnarkFoldersErrorEbpfCompiledFolder(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(5)
	fm.EXPECT().CreateFolder(homePath).Return(nil).Times(1)
	sourceEbpfPath, err := GetEbpfSourceFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(sourceEbpfPath).Return(nil).Times(1)
	compileEbpfPath, err := GetEbpfCompiledFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(compileEbpfPath).Return(fmt.Errorf("error")).Times(1)
	err = CreateKubeKnarkFolders(fm)
	assert.Error(t, err)
}

func TestCreateKubeKnarkFoldersErrorSpecAPIFolder(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(7)
	fm.EXPECT().CreateFolder(homePath).Return(nil).Times(1)
	sourceEbpfPath, err := GetEbpfSourceFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(sourceEbpfPath).Return(nil).Times(1)
	compileEbpfPath, err := GetEbpfCompiledFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(compileEbpfPath).Return(nil).Times(1)
	apiSpecPath, err := GetSpecAPIFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(apiSpecPath).Return(fmt.Errorf("error")).Times(1)
	err = CreateKubeKnarkFolders(fm)
	assert.Error(t, err)
}

func TestCreateKubeKnarkFoldersErrorSpecFilesystemFolder(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(9)
	fm.EXPECT().CreateFolder(homePath).Return(nil).Times(1)
	sourceEbpfPath, err := GetEbpfSourceFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(sourceEbpfPath).Return(nil).Times(1)
	compileEbpfPath, err := GetEbpfCompiledFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(compileEbpfPath).Return(nil).Times(1)
	apiSpecPath, err := GetSpecAPIFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(apiSpecPath).Return(nil).Times(1)
	fsSpecPath, err := GetSpecFilesystemFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(fsSpecPath).Return(fmt.Errorf("error")).Times(1)
	err = CreateKubeKnarkFolders(fm)
	assert.Error(t, err)
}

func TestCreateKubeKnarkFoldersErrorCompilePluginSubFolder(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(11)
	fm.EXPECT().CreateFolder(homePath).Return(nil).Times(1)
	sourceEbpfPath, err := GetEbpfSourceFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(sourceEbpfPath).Return(nil).Times(1)
	compileEbpfPath, err := GetEbpfCompiledFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(compileEbpfPath).Return(nil).Times(1)
	apiSpecPath, err := GetSpecAPIFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(apiSpecPath).Return(nil).Times(1)
	fsSpecPath, err := GetSpecFilesystemFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(fsSpecPath).Return(nil).Times(1)
	cPluginPath, err := GetCompilePluginSubFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(cPluginPath).Return(fmt.Errorf("error")).Times(1)
	err = CreateKubeKnarkFolders(fm)
	assert.Error(t, err)
}

func TestCreateKubeKnarkFoldersErrorSourcePluginSubFolder(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(13)
	fm.EXPECT().CreateFolder(homePath).Return(nil).Times(1)
	sourceEbpfPath, err := GetEbpfSourceFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(sourceEbpfPath).Return(nil).Times(1)
	compileEbpfPath, err := GetEbpfCompiledFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(compileEbpfPath).Return(nil).Times(1)
	apiSpecPath, err := GetSpecAPIFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(apiSpecPath).Return(nil).Times(1)
	fsSpecPath, err := GetSpecFilesystemFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(fsSpecPath).Return(nil).Times(1)
	cPluginPath, err := GetCompilePluginSubFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(cPluginPath).Return(nil).Times(1)
	sPluginPath, err := GetPluginSourceSubFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(sPluginPath).Return(fmt.Errorf("error")).Times(1)
	err = CreateKubeKnarkFolders(fm)
	assert.Error(t, err)
}

func TestGetEnv(t *testing.T) {
	p := GetEnv("a", "p")
	assert.Equal(t, p, "p")
	os.Setenv("a", "k")
	r := GetEnv("a", "p")
	assert.Equal(t, r, "k")
}

func TestCreateEbpfCompiledFolderIfNotExist(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(2)
	folder, err := GetEbpfCompiledFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(folder).Return(fmt.Errorf("failed to create folder")).Times(1)
	err = CreateEbpfCompiledFolderIfNotExist(fm)
	assert.Error(t, err)
}
func TestCreateSpecAPIFolderIfNotExistError(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(2)
	folder, err := GetSpecAPIFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(folder).Return(fmt.Errorf("failed to create folder")).Times(1)
	err = CreateSpecAPIFolderIfNotExist(fm)
	assert.Error(t, err)
}
func TestCreateSpecFSFolderIfNotExistError(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(2)
	folder, err := GetSpecFilesystemFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(folder).Return(fmt.Errorf("failed to create folder")).Times(1)
	err = CreateSpecFSFolderIfNotExist(fm)
	assert.Error(t, err)
}
func TestCreateEbpfSourceFolderIfNotExistError(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(2)
	folder, err := GetEbpfSourceFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(folder).Return(fmt.Errorf("failed to create folder")).Times(1)
	err = CreateEbpfSourceFolderIfNotExist(fm)
	assert.Error(t, err)
}
func TestCreateEbpfSourceFolderIfNotExist(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(2)
	folder, err := GetEbpfSourceFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(folder).Return(nil).Times(1)
	err = CreateEbpfSourceFolderIfNotExist(fm)
	assert.NoError(t, err)
}
func TestCreateSpecAPIFolderIfNotExist(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(2)
	folder, err := GetSpecAPIFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(folder).Return(nil).Times(1)
	err = CreateSpecAPIFolderIfNotExist(fm)
	assert.NoError(t, err)
}
func TestCreateSpecFSFolderIfNotExist(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(2)
	folder, err := GetSpecFilesystemFolder(fm)
	assert.NoError(t, err)
	fm.EXPECT().CreateFolder(folder).Return(nil).Times(1)
	err = CreateSpecFSFolderIfNotExist(fm)
	assert.NoError(t, err)
}
func TestKFolder_CreateFolder(t *testing.T) {
	fmReal := NewKFolder()
	folder, err := fmReal.GetHomeFolder()
	if err != nil {
		assert.NoError(t, err)
	}
	err = NewKFolder().CreateFolder(path.Join(folder, "a"))
	assert.NoError(t, err)
}
func TestGetEbpfCompiledFolder(t *testing.T) {
	fmReal := NewKFolder()
	folder, err := GetEbpfCompiledFolder(fmReal)
	assert.NoError(t, err)
	homeFolder, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	assert.Equal(t, folder, path.Join(homeFolder, CompileSubFolder))
}
func TestGetSpecAPIFolder(t *testing.T) {
	fmReal := NewKFolder()
	folder, err := GetSpecAPIFolder(fmReal)
	assert.NoError(t, err)
	homeFolder, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	assert.Equal(t, folder, path.Join(homeFolder, SpecAPISubFolder))
}
func TestGetSpecFilesystemFolder(t *testing.T) {
	fmReal := NewKFolder()
	folder, err := GetSpecFilesystemFolder(fmReal)
	assert.NoError(t, err)
	homeFolder, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	assert.Equal(t, folder, path.Join(homeFolder, SpecFileSystemSubFolder))
}

func TestGetEbpfSourceFolder(t *testing.T) {
	fmReal := NewKFolder()
	folder, err := GetEbpfSourceFolder(fmReal)
	assert.NoError(t, err)
	homeFolder, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	assert.Equal(t, folder, path.Join(homeFolder, SourceSubFolder))
}

func TestGetFiles(t *testing.T) {
	fmReal := NewKFolder()
	folder, err := GetEbpfSourceFolder(fmReal)
	assert.NoError(t, err)
	err = os.RemoveAll(folder)
	assert.NoError(t, err)
	fm := NewKFolder()
	err = CreateHomeFolderIfNotExist(fm)
	assert.NoError(t, err)
	err = CreateEbpfSourceFolderIfNotExist(fm)
	assert.NoError(t, err)
	destination, err := os.Create(path.Join(folder, "test.c"))
	assert.NoError(t, err)
	defer destination.Close()
	destination.WriteString("c code")
	files, err := GetFiles(folder)
	assert.NoError(t, err)
	assert.Equal(t, files[0].Name, "test.c")
	assert.Equal(t, files[0].Data, "c code")

}

func TestCreateHomeFolderIfNotExist(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, nil).Times(1)
	fm.EXPECT().CreateFolder(homePath).Return(nil).Times(1)
	err = CreateHomeFolderIfNotExist(fm)
	assert.NoError(t, err)
}
func TestCreateHomeFolderIfNotExistError(t *testing.T) {
	fmReal := NewKFolder()
	ctl := gomock.NewController(t)
	fm := mocks.NewMockFolderMgr(ctl)
	homePath, err := fmReal.GetHomeFolder()
	assert.NoError(t, err)
	fm.EXPECT().GetHomeFolder().Return(homePath, fmt.Errorf("err")).Times(1)
	err = CreateHomeFolderIfNotExist(fm)
	assert.Error(t, err)
}
