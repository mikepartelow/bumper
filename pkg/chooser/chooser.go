package chooser

import (
	"fmt"
	"slices"
	"strings"

	"github.com/mikepartelow/bumper/pkg/registry"
)

// TagSelector returns true if tag should be considered when choosing.
type TagSelector func(tag string) bool

// MainSelector returns true if tag starts with `.main`.
func MainSelector(tag string) bool {
	return strings.HasPrefix(tag, ".main")
}

// Chooser chooses tags from a registry.
type Chooser struct {
	registry registry.Registry
	selector TagSelector
}

// New returns a new Chooser that chooses tags from registry. If selector is nil, all tags are considered when choosing the latest.
// If selector is not nil, only tags that cause selector to return true are considrered when choosing the latest.
func New(registry registry.Registry, selector TagSelector) *Chooser {
	return &Chooser{
		registry: registry,
		selector: selector,
	}
}

// Latest returns a URI for the latest version of the given image.
func (ch *Chooser) Latest(image string) (string, error) {
	tags, err := ch.registry.Tags(image)
	if err != nil {
		return "", fmt.Errorf("couldn't fetch tags for image %q: %w", image, err)
	}

	slices.Sort(tags)

	for i := len(tags) - 1; i >= 0; i-- {
		if ch.selector == nil || ch.selector(tags[i]) {
			return tags[i], nil
		}
	}

	return "", registry.ErrNotFound
}
