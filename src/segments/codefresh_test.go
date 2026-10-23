package segments

import (
	"path/filepath"
	"testing"

	"github.com/jandedobbeleer/oh-my-posh/src/runtime/mock"
	"github.com/jandedobbeleer/oh-my-posh/src/segments/options"
	"github.com/stretchr/testify/assert"
)

func TestCodefreshSegment(t *testing.T) {
	cases := []struct {
		Case    string
		Content string
		Context string
	}{
		{Case: "production", Content: "current-context: prod", Context: "prod"},
		{Case: "development", Content: "contexts: {}\ncurrent-context: dev", Context: "dev"},
		{Case: "absent"},
		{Case: "empty", Content: ""},
		{Case: "whitespace content", Content: " \n\t"},
		{Case: "malformed", Content: "current-context: ["},
		{Case: "missing field", Content: "contexts: {}"},
		{Case: "empty context", Content: "current-context: ''"},
		{Case: "whitespace context", Content: "current-context: '  '"},
		{Case: "null context", Content: "current-context: null"},
		{Case: "wrong field type", Content: "current-context: [prod]"},
		{Case: "numeric context", Content: "current-context: 123"},
		{Case: "boolean context", Content: "current-context: true"},
	}

	for _, tc := range cases {
		t.Run(tc.Case, func(t *testing.T) {
			env := new(mock.Environment)
			env.On("Home").Return("testhome")
			env.On("FileContent", filepath.Join("testhome", ".cfconfig")).Return(tc.Content)
			cf := &Codefresh{Context: "stale"}
			cf.Init(options.Map{}, env)

			assert.Equal(t, tc.Context != "", cf.Enabled())
			assert.Equal(t, tc.Context, cf.Context)
			if tc.Context != "" {
				assert.Equal(t, " "+tc.Context+" ", renderTemplateNoTrimSpace(env, cf.Template(), cf))
				assert.Equal(t, tc.Context, renderTemplate(env, "{{ .Context }}", cf))
			}
			env.AssertExpectations(t)
		})
	}
}
