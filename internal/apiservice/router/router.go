package router


import (
    "github.com/gin-gonic/gin"

    "fileservice/internal/apiservice/handler"
)


func InitRouter() *gin.Engine {
    r := gin.Default()

    api := r.Group("/api")
    v1 := api.Group("/v1")
    v1.GET("/ping", handler.Ping)

    return r
}