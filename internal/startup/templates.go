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

//GenerateEbpfFiles ebpf file from template
func GenerateEbpfFiles() ([]utils.FilesInfo, error) {
	filesInfo := make([]utils.FilesInfo, 0)
	kprobFile, err := createFileFromTemplate("./../../ebpf/", common.KprobeSourceFile)
	if err != nil {
		return nil, err
	}
	filesInfo = append(filesInfo, kprobFile...)
	// Add bpf header file
	bpfHeaderFile, err := createFileFromTemplate("./../../ebpf/", common.BpfHeaderFile)
	if err != nil {
		return nil, err
	}
	filesInfo = append(filesInfo, bpfHeaderFile...)
	// Add bph_helper header file
	bpfHeaderHelperFile, err := createFileFromTemplate("./../../ebpf/", common.BpfHelperHeaderFile)
	if err != nil {
		return nil, err
	}
	filesInfo = append(filesInfo, bpfHeaderHelperFile...)
	return filesInfo, nil
}

func createFileFromTemplate(path, fileName string) ([]utils.FilesInfo, error) {
	fileInfo := make([]utils.FilesInfo, 0)
	box := packr.NewBox(path)
	// Add ebpf kprobe program
	ksf, err := box.FindString(fileName)
	if err != nil {
		return nil, fmt.Errorf("faild to load ebpf source file %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: fileName, Data: ksf})
	return fileInfo, nil
}

//SaveFilesIfNotExist save files if not exist
func SaveFilesIfNotExist(filesData []utils.FilesInfo, f func() (string, error)) error {
	folder, err := f()
	if err != nil {
		return err
	}
	for _, fileData := range filesData {
		filePath := filepath.Join(folder, fileData.Name)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			f, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf(err.Error())
			}
			_, err = f.WriteString(fileData.Data)
			if err != nil {
				return fmt.Errorf("failed to write files")
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
func CompileEbpfSources(filesData []utils.FilesInfo, cc shell.ClangExecutor) error {
	ebpfSourceFolder, err := utils.GetEbpfSourceFolder()
	if err != nil {
		return err
	}
	ebpfCompiledFolder, err := utils.GetEbpfCompiledFolder()
	if err != nil {
		return err
	}
	for _, fileData := range filesData {
		if !strings.HasSuffix(fileData.Name, ".c") {
			continue
		}
		sourcefilePath := filepath.Join(ebpfSourceFolder, fileData.Name)
		compiledFilePath := filepath.Join(ebpfCompiledFolder, strings.Replace(fileData.Name, ".c", ".elf", -1))
		cmd := cc.NewExecCommand(fmt.Sprintf(".%s", ebpfSourceFolder), sourcefilePath, compiledFilePath)
		cmdResult, err := cc.CompileSourceToElf(cmd)
		if cmdResult.Stderr != "" || err != nil {
			return fmt.Errorf(cmdResult.Stderr)
		}
	}
	return nil
}
