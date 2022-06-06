package factory

type LoggingConfig struct {
	Factory       string            `yaml:"factory"`
	RootName      string            `yaml:"root-name"`
	RootLevel     string            `yaml:"root-level"`
	PackageLevels map[string]string `yaml:"package-levels"`
	Formatter     string            `yaml:"formatter"`
	Appenders     []AppenderConfig  `yaml:"appenders"`
}

type AppenderConfig struct {
	Type    string            `yaml:"type"` // stdout | file | kafka ...
	Options map[string]string `yaml:"options"`
}
