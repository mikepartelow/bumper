package chooser_test

import (
	"strings"
	"testing"

	"github.com/mikepartelow/bumper/pkg/chooser"
	"github.com/mikepartelow/bumper/pkg/registry"
	"github.com/stretchr/testify/assert"
)

type stubRegistry struct {
	tags map[string][]string
}

func TestLast(t *testing.T) {
	testCases := []struct {
		image    string
		selector chooser.TagSelector
		tags     map[string][]string
		wantErr  error
		wantURI  string
	}{
		{
			image:   "bar",
			wantErr: registry.ErrNotFound,
		},
		{
			image: "foo",
			tags: map[string][]string{
				"foo": {
					"123",
					"000",
				},
			},
			wantURI: "123",
		},
		{
			image:    "foo",
			selector: func(tag string) bool { return strings.HasPrefix(tag, "eggs.") },
			tags: map[string][]string{
				"foo": {
					"spam.123",
					"eggs.000",
				},
			},
			wantURI: "eggs.000",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.image, func(t *testing.T) {
			reg := &stubRegistry{tags: tC.tags}
			ch := chooser.New(reg, tC.selector)

			got, err := ch.Latest(tC.image)
			assert.ErrorIs(t, err, tC.wantErr)
			assert.Equal(t, tC.wantURI, got)
		})
	}
}

func (r *stubRegistry) SHA256(string, string) (string, error) {
	panic("not supposed to be called here")
}

func (r *stubRegistry) Tags(image string) ([]string, error) {
	tags, ok := r.tags[image]
	if !ok {
		return nil, registry.ErrNotFound
	}
	return tags, nil
}

func (r *stubRegistry) URL() string {
	panic("not supposed to be called here")
}
