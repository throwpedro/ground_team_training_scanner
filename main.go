package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	BaseURL                 = "https://ground.brpsystems.com/brponline/api/ver3/apps/390/groupactivities"
	Facility_GroenloekkeVej = "4"
)

type Dates struct {
	Start string
	End   string
}

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

	dates := BuildWeekDates()

	q := url.Query()
	q.Set("businessUnit", Facility_GroenloekkeVej)
	q.Set("groupActivitiesFor", "web")
	q.Set("period.end", dates.End)
	q.Set("period.start", dates.Start)
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

func BuildWeekDates() Dates {
	now := time.Now().UTC()

	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	end := start.AddDate(0, 0, 6).
		Add(time.Hour*23 + time.Minute*59 + time.Second*59 + time.Millisecond*999)

	layout := "2006-01-02T15:04:05.000Z"
	periodStart := start.Format(layout)
	periodEnd := end.Format(layout)

	returnVal := Dates{
		Start: periodStart,
		End:   periodEnd,
	}

	return returnVal
}
