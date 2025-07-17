package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"time"
	"unicode"

	"github.com/englandrecoil/go-marketplace-service/internal/auth"
	"github.com/englandrecoil/go-marketplace-service/internal/database"
	"github.com/englandrecoil/go-marketplace-service/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

var (
	ErrInvalidLoginLength = errors.New("invalid length of login")
	ErrInvalidLoginFormat = errors.New("invalid format of login")
)

const minEntropyBits = 60
const minLoginLength = 5
const maxLoginLength = 32

type ApiConfig struct {
	Conn *sql.DB
	DB   *database.Queries
}

func (cfg *ApiConfig) HandlerRegister(c *gin.Context) {
	credentials := dto.RegisterRequest{}
	if err := c.Bind(&credentials); err != nil {
		dto.ResponseWithError(c, http.StatusBadRequest, "invalid request body format", err)
		return
	}

	// validate login
	if err := validateLogin(credentials.Login); err != nil {
		dto.ResponseWithError(c, http.StatusBadRequest, err.Error(), err)
		return
	}
	// validate password
	if err := passwordvalidator.Validate(credentials.Password, minEntropyBits); err != nil {
		dto.ResponseWithError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(credentials.Password)
	if err != nil {
		dto.ResponseWithError(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	// TODO: store credentials in db
	user, err := cfg.DB.CreateUser(
		c.Request.Context(),
		database.CreateUserParams{
			Login:          credentials.Login,
			HashedPassword: hashedPassword,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		},
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				dto.ResponseWithError(c, http.StatusBadRequest, "this login's already in use", nil)
				return
			}
		}
		dto.ResponseWithError(c, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	c.JSON(
		http.StatusCreated,
		dto.RegisterResponse{
			ID:        user.ID,
			Login:     user.Login,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	)
}

func validateLogin(login string) error {
	if len(login) < minLoginLength || len(login) > maxLoginLength {
		return ErrInvalidLoginLength
	}
	if !unicode.IsLetter(rune(login[0])) {
		return ErrInvalidLoginFormat
	}
	return nil
}
