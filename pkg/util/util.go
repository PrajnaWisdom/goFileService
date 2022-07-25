package util


import (
    "os"
    "log"
    uuid "github.com/satori/go.uuid"
)


func GenerateUUID() string {
    return uuid.NewV4().String()
}


func IsFile(filepath string) bool {
    fileInfo, err := os.Stat(filepath)
    if err != nil {
        log.Printf("error: %v\n", err)
        return false
    }
    return !fileInfo.IsDir()
}
