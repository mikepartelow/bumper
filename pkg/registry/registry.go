package registry

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

var (
	ErrNotFound = errors.New("not found")
)

// Registry queries a container registry for information about images.
type Registry interface {
	SHA256(image, tag string) (string, error)
	Tags(image string) ([]string, error)
	URL() string
}

// OCIRegistry connects to remote OCI container registries.
type OCIRegistry struct {
	url string
}

// New returns a new Registry hosted at url.
//
// url is everything prior to the image name. For `ghcr.io/mikepartelow/bumper`, the Registry url is `ghcr.io/mikepartelow/`
func New(url string) Registry {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	return &OCIRegistry{
		url: url,
	}
}

// SHA256 returns an image SHA256 digest suitable for pinning the image.
//
// `docker pull ghcr.io/mikepartelow/bumper:main.v1.2.3@{{ sha256 ) }}`
func (r *OCIRegistry) SHA256(image, tag string) (string, error) {
	imageRef := r.url + image + ":" + tag

	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return "", fmt.Errorf("couldn't parse imageRef %q: %w", imageRef, err)
	}

	i, err := remote.Image(ref)
	if err != nil {
		return "", fmt.Errorf("couldn't get remote image %q: %w", imageRef, err)
	}

	d, err := i.Digest()
	if err != nil {
		return "", fmt.Errorf("couldn't get image digest %q: %w", imageRef, err)

	}

	return d.String(), nil
}

// Tags returns a list of container tags found on image.
func (r *OCIRegistry) Tags(image string) ([]string, error) {
	imageRef := r.url + image

	repo, err := name.NewRepository(imageRef)
	if err != nil {
		panic(err)
	}

	return remote.List(repo)
}

// URL returns the URL passed to New()
func (r *OCIRegistry) URL() string {
	return r.url
}
