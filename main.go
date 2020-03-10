package main

import (
	"log"
	"os"

	"github.com/jfrog/kubenab/cmd"
	_log "github.com/sirupsen/logrus"
)

func main() {
	// print Version Informations
	log.Printf("Starting kubenab version %s - %s - %s", version, date, commit)

	// initialize logger
	initLogger()

	cmd.Execute()

	// check if all required Flags are set and in a correct Format
	// checkArguments()
}

// initLogger initializes the logrus logger – since it 'pkg/log' wrapps its
// functions.
func initLogger() {
	// Log as JSON instead of the default ASCII formatter.
	_log.SetFormatter(&_log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	_log.SetOutput(os.Stdout)

	// Only log the info severity or above (default)
	_log.SetLevel(_log.InfoLevel)
}
