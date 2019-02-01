package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kafitz/codimd-hugo-sync/database"
	"github.com/lib/pq"
)

// Config struct for loading parsed settings.json
type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		DBName   string `json:"dbname"`
		User     string `json:"user"`
		Password string `json:"password"`
	}
}

func loadConfig(fp string) (cfg Config, err error) {
	file, err := os.Open(fp)
	defer file.Close()
	if err != nil {
		return cfg, err
	}
	jsonParser := json.NewDecoder(file)
	jsonParser.Decode(&cfg)
	return cfg, nil
}

func main() {
	cfg, err := loadConfig("./settings.json")
	if err != nil {
		log.Fatalln(err)
	}

	connStr := fmt.Sprintf("host=%s port=%d dbname=%s "+
		"user=%s password=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName,
		cfg.Database.User, cfg.Database.Password)
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("adding trigger")
	database.AddTrigger(db)

	listenerCallback := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(connStr, 5*time.Second, time.Minute, listenerCallback)
	err = listener.Listen("events")
	if err != nil {
		panic(err)
	}

	fmt.Println("Start monitoring PostgreSQL...")
	for {
		database.WaitForNotification(listener)
	}
}
