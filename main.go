package main

import (
	"os"
	"os/exec"
	"bytes"
	"io/ioutil"
	_"fmt"
	"strings"
	"./compiler"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic("could not read file")
	}

	code := string(data)
	compiler := compiler.New(code)
	out := new(bytes.Buffer)
	compiler.Compile(out)

	ioutil.WriteFile("program.s", out.Bytes(), os.ModePerm)

	executable := os.Args[1][0:strings.Index(os.Args[1], ".c")]
	exec.Command("gcc", "program.s", "-o", executable).Run()
}
