package startup

import (
	"github.com/chen-keinan/kube-knark/internal/common"
	shell "github.com/chen-keinan/kube-knark/internal/compiler"
	"github.com/chen-keinan/kube-knark/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

//Test_CreateEbfFilesIfNotExist test
func Test_CreateEbfFilesIfNotExist(t *testing.T) {
	bFiles, err := GenerateEbpfFiles()
	if err != nil {
		t.Fatal(err)
	}
	err = SaveFilesIfNotExist(bFiles, utils.GetEbpfSourceFolder)
	if err != nil {
		t.Fatal(err)
	}
	// generate test with packr
	assert.Equal(t, bFiles[0].Name, common.KprobeSourceFile)
	assert.Equal(t, bFiles[1].Name, common.BpfHeaderFile)
	assert.Equal(t, bFiles[2].Name, common.BpfHelperHeaderFile)
}

//Test_CompileEbpfSources test
func Test_CompileEbpfSources(t *testing.T) {
	bFiles, err := GenerateEbpfFiles()
	if err != nil {
		t.Fatal(err)
	}
	err = CompileEbpfSources(bFiles, shell.NewClangCompiler())
	if err != nil {
		t.Fatal(err)
	}

	ebpfCompiledFolder, err := utils.GetEbpfCompiledFolder(utils.NewKFolder())
	if err != nil {
		t.Fatal(err)
	}
	cfiles, err := utils.GetFiles(ebpfCompiledFolder)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, cfiles[0].Name, common.KprobeCompiledFile)

}
func TestGenerateEbpfFiles(t *testing.T) {
	files, err := GenerateEbpfFiles()
	assert.NoError(t, err)
	assert.Equal(t, files[0].Name, common.KprobeSourceFile)
	assert.Equal(t, files[1].Name, common.BpfHeaderFile)
	assert.Equal(t, files[2].Name, common.BpfHelperHeaderFile)
}

func TestSaveFilesIfNotExist(t *testing.T) {
	files, err := GenerateEbpfFiles()
	assert.NoError(t, err)
	fm := utils.NewKFolder()
	err = utils.CreateKubeKnarkFolders(fm)
	assert.NoError(t, err)
	err = SaveFilesIfNotExist(files, utils.GetEbpfSourceFolder)
	assert.NoError(t, err)
}
