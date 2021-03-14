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
func GetEbpfSourceFolder() string {
	err := CreateHomeFolderIfNotExist()
	if err != nil {
		panic("Failed to fetch user home folder")
	}
	ebpfFolder, err := CreateEbpfSourceFolderIfNotExist()
	if err != nil {
		panic("Failed to fetch user home folder")
	}
	return ebpfFolder
}

//GetEbpfCompiledFolder return ebpf compiled folder path
func GetEbpfCompiledFolder() string {
	err := CreateHomeFolderIfNotExist()
	if err != nil {
		panic("Failed to fetch user home folder")
	}
	ebpfFolder, err := CreateEbpfCompiledFolderIfNotExist()
	if err != nil {
		panic("Failed to fetch user home folder")
	}
	return ebpfFolder
}

//GetHomeFolder return beacon home folder
func GetHomeFolder() string {
	usr, err := user.Current()
	if err != nil {
		panic("Failed to fetch user home folder")
	}
	// User can set a custom KUBE_KNARK_HOME from environment variable
	usrHome := GetEnv(common.KubeKnarkHome, usr.HomeDir)
	return path.Join(usrHome, ".kube-knark")
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
	beaconFolder := GetHomeFolder()
	_, err := os.Stat(beaconFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(beaconFolder, 0750)
		if errDir != nil {
			return fmt.Errorf("failed to create beacon home folder at %s", beaconFolder)
		}
	}
	return nil
}

//CreateEbpfSourceFolderIfNotExist create ebpf source folder if not exist
func CreateEbpfSourceFolderIfNotExist() (string, error) {
	ebpfFolder := filepath.Join(GetHomeFolder(), fmt.Sprintf("ebpf/source"))
	_, err := os.Stat(ebpfFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(ebpfFolder, 0750)
		if errDir != nil {
			return "", fmt.Errorf("failed to create ebpf folder folder at %s", ebpfFolder)
		}
	}
	return ebpfFolder, nil
}

//CreateEbpfCompiledFolderIfNotExist create ebpf compiled folder if not exist
func CreateEbpfCompiledFolderIfNotExist() (string, error) {
	ebpfFolder := filepath.Join(GetHomeFolder(), fmt.Sprintf("ebpf/compiled"))
	_, err := os.Stat(ebpfFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(ebpfFolder, 0750)
		if errDir != nil {
			return "", fmt.Errorf("failed to create ebpf folder folder at %s", ebpfFolder)
		}
	}
	return ebpfFolder, nil
}
