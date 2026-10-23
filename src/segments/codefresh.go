package segments

import (
	"path/filepath"
	"strings"

	yaml "go.yaml.in/yaml/v3"
)

type Codefresh struct {
	Base

	Context string
}

type CodefreshConfig struct {
	// A string target would coerce numeric and boolean YAML scalars.
	CurrentContext yaml.Node `yaml:"current-context"`
}

func (cf *Codefresh) Template() string {
	return " {{ .Context }} "
}

func (cf *Codefresh) Enabled() bool {
	cf.Context = ""
	cfconfig := filepath.Join(cf.env.Home(), ".cfconfig")
	content := cf.env.FileContent(cfconfig)
	var config CodefreshConfig
	err := yaml.Unmarshal([]byte(content), &config)
	if err != nil || config.CurrentContext.Tag != "!!str" {
		return false
	}

	if strings.TrimSpace(config.CurrentContext.Value) == "" {
		return false
	}

	cf.Context = config.CurrentContext.Value
	return true
}
