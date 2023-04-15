package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"shortener/db"
	_ "shortener/docs"
	r "shortener/router"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
	}
}

// @title           Shortener service
// @version         1.0
// @description     Сервис, позволяющий укорачивать ссылки.
func main() {
	router := gin.Default()
	caller, err := db.NewDB()
	if err != nil {
		panic(err)
	}
	r.SetupRouter(router, caller)
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	_ = router.Run(":8080")
}
