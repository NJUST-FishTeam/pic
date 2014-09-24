package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "pic"
	app.Usage = "A simple picture bed"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Value: "./img/",
			Usage: "The pictures store path",
		},
	}
	app.Action = func(c *cli.Context) {
		if _, err := os.Stat(c.String("path")); err != nil && !os.IsExist(err) {
			os.MkdirAll(c.String("path"), os.ModePerm)
		}
	}

	app.Run(os.Args)
}
