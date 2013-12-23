package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func htmlGen(m *martini.ClassicMartini, dirname string) string {
	dirAbs, err := filepath.Abs(dirname)
	if err != nil {
		fmt.Println("error")
	}
	fileInfos, err := ioutil.ReadDir(dirAbs)
	if err != nil {
		fmt.Println("error")
	}

	html := "<h1>Dir</h1>"
	html += "<ul>"
	for _, fileInfo := range fileInfos {
		fileName := fileInfo.Name()

		var path string
		var realpath string
		if dirname == "."{
			path = "/:fileName"
			realpath ="/"+fileName;
		}else{
			path = "/" + dirname + "/:fileName"
			realpath = "/" + dirname + "/"+fileName
		}

		if !fileInfo.IsDir() {
			html += "<li><a href=\"" + realpath + "\">" + fileName + "</a></li>"
		} else {
			html += "<li><a href=\"" + realpath + "\"><font color=\"#FF0000\">" + fileName + "</font></a></li>"
		}

		m.Get(path, func(params martini.Params) string {
			fileName := params["fileName"]
			file, err := os.Open(dirAbs+"/"+fileName)
			if err != nil {
				fmt.Println("file open error")
			}
			defer file.Close()
			fileInfo, err := file.Stat()
			if !fileInfo.IsDir() {
				if err != nil {
					fmt.Println("IsDir error")

				}
				data := make([]byte, fileInfo.Size())
				count, err := file.Read(data)
				if err != nil {
					fmt.Println("file read error")
				}
				return string(data[:count])
			} else {
				newdirname := strings.Split(path, ":")[0]+fileName
				return htmlGen(m, newdirname[1:])
			}
		})
	}
	html += "</ul>"
	return html
}

func main() {

	m := martini.Classic()
	m.Get("/favicon.ico", func(){
	})
	html := htmlGen(m,".")
	m.Get("/", func() string {
		return html
	})
	m.Get("/gejigeji", func() string {
		return "Hello gejigeji!"
	})
	m.Run()
}
