package handlers

import (
	"errors"
	"testing"
)

func TestValidateLogin(t *testing.T) {
	tests := map[string]struct {
		login   string
		wantErr error
	}{
		"valid_login":                 {login: "spiderman125", wantErr: nil},
		"login_too_short":             {login: "cool", wantErr: ErrInvalidLoginLength},
		"login_too_big":               {login: "mysupercoolbiglogthatdoesntfitintohere", wantErr: ErrInvalidLoginLength},
		"login_starts_with_no_letter": {login: "_dontstealmyloginplease", wantErr: ErrInvalidLoginFormat},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := validateLogin(tc.login)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("%s: expected: %v, got: %v", name, tc.wantErr, err)
			}
		})
	}

}
