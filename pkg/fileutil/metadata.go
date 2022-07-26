// fileutil contains a package of file handling features
package fileutil


import (
    "os"
    "encoding/gob"
    "path"
    "fmt"
    "log"
    "time"
    "strings"
    "encoding/hex"
    "crypto/md5"
    "io/ioutil"
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


const (
    ChunksFileSuffix = ".gochunks"
)


// FileChunks 文件切块结构
type FileChunks struct {
    Fuid        string         // 文件ID，UUID
    OwnerID     string         // 文件所属者
    Index       int            // 切块编号
    Data        []byte         // 切块数据
}


// UploadFileMetadata 客户端传来的文件元数据结构
type UploadFileMetadata struct {
    FileName    string          // 文件名
    Fuid        string          // 文件ID，随机生成的UUID
    FileSize    int64           // 文件大小
    Md5         string          // 文件MD5值
    ChunksNum   int             // 文件分块数
    ModifyTime  time.Time       // 文件修改时间
    ChunksMD5   *map[int]string // 文件分块MD5
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


// IsExists 判断FileChunks切片文件是否存在
func (this *FileChunks) IsExists(baseurl string) bool {
    filepath := path.Join(baseurl, this.OwnerID, this.Fuid, fmt.Sprintf("%v", this.Index) + ChunksFileSuffix)
    if util.IsFile(filepath) {
        return true
    }
    return false
}


func (this *FileChunks) Save(baseurl, cMd5 string) error {
    filepath := path.Join(baseurl, this.OwnerID, this.Fuid, fmt.Sprintf("%v", this.Index) + ChunksFileSuffix)
    hash := md5.New()
    hash.Write(this.Data)
    sMd5 := hex.EncodeToString(hash.Sum(nil))
    if sMd5 != cMd5 {
        return FileError{
            ErrorCode: consts.MD5Inconsistent,
            Msg:       consts.MD5InconsistentMsg,
        }
    }
    return ioutil.WriteFile(filepath, this.Data, 0666)
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


// CheckFileMd5 校验上传文件md5
func (this *ServerFileMetadata) CheckFileMd5(cMd5 string) (bool, *FileError) {
    if len(*this.ChunksMD5) != this.ChunksNum {
        return false, &FileError{
            ErrorCode: consts.ChunksNumError,
            Msg: consts.ChunksNumErrorMsg,
        }
    }
    hash := md5.New()
    for i := 1; i <= this.ChunksNum; i++ {
        hash.Write([]byte((*this.ChunksMD5)[i]))
    }
    sMd5 := hex.EncodeToString(hash.Sum(nil))
    if sMd5 != cMd5 {
        return false, &FileError{
            ErrorCode: consts.MD5Inconsistent,
            Msg:       consts.MD5InconsistentMsg,
        }
    }
    return true, nil
}


func (this *ServerFileMetadata) GetFileUri() string {
    filename := this.Fuid + path.Ext(this.FileName)
    return strings.Join([]string{this.OwnerID, filename}, "/")
}


func (this * ServerFileMetadata) SetChunksMd5(chunkNumber int, cMd5, baseUri string) {
    if this.ChunksMD5 == nil {
        this.ChunksMD5 = &(map[int]string{})
    }
    (*this.ChunksMD5)[chunkNumber] = cMd5
    filepath := path.Join(baseUri, this.OwnerID, this.Fuid, "."+this.Fuid)
    file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
		log.Println("创建元数据文件失败")
		return
	}
    enc := gob.NewEncoder(file)
	err = enc.Encode(this)
	if err != nil {
		log.Println("写元数据文件失败")
		return
	}
    log.Println("写入成功")
	file.Close()
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
