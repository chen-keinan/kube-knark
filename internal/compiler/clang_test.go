package shell

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Exec(t *testing.T) {
	se := NewClangCompiler()
	cmd := NewExecCommand("", "", "")
	_, err := se.CompileSourceToElf(cmd)
	assert.Error(t, err)
}
