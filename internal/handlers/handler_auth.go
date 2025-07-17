package handlers

import (
	"net/http"
	"time"

	"github.com/englandrecoil/go-marketplace-service/internal/auth"
	"github.com/englandrecoil/go-marketplace-service/internal/dto"
	"github.com/gin-gonic/gin"
)

const tokenExpirationTime time.Duration = time.Minute * 15

func (cfg *ApiConfig) HandlerAuth(c *gin.Context) {
	inputCredentials := dto.CredentialsRequest{}
	if err := c.Bind(&inputCredentials); err != nil {
		dto.ResponseWithError(c, http.StatusBadRequest, "invalid request body format", err)
		return
	}

	// compare given password and hash from db
	dbUser, err := cfg.DB.GetUserByLogin(c.Request.Context(), inputCredentials.Login)
	if err != nil {
		dto.ResponseWithError(c, http.StatusUnauthorized, "invalid login or password", err)
		return
	}
	if err = auth.CheckPasswordHash(inputCredentials.Password, dbUser.HashedPassword); err != nil {
		dto.ResponseWithError(c, http.StatusUnauthorized, "invalid login or password", err)
		return
	}

	// make JWT and send it to user
	token, err := auth.MakeJWT(dbUser.ID, cfg.Secret, tokenExpirationTime)
	if err != nil {
		dto.ResponseWithError(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	c.JSON(
		http.StatusOK,
		dto.AuthResponse{
			Token: token,
		},
	)
}
