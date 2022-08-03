package fsc


import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
    Use:   "fsc",
    Short: "fsc is FileService cli",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("execute fsc cmd.")
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
