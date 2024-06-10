package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ChrisHeptagon/golibase/models"
	"github.com/ChrisHeptagon/golibase/server"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	os.Setenv("MODE", "DEV")
	os.Setenv("DEV_PORT", "5173")
	log.Println("Starting server...")
	db, err := sql.Open("sqlite3", "file:golibase.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	models.InitializeDatabase(db)
	server.StartServer(db)
}
