// fileutil contains a package of file handling features
package fileutil


import (
    "os"
    // "encoding/gob"
    "path"
    "log"
    "time"
)


type FileState   int

const (
    Error          FileState = -1
    Init           FileState = 0
    Active         FileState = 1
    Downloading    FileState = 2
    Uploading      FileState = 3
)


// FileChunks 文件切块结构
type FileChunks struct {
    Fuid        string         // 文件ID，UUID
    Index       int            // 切块编号
    Data        []byte         // 切块数据
}


// UploadFileMetadata 客户端传来的文件元数据结构
type UploadFileMetadata struct {
    FileName    string         // 文件名
    Fuid        string         // 文件ID，随机生成的UUID
    FileSize    int64          // 文件大小
    Md5         string         // 文件MD5值
    ChunksNum   int            // 文件分块数
    ModifyTime  time.Time      // 文件修改时间
}


type ServerFileMetadata struct {
    UploadFileMetadata
    OwnerID     string         // 文件所属者
    State       int            // 文件状态
}


func (this *ServerFileMetadata) SaveToFile(baseUrl string) (*os.File, error) {
    filepath := path.Join(baseUrl, this.OwnerID, this.Fuid, "."+this.Fuid)
    //f, err := os.OpenFile()
    log.Printf("文件地址: %v\n", filepath)
    return nil, nil
}
