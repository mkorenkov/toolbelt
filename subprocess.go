package toolbelt

import (
	"fmt"
	"io"
	"os/exec"
	"time"
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

// ExecuteWithTimeout command with timeout in seconds and explicit stdout, stderr and stdin
func ExecuteWithTimeout(command string, timeoutSeconds int, stdout io.Writer, stderr io.Writer, stdin io.Reader, arguments ...string) (timeout bool, err error) {
	cmd := exec.Command(command, arguments...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = stdin
	done := make(chan error, 1)
	err = cmd.Start()
	if err != nil {
		return false, err
	}
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		err := cmd.Process.Kill()
		if err != nil {
			return true, fmt.Errorf("failed to kill: %v", err)
		}
		return true, nil
	case err := <-done:
		return false, err
	}
}
