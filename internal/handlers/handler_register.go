package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type ApiConfig struct {
	Conn *sql.DB
}

func (cfg *ApiConfig) HandlerRegister(c *gin.Context) {
	return
}
