package factory

import (
	"io"
	"os"
)

var writers = make(map[string]io.Writer)

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
	fWriters := make([]io.Writer, 0)
	for _, appender := range appenders {
		var writer io.Writer
		if "file" == appender.Type {
			writerConfig := toFileWriterConfig(appender)
			if writers[writerConfig.LogFilePath] == nil {
				w := newLumberjackWriter(writerConfig)
				writers[writerConfig.LogFilePath] = w
			}
			writer = writers[writerConfig.LogFilePath]
		} else if "stdout" == appender.Type {
			if writers["stdout"] == nil {
				w := os.Stdout
				writers["stdout"] = w
			}
			writer = writers["stdout"]
		} else if "stderr" == appender.Type {
			if writers["stderr"] == nil {
				w := os.Stderr
				writers["stderr"] = w
			}
			writer = writers["stderr"]
		}
		fWriters = append(fWriters, writer)
	}
	merged := mergeWriter(fWriters...)
	writers[name] = merged
	return merged
}
