package toolbelt

import (
	"io"
	"os/exec"
)

// OutputOf executes shell command with arguments
func OutputOf(command string, arguments ...string) (out string, err error) {
	cmdOut, err := exec.Command(command, arguments...).CombinedOutput()
	if err != nil {
		return string(cmdOut), err
	}
	return string(cmdOut), nil
}

// Execute command and pass stdout, stderr and stdin
func Execute(command string, stdout io.Writer, stderr io.Writer, stdin io.Reader, arguments ...string) error {
	cmd := exec.Command(command, arguments...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = stdin
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
