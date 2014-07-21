gosl
====

This is an application that can make you write script with Go languages.

It is *NOT* an interpreter but the pure Go. The preprocessor tranforms the script into a Go program, instantly compile and run. So it is almost same as the standard Go with the same efficiency.

Benefit
-------
1. Pure Go language. No need to learn a new script language.
1. Pre-imported packages and pre-defined functions make it easy to code.
1. Seamless integration with the Go project. E.g. can easily load configuratio or data file from the Go project.
1. Running efficiency same as Go, much faster than Python.

Example
-------

```go
#!/bin/gosl

import "encoding/json"

toJson := func(lines []string) string {
  res, _ := json.Marshal(struct {
    Lines []string `json:"lines"`
  }{
    Lines: lines,
  })
  return string(res)
}

files, _, _ := BashWithStdout("ls -l /tmp/")

Println(toJson(Split(files, "\n")))
    
```

[Go Search](http://go-search.org/) is now operating with gosl. So you can find some good examples at: https://github.com/daviddengcn/gcse/tree/master/scripts

Installation and Usage
----------------------

#### Download and install the package
```bash
go get github.com/daviddengcn/gosl
go install github.com/daviddengcn/gosl
```

#### Link to `/bin`
```bash
sudo ln -s $GOPATH/bin/gosl /bin/ls
```

#### Run a script
If a script starts with the bash interpreter line: `#!/bin/gosl`. You can run it like this
```bash
chmod a+x example.go
./example.go [params...]
```

Or you can explictly call `gosl` to run it:
```bash
gosl example.go  [params...]
```

Pre-imported Packages
---------------------
The following packages are pre-imported with `.`, i.e. you can directly use the methods exported by them. No complain of the compiler if you don't use them.

`fmt`, `os`, `strings`, `strconv`, `math`, `github.com/daviddengcn/gosl/builtin`

Frequently Used Builtin Functions
---------------------------------

Method | Description | Examples
--------|------------|-----------------------
`S`     | Convert anything to a `string` | `S(1234) == "123"`
`I`     | Convert anything to an `int`   | `I("1234") == 1234`
`BashEval` | Similar to bash backtick substitution. | `lsstr := BashEval("ls -l")`
`Exec`  | Execute an command with arguments  | `err, code := Exec("rm", "-rf" "tmp")`
`Bash`  | Execute a bash line           | `err, code := Bash("rm -rf tmp")`
`ScriptDir` | Returns the directory of the script | `file := ScriptDir() + "/" + fn`

More functions are defined in package [daviddengcn/gosl/builtin/](https://github.com/daviddengcn/gosl/tree/master/builtin) ([godoc](http://godoc.org/github.com/daviddengcn/gosl/builtin))

License
--------
Apache License V2
