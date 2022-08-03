package fsc


import (
    "log"
    "github.com/spf13/cobra"
    "fileservice/pkg/account"
)


var uAccount string
var Pwd     string

var createUser = &cobra.Command{
    Use:        "cuser",
    Short:      "创建一个用户",
    Long:       "创建一个用户",
    Run:        func(cmd *cobra.Command, args []string) {
        if user, _ := account.GetUserByAccount(uAccount); user != nil {
            log.Printf("用户[%v]已存在", uAccount)
            return
        }
        _, err := account.CreateUser(uAccount, Pwd)
        if err != nil {
            log.Printf("创建用户[%v]失败: %v", uAccount, err)
            return
        }
        log.Printf("账号[%v]创建成功", uAccount)
        return
    },
}


func init() {
    rootCmd.AddCommand(createUser)
    createUser.Flags().StringVarP(&uAccount, "account", "a", "admin", "账号")
    createUser.Flags().StringVarP(&Pwd, "pwd", "p", "admin", "密码")
}
