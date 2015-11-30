package main

import (
	"github.com/go-martini/martini"
	"github.com/golang/glog"
	"net/http"
	"path/filepath"
	"k8s-ui/templates"
	"os"
	"runtime"
	"flag"
	"strings"
	"k8s-ui/comm"
)

var(
	addr = ":8080"
)

func main() {
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	flag.StringVar(&addr, "addr", addr, "监听地址,如\":8080\", 默认: \":8080\"")
	flag.Parse()

	defer glog.Flush()
	glog.CopyStandardLogTo("INFO")

	m := martini.Classic()
	m.Use(martini.Static("static", martini.StaticOptions{Prefix:"static"}))

	m.Group("/json", func(r martini.Router) {
		r.Group("/namespace", comm.Namespaces)
	})

	m.Get("/", serveTemplate)
	m.RunOnAddr(addr)
}


func serveTemplate(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Base(r.URL.Path);
	if len(fp) == 0 || fp == "/"{
		fp = "index.html"
	}

	if strings.HasPrefix(r.RequestURI, "/json"){
		jsonHandler(w, r)
		return
	}

	if err := templates.RenderTemplate(w, fp, nil); err != nil{
		glog.Warning(err)
		http.NotFound(w, r)
	}
}

func jsonHandler(w http.ResponseWriter, r *http.Request){

}