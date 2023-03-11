package main

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const PORT = ":80"

func main() {

	app := Config{
		Mailer: createMail(),
	}

	log.Println("Starting mail server on port", PORT)

	URL := ""

	if runtime.GOOS == "windows" {
		URL = "localhost" + PORT
	} else {
		URL = PORT
	}
	srv := &http.Server{
		Addr:    URL,
		Handler: app.routes(),
	}

	// start the server
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromAddress: os.Getenv("MAIL_FROM_ADDR"),
		FromName:    os.Getenv("MAIL_FROM_NAME"),
	}

	return m
}
