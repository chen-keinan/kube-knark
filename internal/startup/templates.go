package startup

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/common"
	shell "github.com/chen-keinan/kube-knark/internal/compiler"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/gobuffalo/packr"
	"os"
	"path/filepath"
	"strings"
)

//GenerateK8sBenchmarkFiles use packr to load benchmark audit test yaml
func GenerateEbpfFiles() ([]utils.FilesInfo, error) {
	fileInfo := make([]utils.FilesInfo, 0)
	box := packr.NewBox("./../../ebpf/")
	// Add ebpf kprobe program
	ksf, err := box.FindString(common.KprobeSourceFile)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load ebpf source file %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.KprobeSourceFile, Data: ksf})
	// Add bpf header file
	bh, err := box.FindString(common.BpfHeaderFile)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load ebpf source file %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.BpfHeaderFile, Data: bh})
	// Add bph_helper header file
	bhh, err := box.FindString(common.BpfHelperHeaderFile)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load ebpf source file %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.BpfHelperHeaderFile, Data: bhh})
	return fileInfo, nil
}

//SaveEbpfFilesIfNotExist create ebpf source files if not exist
func SaveEbpfFilesIfNotExist(filesData []utils.FilesInfo) error {
	ebpfFolder := utils.GetEbpfSourceFolder()
	for _, fileData := range filesData {
		filePath := filepath.Join(ebpfFolder, fileData.Name)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			f, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf(err.Error())
			}
			_, err = f.WriteString(fileData.Data)
			if err != nil {
				return fmt.Errorf("failed to write benchmark file")
			}
			err = f.Close()
			if err != nil {
				return fmt.Errorf("faild to close file %s", filePath)
			}
		}
	}
	return nil
}

//CompileEbpfSources compile ebpf program to elf file
func CompileEbpfSources(filesData []utils.FilesInfo) error {
	ebpfSourceFolder := utils.GetEbpfSourceFolder()
	ebpfCompiledFolder := utils.GetEbpfCompiledFolder()
	for _, fileData := range filesData {
		if !strings.HasSuffix(fileData.Name, ".c") {
			continue
		}
		sourcefilePath := filepath.Join(ebpfSourceFolder, fileData.Name)
		compiledFilePath := filepath.Join(ebpfCompiledFolder, strings.Replace(fileData.Name, ".c", ".elf", -1))
		cmdResult, err := shell.NewClangCompiler().CompileSourceToElf(fmt.Sprintf(".%s", ebpfSourceFolder), sourcefilePath, compiledFilePath)
		if cmdResult.Stderr != "" || err != nil {
			return fmt.Errorf(cmdResult.Stderr)
		}
	}
	return nil
}
