package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/apognu/gocal"
	"github.com/rickar/cal/v2"
)

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Print(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timezone, _ := time.LoadLocation("Europe/Copenhagen")

	cal.DefaultLoc = timezone

	bc := cal.NewBusinessCalendar()
	bc.SetWorkHours(7*time.Hour+59*time.Minute, 12*time.Hour)

	startDate := time.Date(2021, 1, 4, 22, 0, 0, 0, time.UTC)

	if startDate.Before(time.Now()) {
		startDate = time.Now()
	}

	workdayStart := cal.ReplaceLocation(bc.NextWorkdayStart(startDate), timezone)

	f, _ := os.Open("calendar.ics")
	defer f.Close()

	start, end := workdayStart, workdayStart.Add(12*time.Hour)

	c := gocal.NewParser(f)
	c.Start, c.End = &start, &end
	c.Parse()

	sort.Slice(c.Events, func(i, j int) bool {
		return c.Events[i].Start.Before(*c.Events[j].Start)
	})

	skema := []string{}

	for _, e := range c.Events {
		if skipEvent(e) {
			continue
		}

		skema = append(skema, fmt.Sprintf("<say-as interpret-as=\\\"time\\\" format=\\\"hm\\\">%s</say-as>: %s med %s",
			strings.TrimPrefix(e.Start.In(timezone).Format("15:04"), "0"),
			expandSummary(e.Summary),
			organizer(e.Organizer.Cn),
		))
	}

	response(w, "<speak>"+strings.Join(skema, ".<break time=\\\"1s\\\"/>\\n")+"</speak>")
}

func organizer(organizer string) string {
	reg := regexp.MustCompile(`\(.*\)?`)
	name := reg.ReplaceAllString(organizer, "")
	name = strings.TrimSpace(strings.Trim(name, "\""))

	if name == "Lisbet Merlung" {
		return "Flet"
	}

	name = strings.Split(name, " ")[0]

	return name
}

func skipEvent(e gocal.Event) bool {
	if e.Summary == "GS" {
		return true
	}

	return false
}

func expandSummary(summary string) string {
	switch summary {
	case "BIL":
		return "Billedkunst"

	case "DAN":
		return "Dansk"

	case "ENG":
		return "Engelsk"

	case "IDR":
		return "Idræt"

	case "KRI":
		return "Kristendomskundskab"

	case "N/T":
		return "Natur og teknik"

	case "MAT":
		return "Matematik"

	case "MUS":
		return "Musik"

	case "UUV":
		return "<say-as interpret-as=\\\"characters\\\">UUV</say-as>"

	default:
		return summary
	}
}
