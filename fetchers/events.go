package fetchData

import (
	"log"
	"strings"
	"time"

	helpers "github.com/18F/cf-events/helpers"
)

func Events() int {
	// Set Event counter
	counter := 0
	// Get Mongo session
	session := helpers.GetMongoSession()
	defer session.Close()
	collection := session.DB("cloudfoundry").C("events")
	// Make new token
	token := helpers.NewToken()
	// events Generator
	eventsGen := token.EventGen()
	// get event indefinitely
	for _ = range time.Tick(5 * time.Second) {
		apiResponse := eventsGen()
		for _, event := range apiResponse.Resources {
			mongoEvent := helpers.MongoEventResource{event.MetaData.Guid, event}
			err := collection.Insert(mongoEvent)
			if err != nil {
				// Break loop only if there is a serious error
				if strings.Contains(err.Error(), "E11000") == false {
					log.Fatal(err.Error())
				}
			}
			counter += 1
		}
		// Break loop if there are no more urls
		if apiResponse.NextUrl == "" {
			break
		}
	}
	return counter
}
