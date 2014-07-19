gosl
====

This is an application that can make you write script with Go languages.

Example
-------

```go
#!/bin/gosl

import "encoding/json"

toJson := func(lines []string) string {
  res, _ := json.Marshal(lines)
  return string(res)
}

files, _, _ := BashWithStdout("ls -l /tmp/")

Println(toJson(Split(files, "\n")))
    
```

Installation
------------

#### Download and install the package
```bash
go get github.com/daviddengcn/gosl
go install github.com/daviddengcn/gosl
```

#### Link to `/bin`
```bash
sudo ln -s $GOPATH/bin/gosl /bin/ls
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
