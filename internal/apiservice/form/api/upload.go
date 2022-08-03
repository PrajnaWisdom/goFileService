// form request params validation and formatting
package form


type ChunksMetaDataForm struct {
    FileName   string     `json:"filename" binding:"required"`
    FileSize   int64      `json:"filesize" binding:"required"`
}


type GetMissChunksNumberForm struct {
    OwnerID    string     `form:"ownerid" binding:"required"`
    Fuid       string     `form:"fuid" binding:"required"`
}


type UploadChunksForm struct {
    Fuid        string    `form:"fuid" binding:"required"`
    //OwnerID     string    `form:"ownerid" binding:"required"`
    Index       int       `form:"index" binding:"required,gte=1"`
    Md5         string    `form:"md5" binding:"required"`
    //Data        []byte    `form:"data" binding:"required"`
}


type CompleteChunksForm struct {
    Fuid        string    `json:"fuid" binding:"required"`
    //OwnerID     string    `json:"ownerid" binding:"required"`
    Md5         string    `json:"md5" binding:"required"`
}


type GetMetadataForm struct {
    OwnerID    string     `form:"ownerid" binding:"required"`
    Fuid       string     `form:"fuid" binding:"required"`
}
