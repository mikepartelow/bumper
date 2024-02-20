package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mikepartelow/bumper/pkg/bumper"
	"github.com/mikepartelow/bumper/pkg/chooser"
	"github.com/mikepartelow/bumper/pkg/registry"
)

func Usage(exitCode int) {
	p := os.Args[0]
	fmt.Println("Usage: " + p + " image")
	fmt.Println("")
	fmt.Println("Print the latest pinned url for `image` considering all available tags prefixed with `.main`")
	fmt.Println("")
	fmt.Println("Example:")
	fmt.Println(" " + p + " ghcr.io/mikepartelow/bumper")

	os.Exit(exitCode)
}

func main() {
	if len(os.Args) != 2 {
		Usage(1)
	}
	imageRef := os.Args[1]

	reg, image := splitImageRef(imageRef)

	r := registry.New(reg)
	ch := chooser.New(r, chooser.MainSelector)
	b := bumper.New(r, ch)

	fmt.Println(must(b.Bump(image)))
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func splitImageRef(imageRef string) (string, string) {
	parts := strings.Split(imageRef, "/")
	reg := strings.Join(parts[0:len(parts)-1], "/")
	image := parts[len(parts)-1]
	return reg, image
}
