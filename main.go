package main

import (
	"database/sql"
	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer func(todoDB *sql.DB) {
		err := todoDB.Close()
		if err != nil {
			log.Println(err)
		}
	}(todoDB)

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB)
	// TODO: サーバーをlistenする
	err = http.ListenAndServe(port, mux)
	if err != nil {
		return err
	}
	return nil
}
