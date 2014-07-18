package builtin

import (
	"bytes"
	"os"
	"os/exec"
	"syscall"
)

func execCode(err  error) int {
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
            return status.ExitStatus()
        }
	}
	return 0
}

func Exec(exe string, args ...string) (error, int) {
	cmd := exec.Command(exe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return err, execCode(err)
}

func ExecWithOutput(exe string, args ...string) (stdout, stderr []byte, err error, errCode int) {
	var stdoutBuf, stderrBuf bytes.Buffer
	
	cmd := exec.Command(exe, args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	
	return stdoutBuf.Bytes(), stderrBuf.Bytes(), err, execCode(err)
}