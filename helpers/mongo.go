package helpers

import (
	"fmt"
	"log"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"gopkg.in/mgo.v2"
)

func GetMongoSession() *mgo.Session {
	var mongoURI string
	// Get the VCAP env
	appEnv, _ := cfenv.Current()
	if appEnv != nil {
		// Get services
		services, _ := appEnv.Services.WithTag("mongodb")
		// Get mongo service info
		mongoURI = fmt.Sprint(services[0].Credentials["uri"])
	} else {
		log.Println("Defaulting to local mongo database")
		mongoURI = "mongodb://localhost"
	}
	// Setup session
	s, err := mgo.Dial(mongoURI)
	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s
}
