package builtin

import (
	"fmt"
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

func S(v interface{}) string {
	return fmt.Sprint(v)
}

func Exec(exe interface{}, args ...string) (error, int) {
	cmd := exec.Command(S(exe), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return err, execCode(err)
}

func ExecWithStdout(exe interface{}, args ...string) (stdout string, err error, errCode int) {
	var stdoutBuf bytes.Buffer
	
	cmd := exec.Command(S(exe), args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	
	return string(stdoutBuf.Bytes()), err, execCode(err)
}

func ExecWithStdErrOut(exe interface{}, args ...string) (stdout, stderr string, err error, errCode int) {
	var stdoutBuf, stderrBuf bytes.Buffer
	
	cmd := exec.Command(S(exe), args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	
	return string(stdoutBuf.Bytes()), string(stderrBuf.Bytes()), err, execCode(err)
}

func Bash(cmd interface{}) (error, int) {
	return Exec("bash", "-c", S(cmd))
}

func BashWithStdout(cmd interface{}) (string, error, int) {
	return ExecWithStdout("bash", "-c", S(cmd))
}
