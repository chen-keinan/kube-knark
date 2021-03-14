package shell

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Exec(t *testing.T) {
	se := NewClangCompiler()
	execResult, _ := se.CompileSourceToElf("echo test", "", "")
	assert.Equal(t, execResult.Stdout, "")
	assert.True(t, len(execResult.Stderr) > 0)
}
