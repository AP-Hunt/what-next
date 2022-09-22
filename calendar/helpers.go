package calendar

import (
	"time"

	ical "github.com/arran4/golang-ical"
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
	start, err := evt.GetStartAt()
	if err != nil {
		return false, err
	}

	end, err := evt.GetEndAt()
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
