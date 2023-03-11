package main

import (
	"auth/data"
	"database/sql"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const PORT = ":80"

var counts int64

// type Config struct {
// 	DB     *sql.DB
// 	Models data.Models
// }

type AppConfig struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting auth service")

	// connect db

	conn := connectToDB()
	if conn == nil {
		log.Panic("Cant connect to postgres")
	}

	// setup config

	app := AppConfig{
		DB:     conn,
		Models: data.New(conn),
	}

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

	log.Println("Starting auth service int port ", PORT)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready....")
			counts++
		} else {
			log.Println("Connected to postgres")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2s")
		time.Sleep(2 * time.Second)
		continue
	}
}
