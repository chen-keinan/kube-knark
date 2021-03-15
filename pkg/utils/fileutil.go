package utils

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

//GetEbpfSourceFolder return ebpf source folder path
func GetEbpfSourceFolder() (string, error) {
	err := CreateHomeFolderIfNotExist()
	if err != nil {
		return "", err
	}
	ebpfFolder, err := CreateEbpfSourceFolderIfNotExist()
	if err != nil {
		return "", err
	}
	return ebpfFolder, nil
}

//GetEbpfCompiledFolder return ebpf compiled folder path
func GetEbpfCompiledFolder() (string, error) {
	err := CreateHomeFolderIfNotExist()
	if err != nil {
		return "", err
	}
	ebpfFolder, err := CreateEbpfCompiledFolderIfNotExist()
	if err != nil {
		if err != nil {
			return "", err
		}
	}
	return ebpfFolder, nil
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

//GetEbpfFiles return ebpf source files
func GetEbpfFiles(folder string) ([]FilesInfo, error) {
	filesData := make([]FilesInfo, 0)
	filesInfo, err := ioutil.ReadDir(filepath.Join(folder))
	if err != nil {
		return nil, err
	}
	for _, fileInfo := range filesInfo {
		if err != nil {
			return nil, err
		}
		filesData = append(filesData, FilesInfo{Name: fileInfo.Name()})
	}
	return filesData, nil
}

//CreateHomeFolderIfNotExist create ebpf home folder if not exist
func CreateHomeFolderIfNotExist() error {
	knarkFolder, err := GetHomeFolder()
	if err != nil {
		return err
	}
	_, err = os.Stat(knarkFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(knarkFolder, 0750)
		if errDir != nil {
			return fmt.Errorf("failed to create beacon home folder at %s", knarkFolder)
		}
	}
	return nil
}

//CreateEbpfSourceFolderIfNotExist create ebpf source folder if not exist
func CreateEbpfSourceFolderIfNotExist() (string, error) {
	homeFolder, err := GetHomeFolder()
	if err != nil {
		return "", err
	}
	ebpfFolder := filepath.Join(homeFolder, "ebpf/source")
	_, err = os.Stat(ebpfFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(ebpfFolder, 0750)
		if errDir != nil {
			return "", fmt.Errorf("failed to create ebpf homeFolder homeFolder at %s", ebpfFolder)
		}
	}
	return ebpfFolder, nil
}

//CreateEbpfCompiledFolderIfNotExist create ebpf compiled folder if not exist
func CreateEbpfCompiledFolderIfNotExist() (string, error) {
	homeFolder, err := GetHomeFolder()
	if err != nil {
		return "", err
	}
	ebpfFolder := filepath.Join(homeFolder, "ebpf/compiled")
	_, err = os.Stat(ebpfFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(ebpfFolder, 0750)
		if errDir != nil {
			return "", fmt.Errorf("failed to create ebpf folder folder at %s", ebpfFolder)
		}
	}
	return ebpfFolder, nil
}
