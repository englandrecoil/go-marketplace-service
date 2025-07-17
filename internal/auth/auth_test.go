package auth

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tests := map[string]struct {
		password string
		wantErr  error
	}{
		"first_successfully_hashed_password":  {password: "mysuperstrongpassword", wantErr: nil},
		"second_successfully_hashed_password": {password: "qwerty12345", wantErr: nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hash, err := HashPassword(tc.password)
			t.Logf("password: %s, hash: %s\n", tc.password, hash)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("%s: expected: %v, got: %v", name, tc.wantErr, err)
			}
		})
	}
}

func TestCheckHashPassword(t *testing.T) {
	tests := map[string]struct {
		password       string
		hashedPassword string
		wantErr        error
	}{
		"valid_password":   {password: "mysuperstrongpassword", hashedPassword: "$2a$10$EVavr/Uo6GZWIle3ZI1xuOlcbmXeGBfQWshvk2TOdkcOif2nkdCC6", wantErr: nil},
		"invalid_password": {password: "qwerty123456", hashedPassword: "$2a$10$CI7tmnsvcg0odFoUSztIoOQMytmuSApeWTJP4t1X.tR0AwQrcyaAC", wantErr: bcrypt.ErrMismatchedHashAndPassword},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := CheckPasswordHash(tc.password, tc.hashedPassword)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("%s: expected: %v, got: %v", name, tc.wantErr, err)
			}
		})
	}
}

/*
   auth_test.go:20: password: mysuperstrongpassword, hash: $2a$10$EVavr/Uo6GZWIle3ZI1xuOlcbmXeGBfQWshvk2TOdkcOif2nkdCC6
   auth_test.go:20: password: qwerty12345, hash: $2a$10$CI7tmnsvcg0odFoUSztIoOQMytmuSApeWTJP4t1X.tR0AwQrcyaAC
*/
