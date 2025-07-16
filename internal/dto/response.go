package dto

import (
	"log"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type RegisterResponse struct {
	Login        string `json:"login"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func ResponseWithError(c *gin.Context, code int, errMsg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %v", err)
	}

	c.JSON(code, ErrorResponse{
		Error: errMsg,
	})
}
