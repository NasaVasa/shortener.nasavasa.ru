package router

import (
	"github.com/gin-gonic/gin"
	"shortener/db"
)

func SetupRouter(router *gin.Engine, caller db.PgCaller) {
	buildV1Api(router, caller)
}
