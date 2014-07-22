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
	"path"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

/*
S converts anything into a string. If args is specified, v is used as a format
string.
*/
func S(v interface{}, args ...interface{}) string {
	return fmt.Sprintf(fmt.Sprint(v), args...)
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
func Bash(cmd interface{}, args ...interface{}) (error, int) {
	return Exec("bash", "-c", S(cmd, args...))
}

/*
BashWithStdout is similar to Bash but with stdout captured and returned as a
string.
*/
func BashWithStdout(cmd interface{}, args ...interface{}) (string, error, int) {
	return ExecWithStdout("bash", "-c", S(cmd, args...))
}

/*
BashEval is similar to BashWithStdout but with stdout captured and returned
as a string. Trainling newlines are deleted.
*/
func BashEval(cmd interface{}, args ...interface{}) string {
	out, _, _ := BashWithStdout(cmd, args...)
	return strings.TrimRight(out, "\r\n")
}

/*
Similar to os.Getwd() but no error returned.
*/
func Pwd() string {
	pwd, _ := os.Getwd()
	return pwd
}

/*
DefExitCode is the default exit code.
*/
var DefExitCode = 1

/*
Fatalf print a message and exit the program with DefExitCode.
*/
func Fatalf(msg interface{}, args ...interface{}) {
	fmt.Fprintln(os.Stderr, S(msg, args...))
	os.Exit(DefExitCode)
}

/*
Eprintf is similar to fmt.Printf but output is stderr.
*/
func Eprintf(format interface{}, args ...interface{}) {
	fmt.Fprint(os.Stderr, S(format, args...))
}

/*
Eprint is similar to fmt.Print but output is stderr.
*/
func Eprint(args ...interface{}) {
	fmt.Fprint(os.Stderr, args...)
}

/*
Eprintln is similar to fmt.Println but output is stderr.
*/
func Eprintln(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

/*
MustSucc checks the result of Exec/Bash. If not succeed, exit the application.
*/
func MustSucc(err error, code int) {
	if err == nil {
		return
	}

	if code != 0 {
		Fatalf("Failed with error code: %d", code)
	}

	Fatalf("Failed with error: %v", err)
}

type sortI struct {
	l    int
	less func(int, int) bool
	swap func(int, int)
}

func (s *sortI) Len() int {
	return s.l
}

func (s *sortI) Less(i, j int) bool {
	return s.less(i, j)
}

func (s *sortI) Swap(i, j int) {
	s.swap(i, j)
}

/*
SortF sorts the data defined by the length, Less and Swap functions.
*/
func SortF(Len int, Less func(int, int) bool, Swap func(int, int)) {
	sort.Sort(&sortI{l: Len, less: Less, swap: Swap})
}

/*
ScriptDir returns the folder of the current script.
*/
func ScriptDir() string {
	return path.Dir(os.Args[0])
}

/*
Exists checks whether the path exists
*/
func Exists(p interface{}) bool {
	_, err := os.Stat(S(p))
	return err == nil
}

/*
IsDir returns true only if the path exists and indicates a directory
*/
func IsDir(p interface{}) bool {
	info, err := os.Stat(S(p))
	if err != nil {
		// the path does not exist
		return false
	}
	return info.Mode().IsDir()
}

/*
IsFile returns true only if the path exists and indicates a file
*/
func IsFile(p interface{}) bool {
	info, err := os.Stat(S(p))
	if err != nil {
		// the path does not exist
		return false
	}
	return !info.Mode().IsDir()
}
