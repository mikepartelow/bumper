package main

import (
	"log/slog"
	"os"
	"path"

	"github.com/go-git/go-git/v5" // with go modules enabled (GO111MODULE=on or outside GOPATH)

	"github.com/mikepartelow/bumper/pkg/bumper"
	"github.com/mikepartelow/bumper/pkg/chooser"
	"github.com/mikepartelow/bumper/pkg/logging"
	"github.com/mikepartelow/bumper/pkg/registry"
	"github.com/mikepartelow/bumper/pkg/replacer"
)

func main() {
	logger := logging.Init()
	logger.Warn("Start")

	// Git repo containing a file to bump, like git@github.com/mikepartelow/bumper
	gitRepo := mustGetenv("GIT_REPO", logger)
	// the path in the Git repo of the file to bump, like /myapp/Pulumi.dev.yaml
	bumpFilename := mustGetenv("BUMP_FILENAME", logger)
	// a container registry containing versions to bump, like ghcr.io/mikepartelow/
	containerRegistry := mustGetenv("CONTAINER_REGISTRY", logger)

	clonePath := path.Join(os.TempDir(), "bumper.clone")
	_, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		URL: gitRepo,
	})
	if err != nil {
		logger.Error("error cloning git repo", "repo", gitRepo, "err", err)
		panic(err)
	}

	file, err := os.Open(path.Join(clonePath, bumpFilename))
	if err != nil {
		logger.Error("couldn't open bump file", "bump filename", bumpFilename, "repo", gitRepo, "err", err)
		panic(err)
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp("", "bumper")
	if err != nil {
		logger.Error("couldn't create temp file", "repo", gitRepo, "err", err)
		panic(err)
	}
	defer tmpFile.Close()

	reg := registry.New(containerRegistry)
	b := bumper.New(reg, chooser.New(reg, chooser.MainSelector))

	repl := replacer.New(containerRegistry, b)

	err = repl.Replace(tmpFile, file)
	if err != nil {
		logger.Error("couldn't bump file", "bump filename", bumpFilename, "repo", gitRepo, "err", err)
		panic(err)
	}

	// - if diffs
	//   - branch
	//   - push

	// - exit
}

func mustGetenv(name string, logger *slog.Logger) string {
	val := os.Getenv(name)
	if val == "" {
		logger.Error("couldn't get required environment variable", "name", name)
		panic("couldn't get required environment variable: " + name)
	}
	return val
}
