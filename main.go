package main


import (
    "github.com/gin-gonic/gin"
)


func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        //输出json结果给调用方
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Run()
}
