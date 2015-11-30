package templates

import (
	"html/template"
	"path/filepath"
	"log"
	"net/http"
	"fmt"
	"os"
	"github.com/oxtoacart/bpool"
	"github.com/go-fsnotify/fsnotify"
//	"time"
	"path"
)

const Dir_Templates string = "templates/"
const ext_tmpl  string = ".html"

var Templates map[string]*template.Template
var bufpool *bpool.BufferPool
var templatesDir string

func init(){
	bufpool = bpool.NewBufferPool(64)
	templatesDir = Dir_Templates

	initTemplate()

	go reloadTemplate()
}

func initTemplate(){
	Templates = make(map[string]*template.Template)

	files := make([]string, 0)
	filepath.Walk(templatesDir, func(path string, fileInfo os.FileInfo, err error)error{
		if fileInfo.IsDir(){
			tmpTmpls, err := filepath.Glob(path + fmt.Sprintf("/*%s", ext_tmpl))
			if err != nil {
				log.Fatal(err)
			}
			files = append(files, tmpTmpls...)
		}
		return nil
	})

	log.Printf("%v", files)
	// Generate our templates map from our layouts/ and includes/ directories
	for _, layout := range files {
		Templates[filepath.Base(layout)] = template.Must(template.New("").Delims("<%", "%>").ParseFiles(files...))
	}
}

// 监控templates目录,当模板发生改变后自动加载
func reloadTemplate(){
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	filepath.Walk(Dir_Templates, func(path string, fileinfo os.FileInfo, err error)error{
		if fileinfo.IsDir(){
			err = watcher.Add(path)
			if err != nil {
				log.Fatal(err)
			}
		}
		return err
	})

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op == fsnotify.Remove|fsnotify.Rename{
					watcher.Remove(event.Name)
				}else{
					if info, err := os.Stat(event.Name); err == nil{
						if info.IsDir(){
							if err = watcher.Add(event.Name);err != nil{
								log.Println(err)
							}
						}
					}
				}
				if path.Ext(event.Name) == ext_tmpl{
//					time.Sleep(3 * time.Second)
					initTemplate()
				}
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()
}

// renderTemplate is a wrapper around template.ExecuteTemplate.
func RenderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	// Ensure the template exists in the map.
	tmpl, ok := Templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	// Create a buffer to temporarily write to and check if any errors were encounted.
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}

	// Set the header and write the buffer to the http.ResponseWriter
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
	return nil
}