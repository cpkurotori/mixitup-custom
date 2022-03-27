package cmd

import (
	"flag"
	"mixitup-custom/cmd/usercount"
	"mixitup-custom/internal/logger"
	"os"

	"github.com/go-kit/log/level"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "mixitup-custom",
}

func init() {
	rootCmd.AddCommand(usercount.UserCountCmd)
}

func AddFlagSet(set *flag.FlagSet) {
	rootCmd.PersistentFlags().AddGoFlagSet(set)
}

func Execute() error {
	if logger.Output != nil && logger.Output != os.Stderr {
		rootCmd.SetErr(logger.GlobalLogger().Writer(logger.WithKey{Key: "log"}, logger.WithKV{"context", "root_err"}, logger.AlterLogger(level.Error)))
		rootCmd.SetOut(logger.GlobalLogger().Writer(logger.WithKey{Key: "log"}, logger.WithKV{"context", "root_out"}, logger.AlterLogger(level.Info)))
		usercount.UserCountCmd.SetErr(logger.GlobalLogger().Writer(logger.WithKey{Key: "log"}, logger.WithKV{"context", "user_count_err"}, logger.AlterLogger(level.Error)))
		usercount.UserCountCmd.SetOut(logger.GlobalLogger().Writer(logger.WithKey{Key: "log"}, logger.WithKV{"context", "user_count_out"}, logger.AlterLogger(level.Info)))
	}

	return rootCmd.Execute()
}
