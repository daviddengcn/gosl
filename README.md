gosl
====

This is an application that can make you write script with Go languages.

It is *NOT* an interpreter but the pure Go. The preprocessor tranforms the script into a Go program, instantly compile and run. So it is almost same as the standard Go with the same efficiency.

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

Builtin Functions
-----------------

Method | Description
--------|------------------------------------
`Exec` | Execute an command with arguments 
`Bash` | Execute a bash line
`S`    | Convert anything to a `string`
`I`    | Convert anything to an `int`

License
--------
Apache License V2
