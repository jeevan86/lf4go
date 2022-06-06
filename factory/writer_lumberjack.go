package factory

import (
	"fmt"
	"github.com/natefinch/lumberjack/v3"
	"io"
	"os"
	"time"
)

const defaultMaxFileSize = 100 * 1024 * 1024
const defaultMaxFileBackups = 20
const defaultMaxFileAge = time.Hour * 72
const defaultLocalTime = true
const defaultCompress = true

func defaultLumberjackWriter(filePath string) io.Writer {
	options := lumberjack.Options{
		MaxAge:     defaultMaxFileAge,
		MaxBackups: defaultMaxFileBackups,
		LocalTime:  defaultLocalTime,
		Compress:   defaultCompress,
	}
	writer, _ := lumberjack.NewRoller(filePath, defaultMaxFileSize, &options)
	return writer
}

func newLumberjackWriter(config *fileWriterConfig) io.Writer {
	options := lumberjack.Options{
		MaxAge:     config.MaxFileAge,
		MaxBackups: config.MaxFileBackups,
		LocalTime:  config.LocalTime,
		Compress:   config.Compress,
	}
	writer, err := lumberjack.NewRoller(config.LogFilePath, defaultMaxFileSize, &options)
	if err != nil {
		fmt.Println(fmt.Sprintf("Fatal! %s", err.Error()))
		os.Exit(-1)
	}
	return writer
}
