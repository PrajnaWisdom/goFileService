// db 是处理db连接配置相关的包
package db


import (
    "log"
    "gorm.io/driver/mysql"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm/schema"
    "gorm.io/gorm"
    "fileservice/pkg/config"
)


var DB *gorm.DB


func init() {
    var (
        driver gorm.Dialector
        err error
    )
    if config.Config.DB.Driver == "mysql" {
        driver = mysql.Open(config.Config.DB.DSN)
    } else {
        driver = sqlite.Open(config.Config.DB.DSN)
    }
    DB, err = gorm.Open(driver, &gorm.Config{
        PrepareStmt: true,
        NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
    })
    if err != nil {
        log.Fatalf("%v init faile: %v", config.Config.DB.Driver, err)
    }
}
