package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlatformContextSegments(t *testing.T) {
	out := string(renderTheme(t, "../../themes/platform.omp.json", fixturePath))
	assert.Contains(t, out, "contoso-prod")
	assert.Contains(t, out, "\ue710 dev ")
	assert.NotContains(t, out, "<no value>")
	assert.NotContains(t, out, "unable to map writer")
}
