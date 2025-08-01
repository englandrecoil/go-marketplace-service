package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"time"
	"unicode"

	"github.com/englandrecoil/go-marketplace-service/internal/auth"
	"github.com/englandrecoil/go-marketplace-service/internal/constants"
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

type ApiConfig struct {
	Conn   *sql.DB
	DB     *database.Queries
	Secret string
}

// HandlerRegister godoc
//
//	@Summary		Зарегистрировать нового пользователя
//	@Description	Создаёт нового пользователя с заданным логином и паролем
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		dto.CredentialsRequest	true	"Данные пользователя для входа"
//	@Success		201			{object}	dto.RegisterResponse	"Успешная регистрация"
//	@Failure		400			{object}	dto.ErrorResponse		"Неверный формат запроса, логин уже используется или пароль слишком слабый"
//	@Failure		500			{object}	dto.ErrorResponse		"Внутренняя ошибка сервера"
//	@Router			/api/reg [post]
func (cfg *ApiConfig) HandlerRegister(c *gin.Context) {
	inputCredentials := dto.CredentialsRequest{}
	if err := c.BindJSON(&inputCredentials); err != nil {
		dto.ResponseWithError(c, http.StatusBadRequest, "invalid request body format", err)
		return
	}

	// validate login
	if err := validateLogin(inputCredentials.Login); err != nil {
		dto.ResponseWithError(c, http.StatusBadRequest, err.Error(), err)
		return
	}
	// validate password
	if err := passwordvalidator.Validate(inputCredentials.Password, constants.MinEntropyBits); err != nil {
		dto.ResponseWithError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(inputCredentials.Password)
	if err != nil {
		dto.ResponseWithError(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	// store credentials in db
	user, err := cfg.DB.CreateUser(
		c.Request.Context(),
		database.CreateUserParams{
			Login:          inputCredentials.Login,
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
	if len(login) < constants.MinLoginLength || len(login) > constants.MaxLoginLength {
		return ErrInvalidLoginLength
	}
	if !unicode.IsLetter(rune(login[0])) {
		return ErrInvalidLoginFormat
	}
	return nil
}
