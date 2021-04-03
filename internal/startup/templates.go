package startup

import (
	"fmt"
	shell "github.com/chen-keinan/kube-knark/internal/compiler"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

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
