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
	"github.com/klauspost/lctime"
	"github.com/rickar/cal/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type lesson struct {
	time       string
	summary    string
	organizers []string
}

func (l lesson) String() string {
	return fmt.Sprintf("<say-as interpret-as=\\\"time\\\" format=\\\"hm\\\">%s</say-as>: %s med %s",
		l.time,
		expandSummary(l.summary),
		organizerJoin(l.organizers),
	)
}

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
	bc.SetWorkHours(7*time.Hour+59*time.Minute, 17*time.Hour)

	workdayStart := cal.ReplaceLocation(bc.NextWorkdayStart(time.Now()), timezone)

	f, _ := os.Open("calendar.ics")
	defer f.Close()

	start, end := workdayStart, workdayStart.Add(12*time.Hour)

	c := gocal.NewParser(f)
	c.Start, c.End = &start, &end
	_ = c.Parse()

	sort.Slice(c.Events, func(i, j int) bool {
		return c.Events[i].Start.Before(*c.Events[j].Start)
	})

	skema := []lesson{}

	var prev *gocal.Event

	for _, e := range c.Events {
		e := e

		// Skip events spanning more than the current day.
		if e.Start.Before(start) || e.End.After(end) {
			continue
		}

		if prev != nil && e.Start.Equal(*prev.Start) {
			skema[len(skema)-1].organizers = append(skema[len(skema)-1].organizers, organizer(e.Organizer.Cn))

			if e.Summary != "GS" && e.Summary != "INK" {
				skema[len(skema)-1].summary = e.Summary
			}

			continue
		}

		skema = append(skema, lesson{
			time:       strings.TrimPrefix(e.Start.In(timezone).Format("15:04"), "0"),
			summary:    e.Summary,
			organizers: []string{organizer(e.Organizer.Cn)},
		})

		prev = &e
	}

	weekday, _ := lctime.StrftimeLoc("da_DK", "%A", workdayStart)

	if len(skema) == 0 {
		response(w, "<speak>Jeg kender ikke skemaet for "+weekday+".</speak>")

		return
	}

	var skemaStrings []string
	for _, s := range skema {
		s := s
		skemaStrings = append(skemaStrings, s.String())
	}

	response(w, "<speak>"+cases.Title(language.Danish).String(weekday)+":<break time=\\\"1s\\\"/>\\n\\nKlokken "+strings.Join(skemaStrings, ".<break time=\\\"1s\\\"/>\\n")+".</speak>")
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

func organizerJoin(organizers []string) string {
	length := len(organizers)
	switch length {
	case 0:
		return ""
	case 1:
		return organizers[0]
	case 2:
		return fmt.Sprintf("%s og %s", organizers[0], organizers[1])
	default:
		return fmt.Sprintf("%s og %s", strings.Join(organizers[0:length-1], ", "), organizers[length-1])
	}
}

func expandSummary(summary string) string {
	re := regexp.MustCompile("^UUV.*")
	summary = re.ReplaceAllString(summary, "UUV")

	switch summary {
	case "BIL":
		return "Billedkunst"

	case "DAN":
		return "Dansk"

	case "ENG":
		return "Engelsk"

	case "HDS":
		return "Håndværk og Design"

	case "HIS":
		return "Historie"

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
