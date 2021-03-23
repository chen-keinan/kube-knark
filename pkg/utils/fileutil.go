package utils

import (
	"github.com/chen-keinan/kube-knark/internal/common"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const SpecSubFolder = "spec/api"
const SourceSubFolder = "ebpf/source"
const CompileSubFolder = "ebpf/compile"

//FolderMgr defines the interface for kube-knark folder
//fileutil.go
//go:generate mockgen -destination=./mocks/mock_FolderMgr.go -package=mocks . FolderMgr
type FolderMgr interface {
	CreateFolder(folderName string) error
}

//KFolder kube-knark folder object
type KFolder struct {
}

//NewKFolder return KFolder instance
func NewKFolder() FolderMgr {
	return &KFolder{}
}

//CreateFolder create new kube knark folder
func (kf KFolder) CreateFolder(folderName string) error {
	_, err := os.Stat(folderName)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folderName, 0750)
		if errDir != nil {
			return err
		}
	}
	return nil
}

//GetSpecAPIFolder return spec files source folder path
func GetSpecAPIFolder() (string, error) {
	folder, err := GetHomeFolder()
	if err != nil {
		return "", err
	}
	return path.Join(folder, SpecSubFolder), nil
}

//GetEbpfSourceFolder return ebpf source folder path
func GetEbpfSourceFolder() (string, error) {
	folder, err := GetHomeFolder()
	if err != nil {
		return "", err
	}
	return path.Join(folder, SourceSubFolder), nil
}

//GetEbpfCompiledFolder return ebpf compiled folder path
func GetEbpfCompiledFolder() (string, error) {
	folder, err := GetHomeFolder()
	if err != nil {
		return "", err
	}
	return path.Join(folder, CompileSubFolder), nil
}

//GetHomeFolder return beacon home folder
func GetHomeFolder() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	// User can set a custom KUBE_KNARK_HOME from environment variable
	usrHome := GetEnv(common.KubeKnarkHome, usr.HomeDir)
	return path.Join(usrHome, ".kube-knark"), nil
}

//GetEnv Get Environment Variable value or return default
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//FilesInfo file data
type FilesInfo struct {
	Name string
	Data string
}

//GetFiles return ebpf source files
func GetFiles(folder string) ([]FilesInfo, error) {
	filesData := make([]FilesInfo, 0)
	filesInfo, err := ioutil.ReadDir(filepath.Join(folder))
	if err != nil {
		return nil, err
	}
	for _, fileInfo := range filesInfo {
		f, err := os.Open(filepath.Join(folder, fileInfo.Name()))
		if err != nil {
			return nil, err
		}
		fData, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		filesData = append(filesData, FilesInfo{Name: fileInfo.Name(), Data: string(fData)})
	}
	return filesData, nil
}

//CreateHomeFolderIfNotExist create ebpf home folder if not exist
func CreateHomeFolderIfNotExist(fm FolderMgr) error {
	knarkFolder, err := GetHomeFolder()
	if err != nil {
		return err
	}
	return fm.CreateFolder(knarkFolder)
}

//CreateEbpfSourceFolderIfNotExist create ebpf source folder if not exist
func CreateEbpfSourceFolderIfNotExist(fm FolderMgr) error {
	ebpfFolder, err := GetEbpfSourceFolder()
	if err != nil {
		return err
	}
	return fm.CreateFolder(ebpfFolder)
}

//CreateSpecAPIFolderIfNotExist create spec api folder if not exist
func CreateSpecAPIFolderIfNotExist(fm FolderMgr) error {
	specFolder, err := GetSpecAPIFolder()
	if err != nil {
		return err
	}
	return fm.CreateFolder(specFolder)
}

//CreateEbpfCompiledFolderIfNotExist create ebpf compiled folder if not exist
func CreateEbpfCompiledFolderIfNotExist(fm FolderMgr) error {
	ebpfFolder, err := GetEbpfCompiledFolder()
	if err != nil {
		return err
	}
	return fm.CreateFolder(ebpfFolder)
}

func CreateKubeKnarkFolders() error {
	fm := NewKFolder()
	err := CreateHomeFolderIfNotExist(fm)
	if err != nil {
		return err
	}
	err = CreateEbpfSourceFolderIfNotExist(fm)
	if err != nil {
		return err
	}
	err = CreateEbpfCompiledFolderIfNotExist(fm)
	if err != nil {
		return err
	}
	err = CreateSpecAPIFolderIfNotExist(fm)
	if err != nil {
		return err
	}
	return nil
}
