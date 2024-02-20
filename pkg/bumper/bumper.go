package bumper

import (
	"fmt"

	"github.com/homeslice-ee/bumper/pkg/chooser"
	"github.com/homeslice-ee/bumper/pkg/registry"
)

// Bumper returns a URL for the latest version of an image.
type Bumper interface {
	Bump(string) (string, error)
}

// RegistryChoosingBumper chooses tags from a registry to bump images to their latest version.
type RegistryChoosingBumper struct {
	chooser  *chooser.Chooser
	registry registry.Registry
}

// New returns a new Bumper that chooses tags from a registry when bumping images to the latest version.
func New(reg registry.Registry, ch *chooser.Chooser) Bumper {
	return &RegistryChoosingBumper{
		chooser:  ch,
		registry: reg,
	}
}

// Bump returns a URL for the latest version of the given image.
func (b *RegistryChoosingBumper) Bump(image string) (string, error) {
	tag, err := b.chooser.Latest(image)
	if err != nil {
		return "", fmt.Errorf("couldn't choose for image %q: %w", image, err)
	}

	sha, err := b.registry.SHA256(image, tag)
	if err != nil {
		return "", fmt.Errorf("couldn't get sha for image/tag %q/%q: %w", image, tag, err)
	}

	url := b.registry.URL() + image + ":" + tag + "@" + sha
	return url, nil
}
