package main

import (
	"log"

	"github.com/englandrecoil/go-marketplace-service/internal/config"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lib/pq"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(router.Run())
}
