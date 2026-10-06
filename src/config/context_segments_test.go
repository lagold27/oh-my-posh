package config

import (
	"bytes"
	"encoding/gob"
	"path/filepath"
	"testing"

	"github.com/jandedobbeleer/oh-my-posh/src/cache"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime/mock"
	"github.com/jandedobbeleer/oh-my-posh/src/segments"
	"github.com/jandedobbeleer/oh-my-posh/src/segments/options"
	"github.com/jandedobbeleer/oh-my-posh/src/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextSegmentFactoryAndGob(t *testing.T) {
	originalCache := template.Cache
	t.Cleanup(func() { template.Cache = originalCache })
	template.Cache = &cache.Template{}

	cases := []struct {
		Writer  SegmentWriter
		Options options.Map
		Type    SegmentType
		File    string
		Content string
	}{
		{
			Type: CODEFRESH, Options: options.Map{},
			File: filepath.Join("testhome", ".cfconfig"), Content: "current-context: contoso-dev",
			Writer: &segments.Codefresh{},
		},
		{
			Type: PULUMI, Options: options.Map{segments.ContextOnly: true},
			File: filepath.Join("testhome", ".pulumi", "credentials.json"), Content: `{"current":"contoso-dev"}`,
			Writer: &segments.Pulumi{},
		},
	}

	for _, tc := range cases {
		t.Run(string(tc.Type), func(t *testing.T) {
			env := new(mock.Environment)
			env.On("Shell").Return("generic")
			env.On("Home").Return("testhome")
			if tc.Type == PULUMI {
				env.On("Getenv", "PULUMI_HOME").Return("")
			}
			env.On("FileContent", tc.File).Return(tc.Content)
			segment := &Segment{Type: tc.Type, Options: tc.Options}
			require.NoError(t, segment.MapSegmentWithWriter(env))
			writer := segment.Writer()
			assert.IsType(t, tc.Writer, writer)
			require.True(t, writer.Enabled())

			var buf bytes.Buffer
			require.NoError(t, gob.NewEncoder(&buf).Encode(&writer))
			var restored SegmentWriter
			require.NoError(t, gob.NewDecoder(&buf).Decode(&restored))
			assert.IsType(t, tc.Writer, restored)
			segment.writer = restored
			segment.Template = "{{ .Context }}"
			template.Init(env, nil, nil)
			assert.Equal(t, "contoso-dev", segment.string())
			env.AssertExpectations(t)
		})
	}
}
