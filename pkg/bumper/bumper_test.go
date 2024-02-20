package bumper_test

import (
	"slices"
	"testing"

	"github.com/homeslice-ee/bumper/pkg/bumper"
	"github.com/homeslice-ee/bumper/pkg/chooser"
	"github.com/homeslice-ee/bumper/pkg/registry"
	"github.com/stretchr/testify/assert"
)

type imageMap map[string]struct {
	tags []string
	shas []string
}

type stubRegistry struct {
	url    string
	images imageMap
}

func TestBumper(t *testing.T) {
	testCases := []struct {
		registry *stubRegistry
		image    string
		wantErr  error
		wantURI  string
	}{
		{
			registry: &stubRegistry{
				url: "blammo.io/path/to/awesomesauce/",
				images: imageMap{
					"spam": {
						tags: []string{"222", "000", "111"},
						shas: []string{"shafake:222", "shafake:000", "shafake:111"},
					},
				},
			},
			image:   "spam",
			wantURI: "blammo.io/path/to/awesomesauce/spam:222@shafake:222",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.image, func(t *testing.T) {
			ch := chooser.New(tC.registry, nil)

			b := bumper.New(tC.registry, ch)
			got, err := b.Bump(tC.image)
			assert.ErrorIs(t, tC.wantErr, err)
			assert.Equal(t, tC.wantURI, got)
		})
	}
}

func (r *stubRegistry) SHA256(image, tag string) (string, error) {
	info, ok := r.images[image]
	if !ok {
		return "", registry.ErrNotFound
	}

	for i, gotTag := range info.tags {
		if gotTag == tag {
			return info.shas[i], nil
		}
	}
	return "", registry.ErrNotFound
}

func (r *stubRegistry) Tags(image string) ([]string, error) {
	info, ok := r.images[image]
	if !ok {
		return nil, registry.ErrNotFound
	}

	return slices.Clone(info.tags), nil
}

func (r *stubRegistry) URL() string {
	return r.url
}
