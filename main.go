package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	BaseURL                 = "https://ground.brpsystems.com/brponline/api/ver3/apps/390/groupactivities"
	Facility_GroenloekkeVej = "4"
)

type LocationData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DurationData struct {
	End   string `json:"end"`
	Start string `json:"start"`
}

type ReturnData struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Duration  DurationData   `json:"duration"`
	Locations []LocationData `json:"locations"`
}

func Fetch(link string) ([]ReturnData, error) {
	res, err := http.Get(link)
	if err != nil {
		println(err)
	}

	var data []ReturnData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func main() {
	baseUrl := "https://ground.brpsystems.com/brponline/api/ver3/apps/390/groupactivities"
	url, err := url.Parse(baseUrl)
	if err != nil {
		println(err)
	}
	q := url.Query()
	q.Set("businessUnit", Facility_GroenloekkeVej)
	q.Set("groupActivitiesFor", "web")
	q.Set("period.end", "2025-09-28T21:59:59.999Z")
	q.Set("period.start", "2025-09-22T17:55:27.866Z")
	q.Set("webCategory", "2")
	url.RawQuery = q.Encode()

	classes, err := Fetch(url.String())
	if err != nil {
		panic(err)
	}

	for _, c := range classes {
		if c.Locations[0].Name != "Funktionelt Område" {
			continue
		}
		fmt.Printf("ID: %d | Name: %s | Place: %s | Start: %s | End: %s\n",
			c.ID, c.Name, c.Locations[0].Name, c.Duration.Start, c.Duration.End)
	}
}
