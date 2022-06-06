package factory

import (
	"strconv"
	"strings"
	"time"
)

type fileAppenderOptions struct {
	LogFileDir     string `yaml:"log-file-dir"`
	LogFileName    string `yaml:"log-file-name"`
	MaxFileSize    int    `yaml:"max-file-size"`
	MaxFileBackups int    `yaml:"max-file-backups"`
	MaxFileAge     string `yaml:"max-file-age"` // 秒
	LocalTime      bool   `yaml:"local-time"`
	Compress       bool   `yaml:"compress"`
}

var fileAppenderOptionKeyLogFileDir = "log-file-dir"
var fileAppenderOptionKeyLogFileName = "log-file-name"
var fileAppenderOptionKeyMaxFileSize = "max-file-size"
var fileAppenderOptionKeyMaxFileBackups = "max-file-backups"
var fileAppenderOptionKeyMaxFileAge = "max-file-age"
var fileAppenderOptionKeyLocalTime = "local-time"
var fileAppenderOptionKeyCompress = "compress"

type fileWriterConfig struct {
	LogFilePath    string
	MaxFileSize    int
	MaxFileBackups int
	MaxFileAge     time.Duration
	LocalTime      bool
	Compress       bool
}

func toFileWriterConfig(appender AppenderConfig) *fileWriterConfig {
	vLogFileDir := appender.Options[fileAppenderOptionKeyLogFileDir]
	vLogFileName := appender.Options[fileAppenderOptionKeyLogFileName]
	vMaxFileSize, _ := strconv.Atoi(appender.Options[fileAppenderOptionKeyMaxFileSize])
	vMaxFileBackups, _ := strconv.Atoi(appender.Options[fileAppenderOptionKeyMaxFileBackups])
	vMaxFileAge, _ := appender.Options[fileAppenderOptionKeyMaxFileAge]
	vLocalTime, _ := strconv.ParseBool(appender.Options[fileAppenderOptionKeyLocalTime])
	vCompress, _ := strconv.ParseBool(appender.Options[fileAppenderOptionKeyCompress])
	options := &fileAppenderOptions{
		LogFileDir:     vLogFileDir,
		LogFileName:    vLogFileName,
		MaxFileSize:    vMaxFileSize,
		MaxFileBackups: vMaxFileBackups,
		MaxFileAge:     vMaxFileAge,
		LocalTime:      vLocalTime,
		Compress:       vCompress,
	}
	logFileDir := strings.TrimSpace(options.LogFileDir)
	if len(logFileDir) <= 0 {
		logFileDir = "./logs"
	}
	logFileName := strings.TrimSpace(options.LogFileName)
	if len(logFileName) <= 0 {
		logFileName = "./application.log"
	}
	logFilePath := logFileDir + SLASH + logFileName
	maxFileAge, _ := time.ParseDuration(options.MaxFileAge)

	return &fileWriterConfig{
		LogFilePath:    logFilePath,
		MaxFileSize:    options.MaxFileSize,
		MaxFileBackups: options.MaxFileBackups,
		MaxFileAge:     maxFileAge,
		LocalTime:      options.LocalTime,
		Compress:       options.Compress,
	}
}
