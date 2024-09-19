package factory

import (
	"io"
	"os"
	"strings"
	"sync"
)

var writers = make(map[string]io.Writer)
var writersLk = &sync.Mutex{}

type mergedWriter struct {
	delegates []io.Writer
}

func (m *mergedWriter) Write(p []byte) (int, error) {
	var wrote = 0
	var err error
	for _, d := range m.delegates {
		wrote, err = d.Write(p)
	}
	return wrote, err
}

func mergeWriter(w ...io.Writer) io.Writer {
	merged := mergedWriter{
		delegates: make([]io.Writer, 0),
	}
	for _, e := range w {
		merged.delegates = append(merged.delegates, e)
	}
	return &merged
}

func writer(name string, appenders []AppenderConfig) io.Writer {
	writersLk.Lock()
	defer writersLk.Unlock()
	appenderWriters := make([]io.Writer, 0)
	for _, appender := range appenders {
		var wr io.Writer
		if "file" == strings.ToLower(appender.Type) {
			writerConfig := toFileWriterConfig(appender)
			fileWriter, exists := writers[writerConfig.LogFilePath]
			if !exists || fileWriter == nil {
				fileWriter = newLumberjackWriter(writerConfig)
				writers[writerConfig.LogFilePath] = fileWriter
			}
			wr = fileWriter
		} else if "stdout" == strings.ToLower(appender.Type) {
			if writers["stdout"] == nil {
				stdoutWriter := os.Stdout
				writers["stdout"] = stdoutWriter
			}
			wr = writers["stdout"]
		} else if "stderr" == strings.ToLower(appender.Type) {
			if writers["stderr"] == nil {
				stderrWriter := os.Stderr
				writers["stderr"] = stderrWriter
			}
			wr = writers["stderr"]
		}
		appenderWriters = append(appenderWriters, wr)
	}
	merged := mergeWriter(appenderWriters...)
	writers[name] = merged
	return merged
}
