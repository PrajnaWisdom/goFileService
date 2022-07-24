// goFileService web main package
package main


import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
    "fmt"
    "github.com/gin-gonic/gin"

    "fileservice/internal/apiservice/config"
    "fileservice/internal/apiservice/router"
)


func init() {
    log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
    config.SetupConfig()
}


func main() {
    gin.SetMode(config.GlobaConfig.Mode)

    r := router.InitRouter()

    addr := fmt.Sprintf(":%d", config.GlobaConfig.Port)
    server := &http.Server{
        Addr:     addr,
        Handler:  r,
    }

    go func() {
        // 服务连接
        log.Printf("\nFile server start listen: %v\n", config.GlobaConfig.Port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
    }()

    // 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
