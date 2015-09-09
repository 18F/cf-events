package main

import (
	"log"
	"time"

	fetchData "github.com/18f/cf-events/fetchers"
)

func main() {
	// Collect events every hour
	for _ = range time.Tick(360 * time.Second) {
		eventsCount := fetchData.Events()
		log.Println("Event Collection Complete. Events Collected: ", eventsCount)
	}
}
