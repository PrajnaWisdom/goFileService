// fileutil contains a package of file handling features
package fileutil


import (
    "os"
    "encoding/gob"
    "path"
    "fmt"
    "log"
    "time"
    "fileservice/pkg/util"
    "fileservice/pkg/consts"
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


type FileError struct {
    ErrorCode   int64
    Msg         string
}


func (e FileError) Error() string {
    return fmt.Sprintf("file error[%v] msg:%v", e.ErrorCode, e.Msg)
}


func (this *ServerFileMetadata) IsExistsMetaDataFile(path string) bool {
    if util.IsFile(path) {
        return true
    }
    return false
}


func (this *ServerFileMetadata) SaveToFile(baseUrl string) error {
    saveDir := path.Join(baseUrl, this.OwnerID, this.Fuid)
    filepath := path.Join(saveDir, "."+this.Fuid)
    //f, err := os.OpenFile()
    log.Printf("文件地址: %v\n", filepath)
    if this.IsExistsMetaDataFile(filepath) {
        return FileError{ErrorCode: consts.FileIsExists, Msg: consts.FileIsExistsMsg}
    }
    err := os.MkdirAll(saveDir, 0766)
    if err != nil {
        log.Println(err)
        return err
    }
    file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
		log.Println("创建元数据文件失败")
		return err
	}
    enc := gob.NewEncoder(file)
	err = enc.Encode(this)
	if err != nil {
		log.Println("写元数据文件失败")
		return err
	}
    log.Println("写入成功")
	file.Close()
    return nil
}


// LoadMetaDataFile 从文件里面加载metadata信息
func LoadMetaDataFile(path string) (*ServerFileMetadata, error) {
    file, err := os.Open(path)
	if err != nil {
		log.Println("获取文件状态失败，文件路径：", path)
		return nil, err
	}

	var metadata ServerFileMetadata
	filedata := gob.NewDecoder(file)
	err = filedata.Decode(&metadata)
	if err != nil {
		log.Println("格式化文件元数据失败, err", err)
		return nil, err
	}
	return &metadata, nil
}
