/*
Package builtin implements some built-in functions for gosl (Go Language
Script Language, github.com/daviddengcn/gosl)

For use of convinience as a script language, the parameters are commonly
defined as an interface{}.
*/
package builtin

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

/*
S converts anything into a string.
*/
func S(v interface{}) string {
	return fmt.Sprint(v)
}

/*
I converts anything into an int. When the value is malformed, if the optional
default value is specified, it is converted to int and returned; otherwise,
0 is returned.
*/
func I(v interface{}, def ...interface{}) int {
	if i, ok := v.(int); ok {
		return i
	}
	if i, ok := v.(int64); ok {
		return int(i)
	}

	i, err := strconv.Atoi(S(v))
	if err != nil && len(def) > 0 {
		return I(def[0])
	}
	return i
}

func execCode(err error) int {
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 0
}

/*
Exec runs a command. exe is the path to the executable and args are arugment
passed to it.

If the command is executed successfuly without mistake, (nil, 0)
will be returned. Otherwise, the error and error code will be returned.

Stdout/stderr are directed the current stdout/stderr.
*/
func Exec(exe interface{}, args ...string) (error, int) {
	cmd := exec.Command(S(exe), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return err, execCode(err)
}

/*
ExecWithStdout is similar to Exec but the stdout is captured and returned as
the first return value.
*/
func ExecWithStdout(exe interface{}, args ...string) (stdout string, err error, errCode int) {
	var stdoutBuf bytes.Buffer

	cmd := exec.Command(S(exe), args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Run()

	return string(stdoutBuf.Bytes()), err, execCode(err)
}

/*
ExecWithStdout is similar to Exec but the stdout/stderr are captured and
returned as the first/second return values.
*/
func ExecWithStdErrOut(exe interface{}, args ...string) (stdout, stderr string, err error, errCode int) {
	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.Command(S(exe), args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	cmd.Stdin = os.Stdin
	err = cmd.Run()

	return string(stdoutBuf.Bytes()), string(stderrBuf.Bytes()), err, execCode(err)
}

/*
Eval is similar to ExecWithStdout but with stdout captured and returned as a
string. Trainling newlines are deleted.
*/
func Eval(exe interface{}, args ...string) string {
	out, _, _ := ExecWithStdout(exe, args...)
	return strings.TrimRight(out, "\r\n")
}

/*
Bash runs a command with bash. Return values are defined in Exec.
*/
func Bash(cmd interface{}) (error, int) {
	return Exec("bash", "-c", S(cmd))
}

/*
BashWithStdout is similar to Bash but with stdout captured and returned as a
string.
*/
func BashWithStdout(cmd interface{}) (string, error, int) {
	return ExecWithStdout("bash", "-c", S(cmd))
}

/*
BashEval is similar to BashWithStdout but with stdout captured and returned
as a string. Trainling newlines are deleted.
*/
func BashEval(cmd interface{}) string {
	out, _, _ := BashWithStdout(cmd)
	return strings.TrimRight(out, "\r\n")
}

/*
Similar to os.Getwd() but no error returned.
*/
func Pwd() string {
	pwd, _ := os.Getwd()
	return pwd
}
