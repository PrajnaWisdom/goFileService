package handler


import (
    "github.com/gin-gonic/gin"
)


type Context struct {
    C  *gin.Context
}


func (this Context) Response(httpStatus int, code int64, data interface{}, msg string) {
    this.C.JSON(httpStatus, map[string]interface{}{
        "code": code,
        "msg": msg,
        "data": data,
    })
}
