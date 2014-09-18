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
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

var (
	// Set DryRun to true will make all Exec related functions in dry run mode, i.e. they only print the command to run.
	DryRun = false
)

/*
S converts anything into a string. If args is specified, v is used as a format
string.
*/
func S(v interface{}, args ...interface{}) string {
	return fmt.Sprintf(fmt.Sprint(v), args...)
}

/*
S2Is returns a slice of interface{} given strings.
*/
func S2Is(args ...string) (ifs []interface{}) {
	ifs = make([]interface{}, len(args))
	for i, arg := range args {
		ifs[i] = arg
	}
	return
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

If the command is executed successfuly without mistakes, (nil, 0) will be
returned. Otherwise, the error and error code will be returned.
NOTE the error code could be 0 with a non-nil error.

Stdout/stderr are directed to the current stdout/stderr.
*/
func Exec(exe interface{}, args ...string) (error, int) {
	if DryRun {
		Eprintfln("Exec: " + S(exe), S2Is(args...)...)
		return nil, 0
	}
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
	if DryRun {
		Eprintfln("ExecWithStdout" + S(exe), S2Is(args...)...)
		return "", nil, 0
	}
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
	if DryRun {
		Eprintfln("ExecWithStdErrOut" + S(exe), S2Is(args...)...)
		return "", "", nil, 0
	}
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
BashEval is similar to BashWithStdout but only returns captured stdout as a
string. Trainling newlines are deleted. It's like the backtick substitution in
Bash.
*/
func BashEval(cmd interface{}, args ...interface{}) string {
	out, _, _ := BashWithStdout(cmd, args...)
	return strings.TrimRight(out, "\r\n")
}

/*
Pwd is similar to os.Getwd() without error returned.
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
Fatalf prints a message and let the program exit with DefExitCode.
*/
func Fatalf(msg interface{}, args ...interface{}) {
	fmt.Fprintln(os.Stderr, S(msg, args...))
	os.Exit(DefExitCode)
}

/*
Eprintf is similar to fmt.Printf but the output is stderr.
*/
func Eprintf(format interface{}, args ...interface{}) {
	fmt.Fprint(os.Stderr, S(format, args...))
}

/*
Eprint is similar to fmt.Print but the output is stderr.
*/
func Eprint(args ...interface{}) {
	fmt.Fprint(os.Stderr, args...)
}

/*
Eprintln is similar to fmt.Println but the output is stderr.
*/
func Eprintln(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

/*
Eprintfln is similar to Eprintf but with a trailing new-line printed
*/
func Eprintfln(format interface{}, args ...interface{}) {
	fmt.Fprintln(os.Stderr, S(format, args...))
}

/*
Printfln is similar to fmt.Printf but with a trailing new-line printed
*/
func Printfln(format interface{}, args ...interface{}) {
	fmt.Println(S(format, args...))
}

/*
Succ check whether return values of Exec and friends means succeed.
*/
func Succ(err error, code int) bool {
	return err == nil
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
ScriptName returns the filename of the current script not including the path.
*/
func ScriptName() string {
	return path.Base(os.Args[0])
}

/*
Exists checks whether the path exists
*/
func Exists(p interface{}, args ...interface{}) bool {
	_, err := os.Stat(S(p, args...))
	return err == nil
}

/*
IsDir returns true only if the path exists and indicates a directory
*/
func IsDir(p interface{}, args ...interface{}) bool {
	info, err := os.Stat(S(p, args...))
	if err != nil {
		// the path does not exist
		return false
	}
	return info.Mode().IsDir()
}

/*
IsFile returns true only if the path exists and indicates a file
*/
func IsFile(p interface{}, args ...interface{}) bool {
	info, err := os.Stat(S(p, args...))
	if err != nil {
		// the path does not exist
		return false
	}
	return !info.Mode().IsDir()
}

/*
Match use regular expression pattern to match str and returns all capturing groups.
*/
func Match(str interface{}, pattern interface{}) []string {
	return regexp.MustCompile(S(pattern)).FindStringSubmatch(S(str))
}
