package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kafitz/codimd-hugo-sync/models"
	"github.com/lib/pq"
)

// WaitForNotification listens for the Postgres event for changes in the Notes table
func WaitForNotification(l *pq.Listener) (msg *models.NotificationMessage, err error) {
	for {
		select {
		case n := <-l.Notify:
			jsonData := []byte(n.Extra)
			json.Unmarshal(jsonData, &msg)
			return msg, err
		case <-time.After(90 * time.Second):
			fmt.Println("Received no events for 90 seconds, checking connection.")
			go func() {
				l.Ping()
			}()
			return nil, nil
		}
	}
}
