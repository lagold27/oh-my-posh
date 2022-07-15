package segments

import (
	"encoding/json"
	"oh-my-posh/environment"
	"oh-my-posh/properties"
	"path/filepath"
)

type Pulumi struct {
	props   properties.Properties
	env     environment.Environment
	Context string
}

type PulumiConfig struct {
	CurrentContext string `json:"current"`
}

func (pul *Pulumi) Template() string {
	return " {{ .Context }} "
}

func (pul *Pulumi) Init(props properties.Properties, env environment.Environment) {
	pul.props = props
	pul.env = env
}

func (pul *Pulumi) Enabled() bool {
	pulconfig := filepath.Join(pul.env.Home(), ".pulumi", "credentials.json")
	content := pul.env.FileContent(pulconfig)
	var config PulumiConfig
	err := json.Unmarshal([]byte(content), &config)
	if err != nil {
		pul.setError("Pulumi Error")
	}
	pul.Context = config.CurrentContext
	return true
}

func (pul *Pulumi) setError(message string) {
	if len(pul.Context) == 0 {
		pul.Context = message
	}
}
