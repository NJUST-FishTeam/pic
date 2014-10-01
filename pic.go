package main

import (
	"html/template"
	"os"

	"io/ioutil"
	"log"
	"net/http"

	"github.com/NJUST-FishTeam/pic/picfs"
	"github.com/codegangsta/cli"
)

var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
	storePath     = "./img/"
)

func indexHandle(w http.ResponseWriter, r *http.Request) {
	if err := indexTemplate.Execute(w, nil); err != nil {
		log.Fatal("Execute: ", err.Error())
		return
	}
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, _, err := r.FormFile("file")
		if err != nil {
			log.Fatal("FormFile: ", err.Error())
			return
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal("Close: ", err.Error())
				return
			}
		}()

		bytes, err := ioutil.ReadAll(file)
		fileName, err := picfs.SavePic(bytes, storePath)
		if err != nil {
			w.Write([]byte("文件上传失败：不支持的文件类型"))
			return
		}

		w.Write([]byte("/static/images/" + fileName))
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "pic"
	app.Usage = "A simple picture bed"
	app.Author = "maemual"
	app.Version = "0.3"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Value: storePath,
			Usage: "The pictures store path",
		},
	}
	app.Action = func(c *cli.Context) {
		if _, err := os.Stat(c.String("path")); err != nil && !os.IsExist(err) {
			os.MkdirAll(c.String("path"), os.ModePerm)
		}
		storePath = c.String("path")

		http.HandleFunc("/", indexHandle)
		http.HandleFunc("/upload/", uploadHandle)
		http.ListenAndServe(":8090", nil)
	}

	app.Run(os.Args)
}
