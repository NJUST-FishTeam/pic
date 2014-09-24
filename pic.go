package main

import (
	"html/template"
	"os"

	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"

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
		fileName := saveFile(bytes)
		if fileName == "" {
			w.Write([]byte("文件上传失败：不支持的文件类型"))
			return
		}

		w.Write([]byte("/static/img/" + fileName))
	}
}

func saveFile(pic []byte) string {
	md5_hash := fmt.Sprintf("%x", md5.Sum(pic))
	suffix := ""
	if IsJPEG(pic) {
		suffix = ".jpg"
	} else if IsPNG(pic) {
		suffix = ".png"
	} else {
		return ""
	}
	err := ioutil.WriteFile(path.Join(storePath, md5_hash+suffix), pic, 0644)
	if err != nil {
		log.Fatal("WriteFile: ", err.Error())
		return ""
	}
	return md5_hash + suffix
}

func main() {
	app := cli.NewApp()
	app.Name = "pic"
	app.Usage = "A simple picture bed"
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
