package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// WaitForNotification listens for the Postgres event for changes in the Notes table
func WaitForNotification(l *pq.Listener) {
	for {
		select {
		case n := <-l.Notify:
			fmt.Println("Received data from channel [", n.Channel, "] :")
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(n.Extra), "", "\t")
			if err != nil {
				fmt.Println("Error processing JSON: ", err)
				return
			}
			fmt.Println(string(prettyJSON.Bytes()))
			return
		case <-time.After(90 * time.Second):
			fmt.Println("Received no events for 90 seconds, checking connection.")
			go func() {
				l.Ping()
			}()
			return
		}
	}
}
