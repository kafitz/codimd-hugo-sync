package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kafitz/codimd-hugo-sync/database"
	"github.com/kafitz/codimd-hugo-sync/hugoio"
	"github.com/kafitz/codimd-hugo-sync/models"
	"github.com/lib/pq"
)

func loadConfig(fp string) (cfg models.Config, err error) {
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
		msg, err := database.WaitForNotification(listener)
		if err != nil {
			log.Fatalln(err)
		}
		// check for empty notification message resulting from keepalive ping
		if msg == nil {
			continue
		}

		postCreatedAt, err := time.Parse(time.RFC3339, msg.Data["createdAt"].(string))
		if err != nil {
			log.Fatalf("Could not parse timestamp (created): %s", msg.Data["createdAt"])
		}
		// postUpdatedAt, err := time.Parse(time.RFC3339, msg.Data["updatedAt"].(string))
		// if err != nil {
		// 	log.Fatalf("Could not parse timestamp (updated): %s", msg.Data["updatedAt"])
		// }
		postID := msg.Data["shortid"].(string)
		postContent := msg.Data["content"].(string)
		postPermission := msg.Data["permission"].(string) == "locked"

		metadata := models.PostMetadata{
			Date:   postCreatedAt,
			PostID: postID}

		hugoio.UpdateMetadataFromMarkdown(&postContent, &metadata)
		if metadata.IsBlogPost != "true" || !postPermission {
			continue
		}

		updatedPost := hugoio.WriteMetadataToMarkdown(postContent, &metadata)
		err = hugoio.WritePostFile(cfg.Hugo.PostsDir, updatedPost, &metadata)
		if err != nil {
			log.Fatalln(err)
		}

		hugoio.RefreshHugoBlog(cfg.Hugo.BaseDir)
	}
}
