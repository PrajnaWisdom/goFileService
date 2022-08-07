// fileutil contains a package of file handling features
package fileutil


import (
    "os"
    "path"
    "fmt"
    "log"
    "io"
    "strings"
    "errors"
    "encoding/hex"
    "crypto/md5"
    "io/ioutil"
    "gorm.io/gorm"
    "fileservice/pkg/consts"
    "fileservice/pkg/db"
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
    MaxChunksNumber  = 10000
    MaxChunksSize    = 1024 * 1024 * 1024
    MinChunksSize    = 1024 * 1024
)


// FileChunks 文件切块结构
type FileChunks struct {
    db.BaseModel
    Fuid        string     `gorm:"size:64;index:chunk_idx_owner_fuid,priority:2;comment:文件ID"`
    OwnerID     string     `gorm:"size:64;index:chunk_idx_owner_fuid,priority:1;comment:文件所属者"`
    Md5         string     `gorm:"size:64;comment:切片MD5值"`
    Index       int        `gorm:"comment:切块编号"`
    // Data        []byte         // 切块数据
}


type FileMetadata struct {
    db.BaseModel
    FileName    string     `gorm:"size:255;comment:文件名"`
    Fuid        string     `gorm:"size:64;index:idx_owner_fuid,priority:2;comment:文件ID"`
    OwnerID     string     `gorm:"size:64;index:idx_owner_fuid,priority:1;comment:文件所属者"`
    FileSize    int64      `gorm:"comment:文件大小"`
    Md5         string     `gorm:"size:64;comment:文件MD5值"`
    ChunksNum   int        `gorm:"default:0;comment:文件分块数"`
    //ChunksMD5   *map[int]string // 文件分块MD5
    State       FileState  `gorm:"default:0;comment:文件状态"`
}


type FileError struct {
    ErrorCode   int64
    Msg         string
}


func (e FileError) Error() string {
    return fmt.Sprintf("file error[%v] msg:%v", e.ErrorCode, e.Msg)
}


// IsExists 判断FileChunks切片文件是否存在
func GetFileChunks(ownerID, fuid string, index int) (*FileChunks, error) {
    var chunk FileChunks
    err := db.DB.Where("owner_id = ? and fuid = ? and `index` = ?", ownerID, fuid, index).First(&chunk).Error
    if err != nil{
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &chunk, nil
}


func (this *FileChunks) Create() error {
    return db.DB.Create(this).Error
}


// SaveData 将分片数据保存到磁盘
func (this *FileChunks) SaveData(reader io.Reader, baseurl string, checkMd5 bool) error {
    filepath := path.Join(baseurl, this.OwnerID, this.Fuid, fmt.Sprintf("%v", this.Index) + ChunksFileSuffix)
    content, err := ioutil.ReadAll(reader)
    if err != nil {
        return err
    }
    hash := md5.New()
    hash.Write(content)
    sMd5 := hex.EncodeToString(hash.Sum(nil))
    if !checkMd5 {
        this.Md5 = sMd5
    } else if sMd5 != this.Md5 {
        return FileError{
            ErrorCode: consts.MD5Inconsistent,
            Msg:       consts.MD5InconsistentMsg,
        }
    }
    return ioutil.WriteFile(filepath, content, 0666)
}


// IsExistsMetaDataFile 判断分片上传原数据文件是否存在
func (this *FileMetadata) IsExistsMetaData() bool {
    var metadata FileMetadata
    err := db.DB.Where("owner_id = ? and fuid = ?", this.OwnerID, this.Fuid).First(&metadata).Error
    if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
        return false
    }
    return true
}


// Create 在数据库中创建一条FileMetadata记录
func (this *FileMetadata) Create(baseUrl string) error {
    saveDir := path.Join(baseUrl, this.OwnerID, this.Fuid)
    log.Printf("文件保存地址: %v\n", saveDir)
    if this.IsExistsMetaData() {
        return FileError{ErrorCode: consts.FileIsExists, Msg: consts.FileIsExistsMsg}
    }
    err := os.MkdirAll(saveDir, 0766)
    if err != nil {
        log.Printf("保存磁盘失败: %v", err)
        return err
    }
    err = db.DB.Create(this).Error
    return err
}


// CheckFileMd5 校验上传文件md5
func (this *FileMetadata) CheckFileMd5(cMd5 string) (bool, *FileError) {
    chunks, err := GetFileChunksByOwnerIDandFuid(this.OwnerID, this.Fuid)
    if err != nil {
        return false, &FileError{
            ErrorCode: consts.ChunksNumError,
            Msg: consts.ChunksNumErrorMsg,
        }
    }
    nums := len(*chunks)
    hash := md5.New()
    for i := 1; i <= nums; i++ {
        chunk := (*chunks)[i - 1]
        if chunk.Index != i {
            return false, &FileError{
                ErrorCode: consts.ChunksNumError,
                Msg: consts.ChunksNumErrorMsg,
            }
        }
        hash.Write([]byte(chunk.Md5))
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


// Complete 更新FileMetadata的md5值和State
func (this *FileMetadata) Complete(cMd5 string) {
    var nums int
    chunks, err := GetFileChunksByOwnerIDandFuid(this.OwnerID, this.Fuid)
    if err == nil {
        nums = len(*chunks)
    }
    this.Md5 = cMd5
    this.ChunksNum = nums
    this.State = Active
    db.DB.Save(this)
}


// SetState 设置文件状态
func (this *FileMetadata) SetState(state FileState, commit bool) {
    this.State = state
    if commit {
        db.DB.Save(this)
    }
}


// GetFileUri 获取上传文件的相对路径
func (this *FileMetadata) GetFileUri() string {
    filename := this.Fuid + path.Ext(this.FileName)
    return strings.Join([]string{this.OwnerID, filename}, "/")
}


// GetMissChunksNumber 获取缺失的分片编号
func (this *FileMetadata) GetUploadedChunksNumber() ([]int, error) {
    numbers := []int{}
    chunks, err := GetFileChunksByOwnerIDandFuid(this.OwnerID, this.Fuid)
    if err != nil {
        return nil, err
    }
    for _, chunk := range *chunks {
        numbers = append(numbers, chunk.Index)
    }
    return numbers, nil
}


// LoadMetaDataFile 从文件里面加载metadata信息
func GetMetaDataByOwnerIDandFuid(ownerID, fuid string) (*FileMetadata, error) {
    var metadata FileMetadata
    err := db.DB.Where("owner_id = ? and fuid = ?", ownerID, fuid).First(&metadata).Error
	return &metadata, err
}


// GetFileChunksByOwnerIDandFuid 根据ownerID和fuid获取已上传的FileChunks列表，并按照index排序
func GetFileChunksByOwnerIDandFuid(ownerID, fuid string) (*[]FileChunks, error) {
    var chunks []FileChunks
    err := db.DB.Where("owner_id = ? and fuid = ?", ownerID, fuid).Order("`index`").Find(&chunks).Error
    if err != nil {
        return nil, err
    }
    return &chunks, nil
}
