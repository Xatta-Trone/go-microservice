package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"runtime"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const PORT = ":80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	fmt.Println("rabbit mq service")

	// connect to rabbit mq
	rabbitConn, err := connect()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitConn.Close()
	
	app := Config{
		Rabbit: rabbitConn,
	}

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
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}


func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// dont continue until rabbit is ready

	for {

		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")

		if err != nil {
			fmt.Println("Rabbit mq not yet ready")
			counts++
		} else {
			log.Println("connected to rabbitMq")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second

		log.Println("backing off ....")
		time.Sleep(backOff)

		continue
	}

	return connection, nil
}
