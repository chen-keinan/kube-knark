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
	fileInfo := make([]utils.FilesInfo, 0)
	box := packr.NewBox("./../spec/api")
	// Add ebpf kprobe program
	ksf, err := box.FindString(common.Workload)
	if err != nil {
		return nil, fmt.Errorf("faild to load workload spec file %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.Workload, Data: ksf})
	// add services spec file
	kserv, err := box.FindString(common.Services)
	if err != nil {
		return nil, fmt.Errorf("faild to load services spec file %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.Services, Data: kserv})
	cns, err := box.FindString(common.ConfigAndStorage)
	if err != nil {
		return nil, fmt.Errorf("faild to load config and storage spec file %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.ConfigAndStorage, Data: cns})
	// Add config storage spec api
	auth, err := box.FindString(common.Authentication)
	if err != nil {
		return nil, fmt.Errorf("faild to load authentication spec file %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.Authentication, Data: auth})
	authz, err := box.FindString(common.Authorization)
	if err != nil {
		return nil, fmt.Errorf("faild to load authorization spec file %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.Authorization, Data: authz})
	policy, err := box.FindString(common.Policy)
	if err != nil {
		return nil, fmt.Errorf("faild to load policy spec file %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.Policy, Data: policy})
	extend, err := box.FindString(common.Extend)
	if err != nil {
		return nil, fmt.Errorf("faild to load extend spec file %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.Extend, Data: extend})
	cluster, err := box.FindString(common.Cluster)
	if err != nil {
		return nil, fmt.Errorf("faild to load cluster spec file %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.Cluster, Data: cluster})
	return fileInfo, nil
}

//GenerateEbpfFiles ebpf file from template
func GenerateEbpfFiles() ([]utils.FilesInfo, error) {
	filesInfo := make([]utils.FilesInfo, 0)
	box := packr.NewBox("./../../ebpf/")
	// Add ebpf kprobe program
	ksf, err := box.FindString(common.KprobeSourceFile)
	if err != nil {
		return nil, fmt.Errorf("faild to load ebpf source file %s", err.Error())
	}
	filesInfo = append(filesInfo, utils.FilesInfo{Name: common.KprobeSourceFile, Data: ksf})
	bph, err := box.FindString(common.BpfHeaderFile)
	if err != nil {
		return nil, fmt.Errorf("faild to load bpf header file %s", err.Error())
	}
	filesInfo = append(filesInfo, utils.FilesInfo{Name: common.BpfHeaderFile, Data: bph})
	bphh, err := box.FindString(common.BpfHelperHeaderFile)
	if err != nil {
		return nil, fmt.Errorf("faild to load bpg helper header file %s", err.Error())
	}
	filesInfo = append(filesInfo, utils.FilesInfo{Name: common.BpfHelperHeaderFile, Data: bphh})
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
