package form


type DownloadForm struct {
    OwnerID      string     `uri:"ownerid" binding:"required"`
    Fuid         string     `uri:"fuid" binding:"required"`
}
