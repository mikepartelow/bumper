package main

import "github.com/mikepartelow/bumper/pkg/logging"

func main() {
	logger := logging.Init()
	logger.Warn("Start")
}
