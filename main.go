package main

import (
	"github.com/projectz-ro/journalz-ro/commands"
	"github.com/projectz-ro/journalz-ro/config"
	"github.com/projectz-ro/journalz-ro/db"
	"log"
)

// TODO Random reminder function to show a random entry to remind you of it
func main() {
	config.LoadConfig()

	err := db.InitializeDB()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer db.USERDB.Close()

	if err := commands.Execute(); err != nil {
		log.Fatal(err)
	}

}
