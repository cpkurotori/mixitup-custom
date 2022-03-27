package main

import (
	"flag"
	"fmt"
	"mixitup-custom/cmd"
	"mixitup-custom/internal/logger"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var (
	logFile  = flag.String("log-file", "-", "file to log to; omit or use `-` to log to stderr")
	logLevel = flag.String("log-level", "INFO", "log level (DEBUG|INFO|WARN|ERROR)")
)

func main() {
	flag.Parse()
	cmd.AddFlagSet(flag.CommandLine)

	switch *logLevel {
	case "ERROR":
		logger.Level = level.AllowError()
	case "WARN", "":
		logger.Level = level.AllowWarn()
	case "INFO":
		logger.Level = level.AllowInfo()
	case "DEBUG":
		logger.Level = level.AllowDebug()
	default:
		panic(fmt.Errorf("Unrecognized log-level: %s", *logLevel))
	}

	if *logFile != "" && *logFile != "-" {
		logWriter, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer logWriter.Close()
		logger.Output = log.NewSyncWriter(logWriter)
	}
	if err := cmd.Execute(); err != nil {
		fmt.Fprint(os.Stdout, "ERROR")
	}
}
