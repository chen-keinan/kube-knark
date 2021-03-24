package startup

import (
	"github.com/chen-keinan/kube-knark/internal/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateSpecFiles(t *testing.T) {
	specFiles, err := GenerateSpecFiles()
	assert.NoError(t, err)
	assert.Equal(t, specFiles[0].Name, common.Workload)
}
func TestGenerateFileSystemSpec(t *testing.T) {
	specFiles, err := GenerateFileSystemSpec()
	assert.NoError(t, err)
	assert.Equal(t, specFiles[0].Name, common.ConfigFilesPermission)
}
