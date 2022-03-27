package logger

import (
	"io"
	"os"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var (
	Output io.Writer = os.Stderr
	Level            = level.AllowWarn()

	once         sync.Once
	globalLogger log.Logger
)

func GlobalLogger() log.Logger {
	once.Do(func() {
		globalLogger = log.NewJSONLogger(Output)
		globalLogger = level.NewFilter(globalLogger, Level)
	})
	return globalLogger
}
