package builtin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/daviddengcn/go-assert"
)

func TestI(t *testing.T) {
	assert.Equals(t, "I(123)", I(123), 123)
	assert.Equals(t, "I(int64(123))", I(int64(123)), 123)
	assert.Equals(t, "I(abc)", I("abc"), 0)
	assert.Equals(t, "I(abc, 50)", I("abc", 50), 50)
	assert.Equals(t, "I(abc, def)", I("abc", "def"), 0)
}

func TestS(t *testing.T) {
	assert.Equals(t, "S(123)", S(123), "123")
	assert.Equals(t, "S(a%db, 123)", S("a%db", 123), "a123b")
}

func TestEprint(t *testing.T) {
	Eprintf("Abc %d def\n", 123)
	Eprint("%d% --- ")
	Eprintln("Hello world %d!")
}

func TestMustSucc(t *testing.T) {
	MustSucc(nil, 0)
	MustSucc(nil, 1)
}

func TestSucc(t *testing.T) {
	assert.Equals(t, "Succ(nil, 0)", Succ(nil, 0), true)
}

func TestSortF(t *testing.T) {
	ints := []int{3, 4, 1, 7, 0}
	SortF(len(ints), func(i, j int) bool {
		return ints[i] < ints[j]
	}, func(i, j int) {
		ints[i], ints[j] = ints[j], ints[i]
	})
	assert.StringEquals(t, "ints", ints, []int{0, 1, 3, 4, 7})

	ints = []int{3, 4, 1, 7, 0}
	SortF(len(ints), func(i, j int) bool {
		return ints[i] > ints[j]
	}, func(i, j int) {
		ints[i], ints[j] = ints[j], ints[i]
	})
	assert.StringEquals(t, "ints", ints, []int{7, 4, 3, 1, 0})
}

func TestExists(t *testing.T) {
	tmpDir := os.TempDir()
	assert.Equals(t, "Exists(tmpDir)", Exists(tmpDir), true)
	assert.Equals(t, "Exists(tmpDir-nonexists)", Exists(tmpDir+"-nonexists"), false)
}

func TestIsDirFile(t *testing.T) {
	tmpDir := os.TempDir()
	assert.Equals(t, "IsDir(tmpDir)", IsDir(tmpDir), true)
	assert.Equals(t, "IsFile(tmpDir)", IsFile(tmpDir), false)

	assert.Equals(t, "IsDir(tmpDir-nonexists)", IsDir(tmpDir+"-nonexists"), false)
	assert.Equals(t, "IsFile(tmpDir-nonexists)", IsFile(tmpDir+"-nonexists"), false)

	fn := filepath.Join(tmpDir, "file")
	if f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644); err == nil {
		f.Close()
		assert.Equals(t, "IsDir(tmpDir/file)", IsDir(fn), false)
		assert.Equals(t, "IsFile(tmpDir/file)", IsFile(fn), true)
	}
}

func TestMatch(t *testing.T) {
	assert.LinesEqual(t, "Match", Match("AAAabc123efgFFF", "[a-z]+([0-9]+)[a-z]+"), []string{
		"abc123efg",
		"123",
	})
}

func TestDryRun(t *testing.T) {
	DryRun = true
	Bash("ls -l")
	Exec("myapp", "arg1", "arg2")
}
