package factory

import (
	"github.com/sirupsen/logrus"
	"strings"
)

var logrusFormatters = map[string]logrus.Formatter{
	"normal": &logrus.TextFormatter{
		ForceColors:               false,
		DisableColors:             true,
		ForceQuote:                false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           DTFormatNormal,
		DisableSorting:            true,
		SortingFunc:               nil,
		DisableLevelTruncation:    true,
		PadLevelText:              true,
		QuoteEmptyFields:          true,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	},
	"json": &logrus.JSONFormatter{
		DisableTimestamp:  false,
		TimestampFormat:   DTFormatNormal,
		DisableHTMLEscape: true,
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       false,
		DataKey:           "msg",
	},
}

func logrusFormatter(name string) logrus.Formatter {
	formatter := strings.ToLower(name)
	if formatter == "normal" {
		formatter = "normal"
	} else if formatter == "json" {
		formatter = "json"
	} else {
		formatter = "normal"
	}
	return logrusFormatters[strings.ToLower(formatter)]
}
