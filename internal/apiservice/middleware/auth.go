// middleware api服务中间件
package middleware


import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "fileservice/pkg/consts"
    "fileservice/internal/apiservice/handler"
    "fileservice/pkg/account"
)


type AuthForm struct {
    AuthSign           string    `header:"Auth-Sign" binding:"required"`
    ApiKey             string    `header:"Api-Key" binding:"required"`
    RequestTimeStamp   int64     `header:"Request-TimeStamp" binding:"required"`
}


// Auth 文件上传认证中间件
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        context := handler.Context{C: c}
        form := AuthForm{}
        log.Println(c.Request.Header)
        if err := c.ShouldBindHeader(&form); err != nil {
            context.Response(
                http.StatusBadRequest,
                consts.ParamError,
                err.Error(),
                consts.ParamErrorMsg,
            )
            c.Abort()
            return
        }
        authorize, err := account.GetAvailableAuthorize(form.ApiKey)
        if err != nil {
            log.Printf("获取authorize[%v]失败: %v\n", form.ApiKey, err)
            context.Response(
                http.StatusForbidden,
                consts.Forbidden,
                nil,
                consts.ForbiddenMsg,
            )
            c.Abort()
            return
        }
        if ok := authorize.CheckSign(form.AuthSign, form.RequestTimeStamp); !ok {
            log.Printf("签名认证失败失败\n")
            context.Response(
                http.StatusForbidden,
                consts.Forbidden,
                nil,
                consts.ForbiddenMsg,
            )
            c.Abort()
            return
        }
        user, err := account.GetUserByID(authorize.UserID)
        if err != nil {
            log.Printf("获取用户[%v]失败: %v\n", authorize.UserID, err)
            context.Response(
                http.StatusForbidden,
                consts.Forbidden,
                nil,
                consts.ForbiddenMsg,
            )
            c.Abort()
            return
        }
        //c.Set("authorize", authorize)
        c.Set("user", user)
        log.Printf("认证token: %v\n", form)
        c.Next()
    }
}
