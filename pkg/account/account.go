// account 账号权限包
package account


import (
    "time"
    "log"
    "fmt"
    "strings"
    "crypto/md5"
    "encoding/hex"
    "fileservice/pkg/db"
    "fileservice/pkg/util"
)


const (
    MaxRequestMistime   int64   =  3 * 1000
)


type User struct {
    db.BaseModel
    Account        string        `gorm:"size:255;index;comment:账号"`
    Password       string        `gorm:"size:255;comment:密码"`
    Name           string        `gorm:"size:64;comment:名称"`
    Enabled        bool          `gorm:"default:true;comment:是否启用"`
    Uid            string        `gorm:"size:64;index;comment:uuid"`
}


type Authorize struct {
    db.BaseModel
    Key            string        `gorm:"size:64;index;comment:apikey"`
    Secret         string        `gorm:"size:64;comment:secret"`
    Enabled        bool          `gorm:"comment:是否启用"`
    ExpiredAt      time.Time     `gorm:"comment:过期时间"`
    NeverExpires   bool          `gorm:"default:false;comment:是否永久有效"`
    UserID         uint64        `gorm:"index;comment:用户ID"`
}


func GetUserByAccount(account string) (*User, error) {
    var user User
    err := db.DB.Where("account = ?", account).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}


func (this *User) CheckPassword(pwd string) bool {
    err := util.CheckPasswordHash(this.Password, pwd)
    if err != nil {
        log.Printf("密码校验错误:%v", err)
        return false
    }
    return true
}


func CreateUser(account, pwd string) (*User, error) {
    hashPwd, err := util.GeneratePasswordHash(pwd)
    if err != nil {
        return nil, err
    }
    user := User{
        Account:  account,
        Password: hashPwd,
        Name:     account,
        Uid:      util.EncodeMD5(util.GenerateUUID()),
    }
    err = db.DB.Create(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}


func (this *User) CreateAuthorize(expiredAt time.Time, neverExpires bool) (*Authorize, error) {
    authorize := Authorize{
        UserID:       this.ID,
        ExpiredAt:    expiredAt,
        NeverExpires: neverExpires,
        Key:          util.EncodeMD5(util.GenerateUUID()),
        Secret:       util.EncodeMD5(util.GenerateUUID()),
        Enabled:      true,
    }
    err := db.DB.Create(&authorize).Error
    if err != nil {
        return nil, err
    }
    return &authorize, nil
}



func (this *Authorize) SetEnabled(enabled bool) {
    this.Enabled = enabled
    db.DB.Save(this)
}



func (this *Authorize) SetExpiredAt(expiredAt time.Time) {
    this.ExpiredAt = expiredAt
    db.DB.Save(this)
}


func (this *Authorize) CheckSign(reqSign string, reqTimeStamp int64) bool {
    nowTimeStamp := time.Now().Unix()
    if reqTimeStamp != 0 && nowTimeStamp > reqTimeStamp + MaxRequestMistime {
        log.Println("请求已过期")
        return false
    }
    var user User
    err := db.DB.Where("id = ?", this.UserID).First(&user).Error
    if err != nil {
        log.Printf("查询user失败：%v\n", err)
        return false
    }
    souceStr := fmt.Sprintf("apikey=%v&uid=%v&timestamp=%v=%v", this.Key, user.Uid, reqTimeStamp, this.Secret)
    m := md5.New()
    m.Write([]byte(souceStr))
    sign := strings.ToLower(hex.EncodeToString(m.Sum(nil)))
    log.Printf("字符串[%v]签名结果:%v\n", souceStr, sign)
    if sign != reqSign {
        log.Println("签名错误")
        return false
    }
    return true
}


func GetAvailableAuthorize(apiKey string) (*Authorize, error) {
    var authorize Authorize
    now := time.Now()
    err := db.DB.Joins("inner join `user` on `user`.id = `authorize`.user_id").Where(
        "`authorize`.key = ? and `authorize`.enabled = 1 and (`authorize`.expired_at > ? " +
        " or never_expires = 1 ) and `user`.enabled = 1", apiKey, now).First(&authorize).Error
    if err != nil {
        return nil, err
    }
    return &authorize, nil
}


func GetUserByID(userID uint64) (*User, error) {
    var user User
    err := db.DB.Where("id = ?", userID).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
