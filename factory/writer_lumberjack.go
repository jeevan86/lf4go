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

func newLumberjackWriter(filePath string) io.Writer {
	options := lumberjack.Options{
		MaxAge:     defaultMaxFileAge,
		MaxBackups: defaultMaxFileBackups,
		LocalTime:  true,
		Compress:   true,
	}
	writer, err := lumberjack.NewRoller(filePath, defaultMaxFileSize, &options)
	if err != nil {
		fmt.Println(fmt.Sprintf("Fatal! %s", err.Error()))
		os.Exit(-1)
	}
	return writer
}
