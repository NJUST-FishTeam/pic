package main

import (
	"html/template"
	"os"

	"log"
	"net/http"

	"github.com/codegangsta/cli"
)

var indexTemplate = template.Must(template.ParseFiles("index.html"))

func indexHandle(w http.ResponseWriter, r *http.Request) {
	if err := indexTemplate.Execute(w, nil); err != nil {
		log.Fatal("Execute: ", err.Error())
		return
	}
}

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

		http.HandleFunc("/", indexHandle)
		http.ListenAndServe(":8090", nil)
	}

	app.Run(os.Args)
}
