package handlers

import (
	"errors"
	"testing"
)

func TestValidateImage(t *testing.T) {
	tests := map[string]struct {
		url     string
		wantErr error
	}{
		"valid_image_url": {url: "https://avatars.githubusercontent.com/u/93984891?v=4", wantErr: nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := validateImage(tc.url)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("%s: expected: %v, got: %v", name, tc.wantErr, err)
			}
		})
	}
}
