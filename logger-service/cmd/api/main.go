package main

import (
	"context"
	"log"
	"logger-service/data"
	"net/http"
	"runtime"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	PORT     = ":80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo

	mongoClient, err := connectToMongo()

	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context to disconnect

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	// close connection

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// go app.serve()
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
	log.Println("starting logger service")
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

// func (app *Config) serve() {
// 	URL := ""

// 	if runtime.GOOS == "windows" {
// 		URL = "localhost" + PORT
// 	} else {
// 		URL = PORT
// 	}
// 	srv := &http.Server{
// 		Addr:    URL,
// 		Handler: app.routes(),
// 	}

// 	// start the server
// 	err := srv.ListenAndServe()

// 	if err != nil {
// 		log.Panic(err)
// 	}
// }

func connectToMongo() (*mongo.Client, error) {
	// connection options

	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect

	c, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println("error connecting :", err)
		return nil, err
	}

	log.Println("connected to mongo")

	return c, nil

}
