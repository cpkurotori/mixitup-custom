package usercount

import (
	"context"
	"errors"
	"fmt"
	internallog "mixitup-custom/internal/logger"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	counterFile string
	userID      string
	username    string
	add         int64
)

func init() {
	UserCountCmd.Flags().StringVar(&counterFile, "counter-file", "", "file to store counters in (required); will create this file if it does not exist")
	UserCountCmd.Flags().StringVar(&userID, "user-id", "", "user id to add counter to (required)")
	UserCountCmd.Flags().StringVar(&username, "user-name", "", "user name to attach to id")
	UserCountCmd.Flags().Int64Var(&add, "add", 1, "the amount to add")

	UserCountCmd.MarkFlagRequired("counter-file")
	UserCountCmd.MarkFlagRequired("user-id")
}

func initializeDB() error {
	if _, err := os.Stat(counterFile); os.IsNotExist(err) {
		file, err := os.Create(counterFile)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

type UserCounter struct {
	UserID  string `gorm:"primaryKey"`
	Name    string
	Counter int64
}

type dbLogger struct {
	*internallog.Logger
}

func (lgr *dbLogger) LogMode(lvl logger.LogLevel) logger.Interface {
	var l log.Logger
	switch lvl {
	case logger.Error:
		l = level.Error(lgr.Logger)
	case logger.Info:
		l = level.Info(lgr.Logger)
	case logger.Warn:
		l = level.Warn(lgr.Logger)
	default:
		l = log.NewNopLogger()
	}
	return &dbLogger{
		Logger: &internallog.Logger{
			Logger: l,
		},
	}
}

func (logger *dbLogger) Info(_ context.Context, msg string, args ...interface{}) {
	_ = level.Info(logger.Logger).Log("msg", fmt.Sprintf(msg, args...))
}

func (logger *dbLogger) Warn(_ context.Context, msg string, args ...interface{}) {
	_ = level.Warn(logger.Logger).Log("msg", fmt.Sprintf(msg, args...))
}

func (logger *dbLogger) Error(_ context.Context, msg string, args ...interface{}) {
	_ = level.Error(logger.Logger).Log("msg", fmt.Sprintf(msg, args...))
}

func (logger *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	_ = level.Debug(logger.Logger).Log("msg", "trace", "begin", begin, "sql", sql, "rows_affected", rows, "err", err)
}

var UserCountCmd = &cobra.Command{
	Use: "user-count",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		defer func() {
			if err != nil {
				_ = level.Error(internallog.GlobalLogger()).Log("msg", "error while handling user count", "err", err)
			}
		}()
		if err := initializeDB(); err != nil {
			return err
		}
		db, err := gorm.Open(sqlite.Open(counterFile), &gorm.Config{
			Logger: &dbLogger{
				Logger: internallog.GlobalLogger(),
			},
		})
		if err != nil {
			return err
		}

		if err := db.AutoMigrate(&UserCounter{}); err != nil {
			return err
		}

		var user UserCounter
		if err := db.First(&user, "user_id = ?", userID).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			user.UserID = userID
		}

		if username != "" {
			user.Name = username
		}

		user.Counter += add

		if err := db.Save(&user).Error; err != nil {
			return err
		}

		level.Info(internallog.GlobalLogger()).Log("msg", "user count complete", "count", user.Counter, "user_id", user.UserID, "username", username)
		fmt.Fprintf(os.Stdout, "%d", user.Counter)
		return nil
	},
}
