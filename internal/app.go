package internal

import (
	"github.com/swayedev/way"
)

func App() {
	// app := Config{}
	// app.Get()
	// app.Set()

	way := way.New()
	if err := way.Db().PgxOpen(); err != nil {
		panic(err)
	}
	defer way.Db().Close()

	way.GET("/", RootHandler)
	way.GET("/upload", DemoHandler)
	way.POST("/upload", UploadWithOptionalEncryptionHandler)
	way.GET("/view/{slug}", ViewHandler)
	way.GET("/public", PublicHandler)

	way.Start(":8080")
}
