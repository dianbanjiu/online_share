package main

import (
	"log"
	"os"

	_ "github.com/dianbanjiu/online_share/db"
	_ "net/http/pprof"
	_ "github.com/dianbanjiu/online_share/docs"
	"github.com/dianbanjiu/online_share/router"
	"github.com/urfave/cli/v2"
)

// @title           Online Share Swagger API
// @version         1.0
// @description     This is a server for text and file share.
// @termsOfService  http://swagger.io/terms/

// @contact.name   dianbanjiu
// @contact.url    https://github.com/dianbanjiu/online_share
// @contact.email  dianbanjiu@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "addr",
				Value: ":8080",
				Usage: "server listen address",
			},
		},
	}

	app.Action = func(ctx *cli.Context) error {
		router.Start(ctx.String("addr"))
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}

}
