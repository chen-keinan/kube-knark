package shell

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Exec(t *testing.T) {
	se := NewClangCompiler()
	execResult, err := se.CompileSourceToElf("echo test")
	assert.NoError(t, err)
	assert.Equal(t, execResult.Stdout, "test\n")
	assert.Equal(t, execResult.Stderr, "")
}
