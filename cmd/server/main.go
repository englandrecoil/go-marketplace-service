package main

import (
	"github.com/englandrecoil/go-marketplace-service/internal/config"
	"github.com/gin-gonic/gin"
)

// @Title       Go Marketplace Service
// @version     1.0
// @description API для маркетплейса.
// @host        localhost:8080
// @BasePath    /
func main() {
	apiCfg := config.Init()
	defer apiCfg.Conn.Close()

	router := gin.Default()
	router.POST("/api/reg", apiCfg.HandlerRegister)
}
