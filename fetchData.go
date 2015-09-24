package main

import (
	"log"
	"net/http"
	"os"

	fetchData "github.com/18f/cf-events/fetchers"
	"github.com/robfig/cron"
)

func main() {
	log.Println("Starting app")
	// Start cron job for every day
	c := cron.New()
	c.AddFunc("0 0 0 * * *", func() {
		log.Println("Starting Collection")
		eventsCount := fetchData.Events()
		log.Println("Event Collection Complete. Events Collected: ", eventsCount)
	})
	c.Start()
	// Start server
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
