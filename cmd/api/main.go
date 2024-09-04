package main

import (
	"fmt"
	"svg-logos-uploader/internal/config"
	"svg-logos-uploader/internal/server"
)

func main() {
	config := config.MustLoadConfig()
	server := server.NewServer(config)

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
