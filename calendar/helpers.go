package calendar

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/golang-module/carbon/v2"
	"golang.org/x/exp/slices"
)

var midnightToday = time.Now().Truncate(24 * time.Hour)
var midnightTomorrow = midnightToday.Add(24 * time.Hour)

// Date format as defined in RFC 5545
// https://www.rfc-editor.org/rfc/rfc5545#section-3.3.4
var RegexIcalDate *regexp.Regexp = regexp.MustCompile("^[0-9]{4}(0[1-9]|1[0-2])([0-2][0-9]|3[0-1])$")

func EventStartsToday(evt *ical.VEvent) (bool, error) {
	start, err := evt.GetStartAt()
	if err != nil {
		return false, err
	}

	startsAtMidnight := start.Equal(midnightToday)
	startsDuringToday := start.After(midnightToday) && start.Before(midnightTomorrow)
	return (startsAtMidnight || startsDuringToday), nil
}

func EventEndsToday(evt *ical.VEvent) (bool, error) {
	end, err := evt.GetEndAt()
	if err != nil {
		return false, err
	}

	endsAtMidnight := end.Equal(midnightToday)
	endsDuringToday := end.After(midnightToday) && end.Before(midnightTomorrow)
	return (endsAtMidnight || endsDuringToday), nil
}

func EventIsCurrentlyHappening(evt *ical.VEvent, now time.Time) (bool, error) {
	start, end, err := EventStartAndEnd(evt)
	if err != nil {
		return false, err
	}

	if now.Equal(start) || now.Equal(end) {
		return true, nil
	}

	if start.Before(now) && end.After(now) {
		return true, nil
	}

	return false, nil
}

func EventStartAndEnd(evt *ical.VEvent) (time.Time, time.Time, error) {
	tZero := time.Time{}
	DTSTART := evt.GetProperty(ical.ComponentPropertyDtStart)
	DTEND := evt.GetProperty(ical.ComponentPropertyDtEnd)
	DURATION := evt.GetProperty("DURATION")

	if DTSTART == nil {
		// -DSTART
		return tZero, tZero, fmt.Errorf("no start time specified")
	} else if DTSTART != nil && DTEND != nil {
		// DTSTART + DTEND

		var err error = nil
		start, startErr := evt.GetStartAt()
		if startErr != nil {
			err = fmt.Errorf("cannot fetch start time: %s", startErr)
		}

		end, endErr := evt.GetEndAt()
		if endErr != nil {
			// The specification in RFC 5545 says that an event
			// with a DTSTART of type DATE-TIME and no DTEND
			// should be interpreted as starting and ending at
			// the same time
			// https://www.rfc-editor.org/rfc/rfc5545#section-3.6.1
			end = start
		}

		return start, end, err
	} else if DTSTART != nil && (DTEND == nil && DURATION == nil) {
		// DTSRT - DTEND - DURATION

		startIsDate := RegexIcalDate.Match([]byte(DTSTART.Value))
		if startIsDate {
			start, err := evt.GetAllDayStartAt()
			if err != nil {
				return tZero, tZero, err
			}

			// When DTSTART is present with type of DATE
			// but there is no DTEND or DURATION, the end
			// time is the start of the next day
			end := carbon.Time2Carbon(start).Tomorrow().StartOfDay().Carbon2Time()
			return start, end, nil
		} else {
			start, err := evt.GetStartAt()
			if err != nil {
				return tZero, tZero, err
			}

			// When DTSTART is present with type of DATETIME
			// but there is no DTEND or DURATION, the end time
			// is the same as the start time
			return start, start, nil
		}
	} else if DTSTART != nil && DURATION != nil {
		// DTSTART + DURATION

		duration, err := time.ParseDuration(strings.ToLower(DURATION.Value))
		if err != nil {
			return tZero, tZero, fmt.Errorf("parsing duration: %s", err)
		}

		start, err := evt.GetStartAt()
		if err != nil {
			return tZero, tZero, err
		}

		return start, start.Add(duration), nil
	}

	return tZero, tZero, nil
}

func IsAllDayEvent(evt *ical.VEvent) (bool, error) {
	DTSTART := evt.GetProperty(ical.ComponentPropertyDtStart)
	if DTSTART == nil {
		return false, fmt.Errorf("cannot find %s property", ical.ComponentPropertyDtStart)
	}

	DTEND := evt.GetProperty(ical.ComponentPropertyDtEnd)
	DURATION := evt.GetProperty("DURATION")

	if DTEND != nil {
		startMatches := RegexIcalDate.Match([]byte(DTSTART.Value))
		endMatches := RegexIcalDate.Match([]byte(DTEND.Value))

		return startMatches && endMatches, nil
	} else {
		// DURATION + DTSTART replaces DTSTART + DTEND
		if DURATION != nil {
			return false, nil
		}

		// DTSTART with no DTEND or DURATION is an all day event
		return true, nil
	}
}

func SortEventsByStartDateAscending(events []*ical.VEvent) error {
	var err error = nil
	slices.SortFunc(events, func(evtA *ical.VEvent, evtB *ical.VEvent) bool {

		aStart, aEnd, err := EventStartAndEnd(evtA)
		if err != nil {
			err = fmt.Errorf("sorting calendar entries: %s", err)
			return false
		}

		bStart, bEnd, err := EventStartAndEnd(evtB)
		if err != nil {
			err = fmt.Errorf("sorting calendar entries: %s", err)
			return false
		}

		if aStart.Equal(bStart) {
			aDuration := aEnd.Sub(aStart)
			bDuration := bEnd.Sub(bStart)

			if aDuration == bDuration || aDuration > bDuration {
				return true
			} else {
				return false
			}
		} else if aStart.Before(bStart) {
			return true
		}

		return false
	})

	return err
}
