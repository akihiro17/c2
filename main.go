package main

import (
	"bytes"
	"c2/compiler"
	_ "fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
