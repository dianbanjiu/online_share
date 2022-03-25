package router

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dianbanjiu/online_share/api"
	"github.com/dianbanjiu/online_share/db"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Route() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	textShare := engine.Group("/v1/text")
	{
		textShare.GET("", api.GetTextHistory)
		textShare.GET("/:id", api.GetTextById)
		textShare.POST("", api.AddText)
		textShare.DELETE("/:id", api.DeleteTextRecord)
	}

	fileShare := engine.Group("/v1/file")
	{
		fileShare.GET("", api.GetFileHistory)
		fileShare.GET("/:name", api.DownloadFileByName)
		fileShare.HEAD("/:name", api.DownloadFileByNameHead)
		fileShare.POST("", api.UploadFile)
		fileShare.DELETE("/:name", api.DeleteFileRecord)
	}

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]bool{"ok": true})
	})

	return engine
}

func Start(addr string) {
	err := db.DB.AutoMigrate(&api.TextMessage{})
	if err != nil {
		log.Fatalln(err)
	}
	if _, err := os.Stat(api.BaseDir); os.IsNotExist(err) {
		err := os.Mkdir(api.BaseDir, fs.FileMode(os.ModeDir|os.ModePerm))
		if err != nil {
			log.Fatalln(err)
		}
	}

	engine := Route()
	srv := &http.Server{Addr: addr, Handler: engine}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalln("server shutdown err, ", err)
		}
		log.Println("server is exiting. ")
	}()

}
