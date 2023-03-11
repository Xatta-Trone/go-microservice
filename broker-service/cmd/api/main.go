package main

import (
	"log"
	"net/http"
	"runtime"
)

const PORT = ":80"

type Config struct {
}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %s",PORT)

	// define http server 

	URL := ""

	if runtime.GOOS == "windows" {
		URL = "localhost" + PORT
	} else {
		URL = PORT
	}
	srv := &http.Server{
		Addr: URL,
		Handler:app.routes(),
	}

	// start the server 
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}