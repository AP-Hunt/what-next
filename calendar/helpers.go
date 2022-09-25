package calendar

import (
	"fmt"
	"time"

	ical "github.com/arran4/golang-ical"
	"golang.org/x/exp/slices"
)

var midnightToday = time.Now().Truncate(24 * time.Hour)
var midnightTomorrow = midnightToday.Add(24 * time.Hour)

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
	var err error = nil
	start, startErr := evt.GetStartAt()
	if startErr != nil {
		err = fmt.Errorf("cannot fetch start time: %s", startErr)
	}

	end, endErr := evt.GetEndAt()
	if endErr != nil {
		err = fmt.Errorf("%s\ncannot fetch end time: %s", err, endErr)
	}

	return start, end, err
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
