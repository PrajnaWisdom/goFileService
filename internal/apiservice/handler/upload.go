package handler


import (
    "net/http"
    "log"
    "path"
    "fmt"
    "strings"
    "io"
    "os"
    "io/ioutil"
    "encoding/hex"
    "crypto/md5"
    "github.com/gin-gonic/gin"
    "fileservice/internal/apiservice/config"
    "fileservice/pkg/fileutil"
    "fileservice/pkg/consts"
    "fileservice/pkg/util"
    "fileservice/pkg/account"
    "fileservice/internal/apiservice/form/api"
)


// UploadHandler provides arbitrary file uploads
func UploadHandler(c *gin.Context) {
    var (
        context = Context{C: c}
        form = form.UploadFileForm{}
    )
    cUser, exists := c.Get("user")
    if !exists {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            nil,
            consts.ParamErrorMsg,
        )
        return
    }
    user, ok := cUser.(*account.User)
    if !ok {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            nil,
            consts.ParamErrorMsg,
        )
        return
    }
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
    if err := c.ShouldBind(&form); err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            err.Error(),
            consts.ParamErrorMsg,
        )
        return
    }
    metadata := fileutil.FileMetadata{
        FileName:   fh.Filename,
        FileSize:   fh.Size,
        Fuid:       util.EncodeMD5(util.GenerateUUID()),
        OwnerID:    user.Uid,
    }
    err = metadata.Create(config.GlobaConfig.FileBaseUri)
    if err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            fmt.Sprintf("%v", err),
            consts.ParamErrorMsg,
        )
        return
    }
    chunks := fileutil.FileChunks{
        Fuid:      metadata.Fuid,
        OwnerID:   user.Uid,
        Index:     1,
        Md5:       form.Md5,
    }
    if err := chunks.SaveData(f, config.GlobaConfig.FileBaseUri, false); err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            err.Error(),
            consts.ParamErrorMsg,
        )
        return
    }
    if err := chunks.Create(); err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            err.Error(),
            consts.ParamErrorMsg,
        )
        return
    }
    hash := md5.New()
    hash.Write([]byte(chunks.Md5))
    fMd5 := hex.EncodeToString(hash.Sum(nil))
    metadata.Complete(fMd5)
    uri := metadata.GetFileUri()
    url := strings.Join([]string{config.GlobaConfig.Domain, config.DownloadUri, uri}, "/")
    context.Response(
        http.StatusOK,
        consts.Success,
        map[string]interface{}{"url": url},
        consts.SuccessMsg,
    )
}


// ChunksMetaDataHandler 分片上传元数据
func ChunksMetaDataHandler(c *gin.Context) {
    var (
        context = Context{C: c}
        form = form.ChunksMetaDataForm{}
    )
    cUser, exists := c.Get("user")
    if !exists {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            nil,
            consts.ParamErrorMsg,
        )
        return
    }
    user, ok := cUser.(*account.User)
    if !ok {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            nil,
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
    metadata := fileutil.FileMetadata{
        FileName:   form.FileName,
        FileSize:   form.FileSize,
        Fuid:       util.EncodeMD5(util.GenerateUUID()),
        OwnerID:    user.Uid,
    }
    err := metadata.Create(config.GlobaConfig.FileBaseUri)
    if err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            fmt.Sprintf("%v", err),
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
    metadata, err := fileutil.GetMetaDataByOwnerIDandFuid(form.OwnerID, form.Fuid)
    if err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            "文件元数据不存在",
            consts.ParamErrorMsg,
        )
        return
    }
    numbers, err := metadata.GetUploadedChunksNumber()
    if err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            "文件元数据不存在",
            consts.ParamErrorMsg,
        )
        return
    }
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
    cUser, exists := c.Get("user")
    if !exists {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            nil,
            consts.ParamErrorMsg,
        )
        return
    }
    user, ok := cUser.(*account.User)
    if !ok {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            nil,
            consts.ParamErrorMsg,
        )
        return
    }
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
    _, err = fileutil.GetMetaDataByOwnerIDandFuid(user.Uid, form.Fuid)
    if err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            "文件元数据不存在",
            consts.ParamErrorMsg,
        )
        return
    }
    chunk, _ := fileutil.GetFileChunks(user.Uid, form.Fuid, form.Index)
    if chunk != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            nil,
            "切片文件已存在",
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
    chunks := fileutil.FileChunks{
        Fuid:      form.Fuid,
        OwnerID:   user.Uid,
        Index:     form.Index,
        Md5:       form.Md5,
    }
    if err := chunks.SaveData(f, config.GlobaConfig.FileBaseUri, true); err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            err.Error(),
            consts.ParamErrorMsg,
        )
        return
    }
    if err := chunks.Create(); err != nil {
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


func CompleteChunksHandler(c * gin.Context) {
    var (
        context = Context{C: c}
        form = form.CompleteChunksForm{}
    )
    cUser, exists := c.Get("user")
    if !exists {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            nil,
            consts.ParamErrorMsg,
        )
        return
    }
    user, ok := cUser.(*account.User)
    if !ok {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            nil,
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
    metadata, err := fileutil.GetMetaDataByOwnerIDandFuid(user.Uid, form.Fuid)
    if err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            "文件元数据不存在",
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
    metadata.Complete(form.Md5)
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
    metadata, err := fileutil.GetMetaDataByOwnerIDandFuid(form.OwnerID, form.Fuid)
    if err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            "文件元数据不存在",
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
    metadata, err := fileutil.GetMetaDataByOwnerIDandFuid(form.OwnerID, fuid)
    if err != nil {
        context.Response(
            http.StatusBadRequest,
            consts.ParamError,
            "文件元数据不存在",
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
