package routes

import (
	"encoding/json"
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

type ResponseData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	StartTime   string `json:"start_time"`
	StartTimeDK string `json:"start_time_dk"`
	EndTime     string `json:"end_time"`
	EndTimeDK   string `json:"end_time_dk"`
	Area        string `json:"area"`
}

func GetGroundTimes(w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(BaseURL)
	if err != nil {
		println(err)
	}

	buildQuery(url)

	fitnessClasses, err := Fetch(url.String())
	if err != nil {
		panic(err)
	}

	var responseData []ResponseData

	for _, c := range fitnessClasses {
		if c.Locations[0].Name == "Funktionelt Område" || c.Locations[0].ID == 201 {
			responseData = append(responseData, ResponseData{
				ID:          c.ID,
				Name:        c.Name,
				Area:        c.Locations[0].Name,
				StartTime:   c.Duration.Start,
				StartTimeDK: DkTime(c.Duration.Start),
				EndTime:     c.Duration.End,
				EndTimeDK:   DkTime(c.Duration.End),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DkTime(timeStr string) string {
	utcTime, err := time.Parse("2006-01-02T15:04:05.000Z", timeStr)
	if err != nil {
		println(err)
	}
	timeLocation, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		println(err)
	}
	dkTime := utcTime.In(timeLocation)
	return dkTime.Format("2006-01-02 15:04:05.000 MST")
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

func BuildWeekDates() Dates {
	now := time.Now().UTC()

	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Add(-time.Hour * 2)

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

func buildQuery(url *url.URL) {
	dates := BuildWeekDates()

	q := url.Query()
	q.Set("businessUnit", Facility_GroenloekkeVej)
	q.Set("groupActivitiesFor", "web")
	q.Set("period.end", dates.End)
	q.Set("period.start", dates.Start)
	q.Set("webCategory", "2")
	url.RawQuery = q.Encode()
}
