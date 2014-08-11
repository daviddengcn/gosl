/*
This is an application that can make you write script with Go languages.

It is NOT an interpreter but the pure Go. The preprocessor tranforms the script into a Go program, instantly compile and run. So it is almost same as the standard Go with the same efficiency.
*/
package main

import (
	"bytes"
	"fmt"
	"flag"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/daviddengcn/go-villa"
)

const (
	STAGE_READY  = iota // start
	STAGE_IMPORT        // import sentences
	STAGE_MAIN          // main part
)

func genFilename(suffix villa.Path) villa.Path {
	if !strings.HasSuffix(suffix.S(), ".go") {
		suffix += ".go"
	}
	dir := villa.Path(os.TempDir())
	for {
		base := villa.Path(fmt.Sprintf("gosl-%08x-%s", rand.Int63n(math.MaxInt64), suffix))
		fn := dir.Join(base)
		if !fn.Exists() {
			return fn
		}
	}
}

func execCode(err error) int {
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 0
}

var (
	DEFAULT_IMPORT = []string{
		"fmt", "Printf",
		"os", "Exit",
		"strings", "Contains",
		"math", "Abs",
		"strconv", "Atoi",
		"time", "Sleep",
		"github.com/daviddengcn/gosl/builtin", "Exec",
	}
)

func appendInitAndMainHead(code *bytes.Buffer) {
	code.WriteString(`func init() {`)
	code.WriteString(` Args = Args[1:];`)
	for i := 1; i < len(DEFAULT_IMPORT); i += 2 {
		code.WriteString(` _ = `)
		code.WriteString(DEFAULT_IMPORT[i])
		code.WriteString(`;`)
	}
	code.WriteString(` }; `)
	code.WriteString("func main() { ")
}

type Options struct {
	ShowSource bool
	NoClean bool
}

func Process(gf string, opts Options) error {
	fn := villa.Path(gf)
	buffer, err := ioutil.ReadFile(fn.S())
	if err != nil {
		return err
	}

	var code bytes.Buffer
	code.WriteString(`package main;`)
	for i := 0; i < len(DEFAULT_IMPORT); i += 2 {
		code.WriteString(` import . "`)
		code.WriteString(DEFAULT_IMPORT[i])
		code.WriteString(`";`)
	}

	stage := STAGE_READY

	buf := buffer
	for len(buf) > 0 {
		p := bytes.IndexByte(buf, byte('\n'))
		var line []byte
		if p < 0 {
			line = buf
			buf = nil
		} else {
			line = buf[:p]
			buf = buf[p+1:]
		}

		if len(line) == 0 {
			code.WriteRune('\n')
			continue
		}

		for {
			switch stage {
			case STAGE_READY:
				if line[0] != '#' {
					stage = STAGE_IMPORT
					continue
				}
				line = nil
			case STAGE_IMPORT:
				trimmed := bytes.TrimSpace(line)
				if len(trimmed) > 0 && !bytes.HasPrefix(trimmed, []byte("import ")) && !bytes.HasPrefix(trimmed, []byte("//")) {
					stage = STAGE_MAIN
				}
			}
			break
		}

		if stage == STAGE_MAIN {
			appendInitAndMainHead(&code)
		}
		code.Write(line)
		code.WriteRune('\n')
		if stage == STAGE_MAIN {
			break
		}
	}
	if stage == STAGE_MAIN {
		code.Write(buf)
	} else {
		appendInitAndMainHead(&code)
	}

	code.WriteString("\n}\n")
	
	if opts.ShowSource {
		fmt.Println(string(code.Bytes()))
	}

	codeFn := genFilename(fn.Base())
	if err := codeFn.WriteFile(code.Bytes(), 0644); err != nil {
		return err
	}
	if !opts.NoClean {
		defer codeFn.Remove()
	}

	exeFn := codeFn + ".exe" // to be compatible with Windows

	cmd := villa.Path("go").Command("build", "-o", exeFn.S(), codeFn.S())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}
	if !opts.NoClean {
		defer exeFn.Remove()
	}

	if p, err := filepath.Abs(os.Args[1]); err == nil {
		os.Args[1] = p
	}
	cmd = exeFn.Command(os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	ec := execCode(err)
	if ec != 0 {
		os.Exit(ec)
	}
	return err
}

func main() {
	var opts Options
	flag.BoolVar(&opts.ShowSource, "showsource", false, "Show generated source code")
	flag.BoolVar(&opts.NoClean, "noclean", false, "No cleaning of generated files")
	
	flag.Parse()
	
	if len(flag.Args()) < 1 {
		return
	}

	if err := Process(flag.Args()[0], opts); err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(-1)
	}
}
