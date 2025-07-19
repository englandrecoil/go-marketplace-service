package main

import (
	"log"

	"github.com/englandrecoil/go-marketplace-service/internal/config"
	"github.com/gin-gonic/gin"

	_ "github.com/englandrecoil/go-marketplace-service/docs"
	_ "github.com/lib/pq"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @Title       Go Marketplace Service
// @version     1.0
// @description API для маркетплейса
// @host        localhost:8080
// @BasePath    /
func main() {
	apiCfg := config.Init()
	defer apiCfg.Conn.Close()

	router := gin.Default()
	router.POST("/api/reg", apiCfg.HandlerRegister)
	router.POST("/api/auth", apiCfg.HandlerAuth)
	router.POST("/api/ads", apiCfg.HandlerCreateAd)

	router.GET("/api/ads", apiCfg.HandlerGetAds)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(router.Run())
}
