package internal

import (
	"os"

	"github.com/swayedev/way"
)

func App() {
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	way.Start(":" + port)
}
