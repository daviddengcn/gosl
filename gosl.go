package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strings"
	"os/exec"
	"syscall"

	"github.com/daviddengcn/go-villa"
)

const (
	STAGE_READY = iota
	STAGE_IMPORT
	STAGE_MAIN
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

func execCode(err  error) int {
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
            return status.ExitStatus()
        }
	}
	return 0
}

func process() error {
	fn := villa.Path(os.Args[1])
	buffer, err := ioutil.ReadFile(fn.S())
	if err != nil {
		return err
	}

	var code bytes.Buffer
	code.WriteString(`package main; import . "fmt"; import . "os"; import . "github.com/daviddengcn/gosl/builtin"; import . "strings"; `)
	code.WriteString(`func init() {_ = Printf; _ = Exit; _ = Exec; _ = Contains; } `)

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
			if _, err := code.WriteRune('\n'); err != nil {
				return err
			}
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
				if !bytes.HasPrefix(line, []byte("import ")) {
					stage = STAGE_MAIN
				}
			}
			break
		}

		if stage == STAGE_MAIN {
			if _, err := code.WriteString("func main() { "); err != nil {
				return err
			}
		}
		if _, err := code.Write(line); err != nil {
			return err
		}
		if _, err := code.WriteRune('\n'); err != nil {
			return err
		}
		if stage == STAGE_MAIN {
			break
		}
	}
	if stage == STAGE_MAIN {
		if _, err := code.Write(buf); err != nil {
			return err
		}
	} else {
		if _, err := code.WriteString("\nfunc main() { "); err != nil {
			return err
		}
	}
	
	if _, err := code.WriteString("\n}\n"); err != nil {
		return err
	}

	codeFn := genFilename(fn.Base())
	if err := codeFn.WriteFile(code.Bytes(), 0644); err != nil {
		return err
	}
	defer codeFn.Remove()

	cmd := villa.Path("go").Command(append([]string{"run", codeFn.S()}, os.Args[2:]...)...)
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
	if len(os.Args) < 2 {
		return
	}

	if err := process(); err != nil {
		fmt.Printf("Failed: %v\n", err)
	}
}
