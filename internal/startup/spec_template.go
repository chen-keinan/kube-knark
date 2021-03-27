package startup

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/gobuffalo/packr"
)

//GenerateSpecFiles generate spec api files from template
func GenerateSpecFiles() ([]utils.FilesInfo, error) {
	fileInfo := make([]utils.FilesInfo, 0)
	box := packr.NewBox("./../spec/api")
	// Add workload spec api
	ksf, err := box.FindString(common.Workload)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load workload spec api %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.Workload, Data: ksf})
	// Add services spec api
	sb, err := box.FindString(common.Services)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load services spec api %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.Services, Data: sb})
	return fileInfo, nil
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
