package handler


import (
    "net/http"
    "log"
    "path"
    "fmt"
    "time"
    "strings"
    "io"
    "os"
    "io/ioutil"
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
    if form.ChunksNum > fileutil.MaxChunksNumber {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            fmt.Sprintf("ChunksNum 不能大于:%v", fileutil.MaxChunksNumber),
            consts.ParamErrorMsg,
        )
        return
    }
    metadata := fileutil.ServerFileMetadata{
        UploadFileMetadata: fileutil.UploadFileMetadata{
            FileName:   form.FileName,
            FileSize:   form.FileSize,
            Fuid:       util.GenerateUUID(),
            ChunksNum:  form.ChunksNum,
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


// GetMissChunksNumberHandler 获取缺失的分片编号
func GetMissChunksNumberHandler(c *gin.Context) {
    var (
        context = Context{C: c}
        form = form.GetMissChunksNumberForm{}
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
    log.Println(metadata)
    numbers := metadata.GetMissChunksNumber()
    context.Response(
        http.StatusOK,
        consts.Success,
        numbers,
        consts.SuccessMsg,
    )
    return
}


func UploadFileChunksHandler(c *gin.Context) {
    var (
        context = Context{C: c}
        form = form.UploadChunksForm{}
    )
    fh, err := c.FormFile("file")
    if err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            err.Error(),
            consts.ParamErrorMsg,
        )
        return
    }
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
    if form.Index > metadata.ChunksNum {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            "切片索引[index]不能大于切片数",
            consts.ParamErrorMsg,
        )
        return
    }
    f, err := fh.Open()
    if err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            err.Error(),
            consts.ParamErrorMsg,
        )
        return
    }
    content, err := ioutil.ReadAll(f)
    if err != nil {
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
        Data:      content,
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
    if err := chunks.Save(config.GlobaConfig.FileBaseUri, form.Md5); err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            err.Error(),
            consts.ParamErrorMsg,
        )
        return
    }
    go metadata.SetChunksMd5(chunks.Index, form.Md5, config.GlobaConfig.FileBaseUri)
    context.Response(
        http.StatusOK,
        consts.Success,
        nil,
        consts.SuccessMsg,
    )
}


func CompleteChunksHandler(c * gin.Context) {
    var (
        context = Context{C: c}
        form = form.CompleteChunksForm{}
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
    ok, fileErr := metadata.CheckFileMd5(form.Md5);
    if !ok {
        context.Response(
            http.StatusOK,
            fileErr.ErrorCode,
            fileErr.Error(),
            fileErr.Msg,
        )
        return
    }
    uri := metadata.GetFileUri()
    url := strings.Join([]string{config.GlobaConfig.Domain, config.DownloadUri, uri}, "/")
    context.Response(
        http.StatusOK,
        consts.Success,
        map[string]interface{}{"url": url},
        consts.SuccessMsg,
    )
    return
}


// GetMetadataHandler 获取上传文件元数据
func GetMetadataHandler(c *gin.Context) {
    var (
        context = Context{C: c}
        form = form.GetMetadataForm{}
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


func DownloadHandler(c *gin.Context) {
    var (
        context = Context{C: c}
        form = form.DownloadForm{}
    )
    if err := c.ShouldBindUri(&form); err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            err.Error(),
            consts.ParamErrorMsg,
        )
        return
    }
    fuid := strings.Split(form.Fuid, ".")[0]
    filepath := path.Join(config.GlobaConfig.FileBaseUri, form.OwnerID, fuid, "."+fuid)
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
    pr, pw := io.Pipe()
    go func() {
        defer pw.Close()
        file_uri := path.Join(config.GlobaConfig.FileBaseUri, metadata.OwnerID, metadata.Fuid)
        for i := 1; i <= metadata.ChunksNum; i++ {
            fileurl := path.Join(file_uri, fmt.Sprintf("%v", i) + fileutil.ChunksFileSuffix)
            file, err := os.Open(fileurl)
            if err != nil {
                log.Println(err)
                return
            }
            content, err := ioutil.ReadAll(file)
            if err != nil {
                file.Close()
                log.Println(err)
                return
            }
            pw.Write(content)
            file.Close()
        }
    }()
    headers := map[string]string{
        "Content-disposition": "attachment; filename=\""+metadata.FileName+"\"",
    }
    c.DataFromReader(
        http.StatusOK,
        metadata.FileSize,
        "application/octet-stream",
        pr,
        headers,
    )
}
