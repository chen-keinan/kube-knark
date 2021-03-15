package shell

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ExecCommand(t *testing.T) {
	exec := NewExecCommand("a", "b", "c")
	args := cmdArgs(exec)
	assert.Equal(t, args[0], "bash")
	assert.Equal(t, args[1], "-c")
	assert.Equal(t, args[2], "clang -Ia -O2 -target bpf -c b -o c")
}

func TestExecRun(t *testing.T) {
	ctl := gomock.NewController(t)
	cmd := mocks.NewMockExecutor(ctl)
	cmd.EXPECT().Run().Return(nil).Times(1)
	re, err := NewClangCompiler().CompileSourceToElf(cmd)
	assert.NoError(t, err)
	assert.Equal(t, re.Stdout, "")
}

func TestExecRunError(t *testing.T) {
	ctl := gomock.NewController(t)
	cmd := mocks.NewMockExecutor(ctl)
	cmd.EXPECT().Run().Return(fmt.Errorf("failed to exec command")).Times(1)
	_, err := NewClangCompiler().CompileSourceToElf(cmd)
	assert.Error(t, err)
}
