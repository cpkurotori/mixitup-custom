package logger

import (
	"fmt"
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
	globalLogger *Logger
)

type Logger struct {
	log.Logger
}

type writeOpt interface{}
type WithKey struct {
	Key    string
	Parser func(b []byte) interface{}
}

type WithKV []interface{}

type AlterLogger func(log.Logger) log.Logger

type writeFunc func(p []byte) (int, error)

func (f writeFunc) Write(p []byte) (int, error) {
	return f(p)
}

func (lgr *Logger) Writer(opts ...writeOpt) io.Writer {
	return writeFunc(func(p []byte) (int, error) {
		logger := lgr.Logger
		kv := []interface{}{}
		for _, o := range opts {
			switch opt := o.(type) {
			case WithKey:
				var val interface{} = string(p)
				if opt.Parser != nil {
					val = opt.Parser(p)
				}
				kv = append(kv, opt.Key, val)
			case WithKV:
				kv = append(kv, opt...)
			case AlterLogger:
				logger = opt(logger)
			default:
				panic(fmt.Sprintf("Unknown write opt: %T", opt))
			}
		}
		return len(p), logger.Log("msg", string(p))
	})
}

func GlobalLogger() *Logger {
	once.Do(func() {
		lgr := log.NewJSONLogger(Output)
		lgr = level.NewFilter(lgr, Level)
		globalLogger = &Logger{
			Logger: lgr,
		}
	})
	return globalLogger
}
