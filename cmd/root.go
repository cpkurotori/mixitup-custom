package cmd

import (
	"fmt"
	"mixitup-custom/cmd/usercount"
	"mixitup-custom/internal/logger"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/cobra"
)

var (
	logFile  string
	logLevel string
)

var rootCmd = &cobra.Command{
	Use: "mixitup-custom",
}

func init() {
	rootCmd.AddCommand(usercount.UserCountCmd)

	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "-", "log file, `-` for StdErr")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "WARN", "log level [DEBUG|INFO|WARN|ERROR]")
}

func Execute() error {
	rootCmd.ParseFlags(os.Args[1:])
	if logFile != "" && logFile != "-" {
		var err error
		logWriter, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		logger.Output = log.NewSyncWriter(logWriter)
		switch logLevel {
		case "ERROR":
			logger.Level = level.AllowError()
		case "WARN":
			logger.Level = level.AllowWarn()
		case "INFO":
			logger.Level = level.AllowInfo()
		case "DEBUG":
			logger.Level = level.AllowDebug()
		default:
			return fmt.Errorf("Unrecognized log-level: %s", logLevel)
		}
		defer logWriter.Close()
		rootCmd.SetErr(logger.Output)
		rootCmd.SetOut(logger.Output)
	}
	return rootCmd.Execute()
}
