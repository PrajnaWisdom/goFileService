// middleware api服务中间件
package middleware


import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "fileservice/pkg/consts"
    "fileservice/internal/apiservice/handler"
)


// Auth 文件上传认证中间件
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        context := handler.Context{C: c}
        token, ok := c.Request.Header["Auth"]
        if !ok {
            context.Response(
                http.StatusForbidden,
                consts.Forbidden,
                nil,
                consts.ForbiddenMsg,
            )
            c.Abort()
        }
        log.Printf("认证token: %v\n", token)
        c.Next()
    }
}
