package handlers

import (
	"errors"
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

func TestValidateAdParams(t *testing.T) {
	tests := map[string]struct {
		title        string
		description  string
		imageAddress string
		price        int
		wantErr      error
	}{
		"invalid_length_title":       {title: "", description: "description", imageAddress: "https://iili.io/2g8pbwu.jpg", price: 55, wantErr: ErrInvalidLengthTitle},
		"invalid_length_description": {title: "myadtitle", description: "desc", imageAddress: "https://iili.io/2g8pbwu.jpg", price: 143, wantErr: ErrInvalidLengthDescription},
		"invalid_price":              {title: "myadtitle2", description: "super cool description", imageAddress: "https://iili.io/2g8pbwu.jpg", price: -100, wantErr: ErrInvalidPrice},
		"valid_ad_parameters":        {title: "Macbook Air 13 M1", description: "Super cool and brand new laptop (2020 lol)", imageAddress: "https://iili.io/2g8pbwu.jpg", price: 78900, wantErr: nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := validateAdParams(tc.title, tc.description, tc.imageAddress, tc.price)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("%s: expected: %v, got: %v", name, tc.wantErr, err)
			}
		})
	}
}
