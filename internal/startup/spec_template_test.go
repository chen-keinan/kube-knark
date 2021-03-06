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
	assert.Equal(t, specFiles[1].Name, common.Services)
	assert.Equal(t, specFiles[2].Name, common.ConfigAndStorage)
	assert.Equal(t, specFiles[3].Name, common.Authentication)
	assert.Equal(t, specFiles[4].Name, common.Authorization)
	assert.Equal(t, specFiles[5].Name, common.Policy)
	assert.Equal(t, specFiles[6].Name, common.Extend)
	assert.Equal(t, specFiles[7].Name, common.Cluster)
}
func TestGenerateFileSystemSpec(t *testing.T) {
	specFiles, err := GenerateFileSystemSpec()
	assert.NoError(t, err)
	assert.Equal(t, specFiles[0].Name, common.ConfigFilesPermission)
}
