package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type EventAPIResponse struct {
	// Struct of API response for quota data
	APIResponse
	Resources []EventResource `json:"resources" bson:"resources"`
}

type EventMetaData struct {
	// Quota meta data struct returned from the CF api
	Guid    string `json:"guid" bson:"guid"`
	Url     string `json:"url" bson:"url"`
	Created string `json:"created_at" bson:"created_at"`
	Updated string `json:"updated_at" bson:"updated_at"`
}

type EventEntity struct {
	// Quota entity sturct returned from the CF api
	Type             string `json:"type" bson:"type"`
	Actor            string `json:"actor" bson:"actor"`
	ActorType        string `json:"actor_type" bson:"actor_type"`
	ActorName        string `json:"actor_name" bson:"actor_name"`
	Actee            string `json:"actee" bson:"actee"`
	ActeeType        string `json:"actee_type" bson:"actee_type"`
	ActeeName        string `json:"actee_name" bson:"actee_name"`
	Timestamp        string `json:"timestamp" bson:"timestamp"`
	Metadata         string `json:"metadata" bson:"metadata"`
	SpaceGuid        string `json:"space_guid" bson:"space_guid"`
	OrganizationGuid string `json:"organization_guid" bson:"organization_guid"`
}

type EventResource struct {
	// Quota resource struct returned from the CF api, composed
	// composed of metadata and entity data.
	MetaData EventMetaData `json:"metadata" bson:"metadata"`
	Entity   EventEntity   `json:"entity" bson:"entity"`
}

type MongoEventResource struct {
	// Struct for making the Id the index, to prevent duplicate
	// events from entering the database.
	Id string `bson:"_id"`
	EventResource
}

func PrepTime() string {
	// This creates a url to collect all the data between 2 days an 2 hours ago
	// and yesterday
	before := time.Now().Add(-24 * time.Hour).Format("2006-01-02T15:04:04Z")
	after := time.Now().Add(-50 * time.Hour).Format("2006-01-02T15:04:04Z")
	return fmt.Sprintf("?q=timestamp>%s&q=timestamp<%s", after, before)
}

func (token *Token) getEvent(eventUrl string) *EventAPIResponse {
	// Get a list of quotas and converts it to the QuotaAPIResponse struct
	req_url := fmt.Sprintf("https://api.%s%s", os.Getenv("API_URL"), eventUrl)
	res := token.make_request(req_url)
	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	var events EventAPIResponse
	if json.Unmarshal(body, &events) == nil {
		fmt.Println("Error")
	}
	return &events
}

func (token *Token) EventGen() func() *EventAPIResponse {
	// Create a generator for the events api
	baseUrl := "/v2/events" + PrepTime()
	currentPage := 1
	eventUrl := baseUrl + fmt.Sprintf("&page=%d", currentPage)
	return func() *EventAPIResponse {
		fmt.Println(eventUrl)
		eventResponse := token.getEvent(eventUrl)
		currentPage += 1
		eventUrl = baseUrl + fmt.Sprintf("&page=%d", currentPage)
		return eventResponse
	}
}
