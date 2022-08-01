package util


import (
    "os"
    "log"
    "crypto/md5"
    "encoding/hex"
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


func EncodeMD5(value string) string {
    m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
