package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

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

	bc := cal.NewBusinessCalendar()

	startDate := time.Date(2021, 1, 4, 22, 0, 0, 0, time.UTC)

	if startDate.Before(time.Now()) {
		startDate = time.Now()
	}

	workdayStart := bc.NextWorkdayStart(startDate)
	f, _ := os.Open("calendar.ics")
	defer f.Close()

	start, end := workdayStart, workdayStart.Add(24*time.Hour)

	c := gocal.NewParser(f)
	c.Start, c.End = &start, &end
	c.Parse()

	skema := []string{}

	for _, e := range c.Events {
		if skipEvent(e) {
			continue
		}

		skema = append(skema, fmt.Sprintf("%s med %s", expandSummary(e.Summary), organizer(e.Organizer.Cn)))
	}

	response(w, strings.Join(skema, ", "))
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

	case "KRI":
		return "Kristendomskundskab"

	case "N/T":
		return "Natur og teknik"

	case "MAT":
		return "Matematik"

	case "MUS":
		return "Musik"

	default:
		return summary
	}
}
