package startup

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/gobuffalo/packr"
)

//GenerateSpecFiles generate spec api files from template
func GenerateSpecFiles() ([]utils.FilesInfo, error) {
	// add workload spec file
	filesInfo := make([]utils.FilesInfo, 0)
	workloadFile, err := createFileFromTemplate("./../spec/api", common.Workload)
	if err != nil {
		return nil, err
	}
	filesInfo = append(filesInfo, workloadFile...)
	// add services spec file
	servicesFile, err := createFileFromTemplate("./../spec/api", common.Services)
	if err != nil {
		return nil, err
	}
	filesInfo = append(filesInfo, servicesFile...)
	// Add config storage spec api
	configStorageFile, err := createFileFromTemplate("./../spec/api", common.ConfigAndStorage)
	if err != nil {
		return nil, err
	}
	filesInfo = append(filesInfo, configStorageFile...)
	authenticationFile, err := createFileFromTemplate("./../spec/api", common.Authentication)
	if err != nil {
		return nil, err
	}
	filesInfo = append(filesInfo, authenticationFile...)
	authorizationFile, err := createFileFromTemplate("./../spec/api", common.Authorization)
	if err != nil {
		return nil, err
	}
	filesInfo = append(filesInfo, authorizationFile...)
	policyFile, err := createFileFromTemplate("./../spec/api", common.Policy)
	if err != nil {
		return nil, err
	}
	filesInfo = append(filesInfo, policyFile...)
	extendFile, err := createFileFromTemplate("./../spec/api", common.Extend)
	if err != nil {
		return nil, err
	}
	filesInfo = append(filesInfo, extendFile...)
	clusterFile, err := createFileFromTemplate("./../spec/api", common.Cluster)
	if err != nil {
		return nil, err
	}
	filesInfo = append(filesInfo, clusterFile...)
	return filesInfo, nil
}

//GenerateFileSystemSpec generate spec file system from template
func GenerateFileSystemSpec() ([]utils.FilesInfo, error) {
	fileInfoCom := make([]utils.FilesInfo, 0)
	boxCom := packr.NewBox("./../spec/filesystem")
	// Add workload spec api
	ksfCom, err := boxCom.FindString(common.ConfigFilesPermission)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load filesystem spec api %s", err.Error())
	}
	fileInfoCom = append(fileInfoCom, utils.FilesInfo{Name: common.ConfigFilesPermission, Data: ksfCom})
	return fileInfoCom, nil
}
