package main


import (
    "log"
    "fileservice/pkg/config"
    "fileservice/pkg/db"
    "fileservice/pkg/fileutil"
)


func main(){
    Db := db.DB.Debug()
    if config.Config.DB.Driver == "mysql" {
        Db = db.DB.Set("gorm:table_options", "ENGINE=InnoDB")
    }
    err := Db.AutoMigrate(
        &fileutil.FileMetadata{},
        &fileutil.FileChunks{},
    )
    if err != nil {
        log.Fatalf("migrate err: %v", err)
    }
    if !Db.Migrator().HasIndex(&fileutil.FileMetadata{}, "idx_owner_fuid"){
        if err := Db.Migrator().CreateIndex(&fileutil.FileMetadata{}, "idx_owner_fuid"); err != nil {
            log.Fatalf("migrate err: %v", err)
        }
    }
    if !Db.Migrator().HasIndex(&fileutil.FileChunks{}, "chunk_idx_owner_fuid"){
        if err := Db.Migrator().CreateIndex(&fileutil.FileChunks{}, "chunk_idx_owner_fuid"); err != nil {
            log.Fatalf("migrate err: %v", err)
        }
    }
}
