package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	fetchData "github.com/18f/cf-events/fetchers"
	helpers "github.com/18f/cf-events/helpers"
	"github.com/robfig/cron"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Env struct {
	// A struct for storing a pointer to a mongo collection
	Collection *mgo.Collection
}

type Result struct {
	// Quota meta data struct returned from the CF api
	Id      string  `json:"_id" bson:"_id"`
	Average float32 `json:"average" bson:"average"`
	Max     float32 `json:"max" bson:"max"`
	Min     float32 `json:"min" bson:"min"`
}

func (env *Env) statsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the average, min, and max of events grouped by actee types for each day.
	var results []Result
	pipe := env.Collection.Pipe([]bson.M{bson.M{"$group": bson.M{"_id": bson.M{"actee": "$eventresource.entity.actee", "actee_type": "$eventresource.entity.actee_type", "day": bson.M{"$substr": []interface{}{"$eventresource.metadata.created_at", 0, 10}}}, "events": bson.M{"$sum": 1}}}, bson.M{"$group": bson.M{"_id": "$_id.actee_type", "average": bson.M{"$avg": "$events"}, "max": bson.M{"$max": "$events"}, "min": bson.M{"$min": "$events"}}}})
	err := pipe.All(&results)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		json.NewEncoder(w).Encode(results)
	}

}

func main() {
	log.Printf("Starting app on :%s", os.Getenv("PORT"))

	// Start cron job for every day
	c := cron.New()
	c.AddFunc("0 0 0 * * *", func() {
		log.Println("Starting Collection")
		eventsCount := fetchData.Events()
		log.Println("Event Collection Complete. Events Collected: ", eventsCount)
	})
	c.Start()

	// Get Mongo session
	session := helpers.GetMongoSession()
	defer session.Close()
	collection := session.DB("cloudfoundry").C("events")
	env := &Env{Collection: collection}

	// Start server
	http.HandleFunc("/", env.statsHandler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
