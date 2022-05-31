package factory

import (
	"io"
	"os"
	"strings"
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

func writer(name string, outputs []string) io.Writer {
	fWriters := make([]io.Writer, 0)
	for _, p := range outputs {
		var writer io.Writer
		if strings.ToLower(p) == "stdout" {
			if writers["stdout"] == nil {
				w := os.Stdout
				writers["stdout"] = w
			}
			writer = writers["stdout"]
		} else if strings.ToLower(p) == "stderr" {
			if writers["stderr"] == nil {
				w := os.Stderr
				writers["stderr"] = w
			}
			writer = writers["stderr"]
		} else {
			if writers[p] == nil {
				w := newLumberjackWriter(p)
				writers[p] = w
			}
			writer = writers[p]
		}
		fWriters = append(fWriters, writer)
	}
	merged := mergeWriter(fWriters...)
	writers[name] = merged
	return merged
}
