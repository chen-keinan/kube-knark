package shell

import (
	"bytes"
	"fmt"
	"os/exec"
)

const command = "clang -I%s -O2 -target bpf -c %s -o %s"

//ShellToUse bash shell
const ShellToUse = "bash"

//Executor defines the interface for clang compiler
//exec.go
//go:generate mockgen -destination=../mocks/mock_Executor.go -package=mocks . Executor
type Executor interface {
	CompileSourceToElf(headerPath,source,destination string) (*CommandResult, error)
}

//ClangCompiler object
type ClangCompiler struct {
}

//NewClangCompiler return new instance of shell executor
func NewClangCompiler() Executor {
	return &ClangCompiler{}
}

//CommandResult return data
type CommandResult struct {
	Stdout string
	Stderr string
}

//CompileSourceToElf execute shell command
// #nosec
func (ce ClangCompiler) CompileSourceToElf(headerPath, source, destination string) (*CommandResult, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	fullCmd := fmt.Sprintf(command,headerPath,source, destination)
	cmd := exec.Command(ShellToUse, "-c", fullCmd)
	fmt.Println(cmd.String())
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return &CommandResult{Stdout: stdout.String(), Stderr: stderr.String()}, err
}
