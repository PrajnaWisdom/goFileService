// form request params validation and formatting
package form


type ChunksMetaDataForm struct {
    FileName   string     `json:"filename" binding:"required"`
    FileSize   int64      `json:"filesize" binding:"required"`
}
