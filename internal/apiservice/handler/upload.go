package handler


import (
    "net/http"
    "log"
    "path"
    "fmt"
    "time"
    "github.com/google/uuid"
    "github.com/gin-gonic/gin"
    "fileservice/internal/apiservice/config"
    "fileservice/pkg/fileutil"
    "fileservice/pkg/consts"
    "fileservice/internal/apiservice/form/api"
)


// UploadHandler provides arbitrary file uploads
func UploadHandler(c *gin.Context) {
    file, _ := c.FormFile("file")
    log.Printf("upload file: %v\n", file.Filename)
    dst := path.Join(config.GlobaConfig.FileBaseUri, file.Filename)
    log.Printf("upload file path: %v\n", dst)
    if err := c.SaveUploadedFile(file, dst); err != nil {
        log.Printf("upload file fail: %v\n", err)
        c.String(http.StatusOK, fmt.Sprintf("'%s' upload fail!", file.Filename))
        return
    }
    c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}


// ChunksMetaDataHandler 分片上传元数据
func ChunksMetaDataHandler(c *gin.Context) {
    var (
        context = Context{C: c}
        form = form.ChunksMetaDataForm{}
        uuid, _ = uuid.NewUUID()
    )
    if err := c.ShouldBind(&form); err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            err.Error(),
            consts.ParamErrorMsg,
        )
        return
    }
    metadata := fileutil.ServerFileMetadata{
        UploadFileMetadata: fileutil.UploadFileMetadata{
            FileName:   form.FileName,
            FileSize:   form.FileSize,
            Fuid:       uuid.String(),
            ChunksNum:  10,
            ModifyTime: time.Now(),
        },
        OwnerID:    uuid.String(),
    }
    metadata.SaveToFile(config.GlobaConfig.FileBaseUri)
    context.Response(
        http.StatusOK,
        consts.Success,
        form,
        consts.SuccessMsg,
    )
    return
}
