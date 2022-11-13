package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var logger = initZerolog("app")

func Log() *zerolog.Logger {
	return &logger
}

func initZerolog(logType string) zerolog.Logger {
	var file = initLog(logType)
	var t = zerolog.New(zerolog.ConsoleWriter{
		Out:        io.MultiWriter(os.Stdout, file),
		NoColor:    true,
		TimeFormat: time.RFC3339,
	}).With().Caller().Timestamp().Logger().Level(zerolog.GlobalLevel())
	return t
}

func initLog(logType string) *os.File {
	var f *os.File
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "data/log/"
	}
	path := logDir + logType + ".log"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		f, _ = os.Create(path)
	} else {
		f, _ = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	}
	return f
}
