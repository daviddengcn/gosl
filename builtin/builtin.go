package builtin

import (
	"fmt"
	"bytes"
	"os"
	"os/exec"
	"syscall"
	"strconv"
)

func S(v interface{}) string {
	return fmt.Sprint(v)
}

func I(v interface{}) int {
	if i, ok := v.(int); ok {
		return i
	}
	if i, ok := v.(int64); ok {
		return int(i)
	}
	
	i, _ := strconv.Atoi(S(v))
	return i
}

func AtoiDef(a, def interface{}) int {
	i, err := strconv.Atoi(S(a))
	if err != nil {
		i = I(def)
	}
	return i
}

func execCode(err  error) int {
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
            return status.ExitStatus()
        }
	}
	return 0
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
