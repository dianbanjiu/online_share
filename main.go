package main

import (
	"online_share/common"
	"online_share/router"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "addr",
				Value: ":8080",
				Usage: "set server listen address. ",
			},
			&cli.StringFlag{
				Name:  "save-dir",
				Value: "/tmp",
				Usage: "set the dir for save upload file. ",
			},
		},
		Action: func(ctx *cli.Context) error {
			common.DefaultListenAddr = ctx.String("addr")
			common.DefaultSaveDir = ctx.String("save-dir")
			router.Start(ctx.String("addr"))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
