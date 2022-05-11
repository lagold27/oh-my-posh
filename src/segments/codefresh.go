package segments

import (
	"gopkg.in/yaml.v3"
	"oh-my-posh/environment"
	"oh-my-posh/properties"
	"path/filepath"
)

type Codefresh struct {
	props   properties.Properties
	env     environment.Environment
	Context string
}

type CodefreshConfig struct {
	CurrentContext string `yaml:"current-context"`
}

func (cf *Codefresh) Template() string {
	return " {{ .Context }} "
}

func (cf *Codefresh) Init(props properties.Properties, env environment.Environment) {
	cf.props = props
	cf.env = env
}

func (cf *Codefresh) Enabled() bool {
	cfconfig := filepath.Join(cf.env.Home(), ".cfconfig")
	content := cf.env.FileContent(cfconfig)
	var config CodefreshConfig
	err := yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		cf.setError("Codefresh Error")
		return false
	}
	cf.Context = config.CurrentContext
	return true
}

func (cf *Codefresh) setError(message string) {
	if len(cf.Context) == 0 {
		cf.Context = message
	}
}
