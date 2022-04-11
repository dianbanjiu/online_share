package router

import (
	"online_share/api"
	"online_share/common"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

func Start(addr string) {
	r := gin.Default()

	r.Static("/index", "./ui")
	clip := r.Group("/clip")
	{
		clip.GET("", api.GetTheLastClip)
		clip.POST("", api.PushNewClip)
	}

	file := r.Group("/file")
	{
		file.StaticFS("/", http.Dir(common.DefaultSaveDir))
		file.POST("", api.PushFile)
	}

	var srv = http.Server{
		Addr:    addr,
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	var idle = make(chan struct{})
	go func() {
		var stop = make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		<-stop
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
		log.Println("server is stopped. ")
		close(idle)
	}()
	<-idle
}
