package main

import (
	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()
	m.Use(martini.Static("static", martini.StaticOptions{Prefix:"static"}))

	m.Get("/", func() string {
		return "Hello world!"
	})
	m.RunOnAddr(":8080")
}