package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendSlash(t *testing.T) {
	r := New("spam").(*OCIRegistry)
	assert.Equal(t, "spam/", r.url)

	r = New("spam/").(*OCIRegistry)
	assert.Equal(t, "spam/", r.url)
}
