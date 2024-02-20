package main

import (
	"fmt"
	"os"

	"github.com/mikepartelow/bumper/pkg/bumper"
	"github.com/mikepartelow/bumper/pkg/chooser"
	"github.com/mikepartelow/bumper/pkg/registry"
	"github.com/mikepartelow/bumper/pkg/replacer"
)

func Usage(exitCode int) {
	p := os.Args[0]
	fmt.Println("Usage: " + p + " registry /path/to/text")
	fmt.Println("")
	fmt.Println("Bump all image refs prefixed with `registry` found in `/path/to/text` and prints the modified `/path/to/text`")
	fmt.Println("Considers only tags in `registry` prefixed with `.main`.")
	fmt.Println("")
	fmt.Println("Example:")
	fmt.Println(" " + p + " ghcr.io/mikepartelow/ Pulumi.prod.yaml")

	os.Exit(exitCode)
}

func main() {
	if len(os.Args) != 3 {
		Usage(1)
	}
	registryname, filename := os.Args[1], os.Args[2]

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reg := registry.New(registryname)
	ch := chooser.New(reg, chooser.MainSelector)
	b := bumper.New(reg, ch)

	rp := replacer.New(registryname, b)

	err = rp.Replace(os.Stdout, file)
	if err != nil {
		panic(err)
	}
}
