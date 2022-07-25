package handler


import (
    "net/http"
    "log"
    "path"
    "fmt"
    "time"
    "github.com/gin-gonic/gin"
    "fileservice/internal/apiservice/config"
    "fileservice/pkg/fileutil"
    "fileservice/pkg/consts"
    "fileservice/pkg/util"
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
            Fuid:       util.GenerateUUID(),
            ChunksNum:  10,
            ModifyTime: time.Now(),
        },
        OwnerID:    util.GenerateUUID(),
    }
    metadata.SaveToFile(config.GlobaConfig.FileBaseUri)
    context.Response(
        http.StatusOK,
        consts.Success,
        metadata,
        consts.SuccessMsg,
    )
    return
}


// LoadMeatDataFileHandler 加载metadata信息 
func LoadMeatDataFileHandler(c *gin.Context) {
    var (
        context = Context{C: c}
        form = form.GetMetaDataForm{}
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
    filepath := path.Join(config.GlobaConfig.FileBaseUri, form.OwnerID, form.Fuid, "."+form.Fuid)
    metadata, err := fileutil.LoadMetaDataFile(filepath)
    if err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            err.Error(),
            consts.ParamErrorMsg,
        )
        return
    }
    context.Response(
        http.StatusOK,
        consts.Success,
        metadata,
        consts.SuccessMsg,
    )
    return
}


func UploadFileChunksHandler(c *gin.Context) {
    var (
        context = Context{C: c}
        form = form.UploadChunksForm{}
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
    chunks := fileutil.FileChunks{
        Fuid:      form.Fuid,
        OwnerID:   form.OwnerID,
        Index:     form.Index,
        Data:      form.Data,
    }
    if chunks.IsExists(config.GlobaConfig.FileBaseUri) {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            nil,
            "切片文件已存在",
        )
        return
    }
    if err := chunks.Save(config.GlobaConfig.FileBaseUri); err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            err.Error(),
            consts.ParamErrorMsg,
        )
        return
    }
    context.Response(
        http.StatusOK,
        consts.Success,
        nil,
        consts.SuccessMsg,
    )
}
