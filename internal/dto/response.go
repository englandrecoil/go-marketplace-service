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

type AuthResponse struct {
	Token string `json:"token"`
}

type CreateAdsResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ImageAddress string    `json:"image_address"`
	Price        int       `json:"price"`
	CreatedAt    time.Time `json:"created_at"`
}

type GetAdsResponse struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	ImageAddress string `json:"image_address"`
	AuthorLogin  string `json:"author_login"`
	Price        int    `json:"price"`
	IsOwner      *bool  `json:"is_owner,omitempty"`
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
