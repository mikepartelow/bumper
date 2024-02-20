package main

import "github.com/mikepartelow/bumper/pkg/logging"

func main() {
	logger := logging.Init()
	logger.Warn("Start")

	// - set up env:
	//   - git repo to monitor
	//   - bumper registry
	//   - branch name template

	// - fetch git repo

	// - bump

	// - if diffs
	//   - branch
	//   - push

	// - exit
}
