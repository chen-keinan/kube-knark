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

//Executor defines the interface for command exec
//clang.go
//go:generate mockgen -destination=./mocks/mock_Executor.go -package=mocks . Executor
type Executor interface {
	Run() error
}

//ClangExecutor defines the interface for clang compiler
//clang.go
//go:generate mockgen -destination=../mocks/mock_ClangExecutor.go -package=mocks . ClangExecutor
type ClangExecutor interface {
	CompileSourceToElf(cmd Executor) (*CommandResult, error)
	NewExecCommand(args ...string) Executor
}

//ClangCompiler object
type ClangCompiler struct {
}

//NewClangCompiler return new instance of shell executor
func NewClangCompiler() ClangExecutor {
	return &ClangCompiler{}
}

//CommandResult return data
type CommandResult struct {
	Stdout string
	Stderr string
}

//NewExecCommand return new exec Command instance
// #nosec
func (ce ClangCompiler) NewExecCommand(args ...string) Executor {
	fullCmd := fmt.Sprintf(Command, args[0], args[1], args[2])
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

// return cmd args
func cmdArgs(executor Executor) []string {
	c, ok := executor.(*exec.Cmd)
	if !ok {
		return []string{}
	}
	return c.Args
}
