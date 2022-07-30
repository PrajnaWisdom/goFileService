package router


import (
    "github.com/gin-gonic/gin"

    "fileservice/internal/apiservice/middleware"
    "fileservice/internal/apiservice/handler"
)


func InitRouter() *gin.Engine {
    r := gin.Default()

    api := r.Group("/api")
    v1 := api.Group("/v1")
    v1.GET("/ping", handler.Ping)
    v1.POST("/upload", handler.UploadHandler)
    chunks := v1.Group("/chunks")
    chunks.Use(middleware.Auth())
    chunks.POST("/metadata", handler.ChunksMetaDataHandler)
    chunks.GET("/metadata", handler.GetMetadataHandler)
    chunks.GET("/missnumbers", handler.GetMissChunksNumberHandler)
    chunks.POST("/data", handler.UploadFileChunksHandler)
    chunks.POST("/complete", handler.CompleteChunksHandler)
    v1.GET("/:ownerid/:fuid", handler.DownloadHandler)

    return r
}
