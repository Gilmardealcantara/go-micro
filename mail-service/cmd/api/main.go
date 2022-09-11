package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
}

const wePort = "80"

func main() {
	app := Config{}

	log.Println("Starting mail service on port: ", wePort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", wePort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
