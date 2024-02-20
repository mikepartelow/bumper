package replacer

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mikepartelow/bumper/pkg/bumper"
)

// Replacer replaces image URLs in text with their latest version.
type Replacer struct {
	bumper       bumper.Bumper
	rxpImageName *regexp.Regexp
}

// New returns a new Replacer that bumps images to versions found in registryURL.
func New(registryURL string, bumper bumper.Bumper) *Replacer {
	return &Replacer{
		bumper:       bumper,
		rxpImageName: regexp.MustCompile(registryURL + `[^\s]+`),
	}
}

// Replace replaces image URLs in r with their latest version and writes the result to w.
func (rp *Replacer) Replace(w io.Writer, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		matches := rp.rxpImageName.FindAll([]byte(line), -1)
		if matches == nil {
			_, _ = w.Write([]byte(line + "\n"))
			continue
		}

		for _, match := range matches {
			imageRef := string(match)
			image, err := imageName(imageRef)
			if err != nil {
				return fmt.Errorf("error parsing image name %q: %w", imageRef, err)
			}

			bumpedImage, err := rp.bumper.Bump(image)
			if err != nil {
				return fmt.Errorf("error bumping image %q: %w", image, err)
			}

			line = strings.ReplaceAll(line, imageRef, bumpedImage)
		}

		_, _ = w.Write([]byte(line + "\n"))
	}

	return nil
}

func imageName(imageRef string) (string, error) {
	rxp := regexp.MustCompile(`.*?/([^/:@]+)[^/]*$`)
	matches := rxp.FindAllStringSubmatch(imageRef, -1)
	if len(matches) != 1 {
		return "", fmt.Errorf("could not parse imageRef: %q", imageRef)
	}

	return matches[0][1], nil
}
