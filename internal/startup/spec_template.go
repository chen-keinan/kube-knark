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

	return fileInfo, nil
}
