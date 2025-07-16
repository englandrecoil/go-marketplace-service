package dto

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type RegisterResponse struct {
	ID        uuid.UUID `json:"id"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
