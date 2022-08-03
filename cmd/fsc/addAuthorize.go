package fsc


import (
    "log"
    "time"
    "github.com/spf13/cobra"
    "fileservice/pkg/account"
)


var Account       string
var ExpiredAt     string
var NerverExpired bool

const defaultTimeFormat = "2006-01-02"

var author = &cobra.Command{
    Use:        "genApikey",
    Short:      "添加apikey",
    Run:        func(cmd *cobra.Command, args []string) {
        expiredAt, err := time.Parse(defaultTimeFormat, ExpiredAt)
        if err != nil {
            log.Printf("过期时间格式错误:%v", err)
            return
        }
        user, err := account.GetUserByAccount(Account)
        if err != nil {
            log.Printf("获取用户[%v]失败: %v", Account, err)
            return
        }

        authorize, err := user.CreateAuthorize(expiredAt, NerverExpired)
        if err != nil {
            log.Printf("生成apikey失败: %v", err)
            return
        }
        log.Printf("账号[%v]生成apikey成功:\n[key]:%v\n[secret]:%v", Account, authorize.Key, authorize.Secret)
        return
    },
}


func init() {
    rootCmd.AddCommand(author)
    author.Flags().StringVarP(&Account, "account", "a", "admin", "账号")
    author.Flags().StringVarP(&ExpiredAt, "expired_at", "e", defaultTimeFormat, "过期时间")
    author.Flags().BoolVarP(&NerverExpired, "nerver_expired", "n", false, "是否永久有效")
}
