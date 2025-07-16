package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"unicode"

	"github.com/englandrecoil/go-marketplace-service/internal/dto"
	"github.com/gin-gonic/gin"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const minEntropyBits = 60
const minLoginLength = 5
const maxLoginLength = 32

type ApiConfig struct {
	Conn *sql.DB
}

func (cfg *ApiConfig) HandlerRegister(c *gin.Context) {
	creds := dto.RegisterRequest{}
	if err := c.Bind(&creds); err != nil {
		dto.ResponseWithError(c, http.StatusBadRequest, "invalid request body format", err)
		return
	}

	// TODO: validate login
	if err := validateLogin(creds.Login); err != nil {
		dto.ResponseWithError(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	// validate password
	if err := passwordvalidator.Validate(creds.Password, minEntropyBits); err != nil {
		dto.ResponseWithError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// TODO: hash password

	// TODO: store credentials in db

}

func validateLogin(login string) error {
	if len(login) < minLoginLength || len(login) > maxLoginLength {
		return errors.New("invalid length of login")
	}
	if !unicode.IsLetter(rune(login[0])) {
		return errors.New("invalid format of login")
	}
	return nil
}
