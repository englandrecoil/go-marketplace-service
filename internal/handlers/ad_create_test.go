package handlers

import (
	"testing"
)

func TestValidateImage(t *testing.T) {
	tests := map[string]struct {
		url         string
		containsErr bool
	}{
		"valid_image_url":   {url: "https://avatars.githubusercontent.com/u/93984891?v=4", containsErr: false},
		"invalid_image_url": {url: "https://broken-image-url", containsErr: true},
		"not_an_image_url":  {url: "https://github.com/englandrecoil", containsErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := validateImage(tc.url)
			if err != nil && tc.containsErr == false {
				t.Fatalf("%s: expected no error, got: %v", name, err)
			}
			if err == nil && tc.containsErr == true {
				t.Fatalf("%s: error expected, got nothing", name)
			}
		})
	}
}
