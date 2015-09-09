package fetchData

import (
	"fmt"
	"strings"
	"time"

	helpers "github.com/18f/cf-events/helpers"
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
	endLoop := false
	for _ = range time.Tick(2 * time.Second) {
		apiResponse := eventsGen()
		if apiResponse.NextUrl == "" {
			// Break loop if there are no more urls
			break
		}
		for _, event := range apiResponse.Resources {
			mongoEvent := helpers.MongoEventResource{event.MetaData.Guid, event}
			err := collection.Insert(mongoEvent)
			if err != nil {
				// Break loop if there are document exist in the database
				if strings.HasPrefix(err.Error(), "E11000") {
					endLoop = true
					break
				}
			}
			counter += 1
			fmt.Println(counter)
		}
		if endLoop == true {
			break
		}
	}
	return counter
}
