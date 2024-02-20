package replacer_test

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/mikepartelow/bumper/pkg/registry"
	"github.com/mikepartelow/bumper/pkg/replacer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubBumper struct {
	bumps map[string]string
}

func TestReplacer(t *testing.T) {
	testCases := []struct {
		filename    string
		registryurl string
		bumps       map[string]string
	}{
		{
			filename:    "mild.txt",
			registryurl: "spamcr.io/eggs/",
			bumps: map[string]string{
				"spam": "spamcr.io/eggs/spam:v1.2.3@sha256:ramalamadingdong123",
			},
		},
		{
			filename:    "medium.txt",
			registryurl: "spamcr.io/eggs/",
			bumps: map[string]string{
				"spam": "spamcr.io/eggs/spam:v1.2.3@sha256:ramalamadingdong123",
				"eggs": "spamcr.io/eggs/eggs:v3.4.6@sha256:ringadingdong789",
			},
		},
		{
			filename:    "spicy.txt",
			registryurl: "spamcr.io/eggs/",
			bumps: map[string]string{
				"spam": "spamcr.io/eggs/spam:v1.2.3@sha256:ramalamadingdong123",
				"eggs": "spamcr.io/eggs/eggs:v3.4.6@sha256:ringadingdong789",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.filename, func(t *testing.T) {
			file, err := os.Open(path.Join("testdata", tC.filename))
			require.NoError(t, err)
			defer file.Close()

			want, err := os.ReadFile(path.Join("testdata", tC.filename+".want"))
			require.NoError(t, err)

			b := stubBumper{bumps: tC.bumps}
			r := replacer.New(tC.registryurl, &b)

			var got bytes.Buffer
			err = r.Replace(&got, file)
			assert.NoError(t, err)
			assert.Equal(t, string(want), got.String())
		})
	}
}

func (b *stubBumper) Bump(instr string) (string, error) {
	outstr, ok := b.bumps[instr]
	if !ok {
		return "", registry.ErrNotFound
	}

	return outstr, nil
}
