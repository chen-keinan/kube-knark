package shell

import (
	"bytes"
	"fmt"
	"os/exec"
)

const (
	//Command compile command
	Command = "clang -I%s -O2 -target bpf -c %s -o %s"
	//ShellToUse bash shell
	ShellToUse = "bash"
)

//Executor defines the interface for clang compiler
//exec.go
//go:generate mockgen -destination=../mocks/mock_Executor.go -package=mocks . Executor
type Executor interface {
	Run() error
}

//ClangCompiler object
type ClangCompiler struct {
}

//NewClangCompiler return new instance of shell executor
func NewClangCompiler() *ClangCompiler {
	return &ClangCompiler{}
}

//CommandResult return data
type CommandResult struct {
	Stdout string
	Stderr string
}

//NewExecCommand return new exec Command instance
// #nosec
func NewExecCommand(headerPath, source, destination string) Executor {
	fullCmd := fmt.Sprintf(Command, headerPath, source, destination)
	return exec.Command(ShellToUse, "-c", fullCmd)
}

//CompileSourceToElf execute shell Command
func (ce ClangCompiler) CompileSourceToElf(cmd Executor) (*CommandResult, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return &CommandResult{Stdout: stdout.String(), Stderr: stderr.String()}, err
}
