// form request params validation and formatting
package form


type ChunksMetaDataForm struct {
    FileName   string     `json:"filename" binding:"required"`
    FileSize   int64      `json:"filesize" binding:"required"`
}


type GetMetaDataForm struct {
    OwnerID    string     `form:"ownerid" binding:"required"`
    Fuid       string     `form:"fuid" binding:"required"`
}


type UploadChunksForm struct {
    Fuid        string    `json:"fuid" binding:"required"`
    OwnerID     string    `json:"ownerid" binding:"required"`
    Index       int       `json:"index" binding:"required"`
    Data        []byte    `json:"data" binding:"required"`
}
