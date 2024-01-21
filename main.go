package main

import (
	"github.com/swayedev/way"
)

var savedFileHashes = make(map[string]bool)
var maxMemory int64 = 1024 * 1024 * 10 // 10 MB

func main() {
	app := Config{}
	app.Get()
	app.Set()

	way := way.New()
	if err := way.Db().PgxOpen(); err != nil {
		panic(err)
	}
	defer way.Db().Close()

	way.GET("/", rootHandler)
	way.GET("/upload", demoHandler)
	way.POST("/upload", uploadHandler)
	way.GET("/view/{slug}", viewFileHandler)

	way.Start(":8080")
}
