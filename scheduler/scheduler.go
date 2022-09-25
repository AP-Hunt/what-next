package scheduler

import (
	"time"

	"github.com/AP-Hunt/what-next/m/calendar"
	"github.com/AP-Hunt/what-next/m/todo"
	ical "github.com/arran4/golang-ical"
	"golang.org/x/exp/slices"
)

type Schedule struct {
	CurrentCalendarEvents      []*ical.VEvent
	NextCalendarEvents         []*ical.VEvent
	TimeUntilNextCalendarEvent *time.Duration
	AchievableTasks            todo.TodoItemCollection
}

// GenerateSchedule takes todays calendar and a todo list
// and produces a schedule which shows
// * any calendar events currently happening
// * the next calendar event after that
// * the time until that event
// * which tasks from the todo list are achievable in that time
func GenerateSchedule(now time.Time, calendars []*ical.Calendar, todoList *todo.TodoItemCollection) (*Schedule, error) {
	schedule := &Schedule{
		CurrentCalendarEvents:      []*ical.VEvent{},
		NextCalendarEvents:         []*ical.VEvent{},
		TimeUntilNextCalendarEvent: nil,
		AchievableTasks:            todo.TodoItemCollection{},
	}

	allEvents := []*ical.VEvent{}
	for _, cal := range calendars {
		for _, e := range cal.Events() {
			allEvents = append(allEvents, e)
		}
	}

	for _, event := range allEvents {
		isHappening, err := calendar.EventIsCurrentlyHappening(event, now)

		if err != nil {
			return nil, err
		}

		if isHappening {
			schedule.CurrentCalendarEvents = append(schedule.CurrentCalendarEvents, event)
		}
	}

	eventsStartingAfterNow := []*ical.VEvent{}
	for _, event := range allEvents {
		start, err := event.GetStartAt()
		if err != nil {
			return nil, err
		}

		if start.After(now) {
			eventsStartingAfterNow = append(eventsStartingAfterNow, event)
		}
	}

	tasksForConsideration := todoList.Filter(func(ti *todo.TodoItem) bool {
		return ti.Completed == false &&
			(ti.DueDate == nil || (ti.DueDate != nil && ti.DueDate.After(now)))
	})

	tasksWithoutDurationSet := *tasksForConsideration.Filter(func(ti *todo.TodoItem) bool {
		return ti.Duration == nil
	})

	if len(eventsStartingAfterNow) > 0 {
		err := calendar.SortEventsByStartDateAscending(eventsStartingAfterNow)
		if err != nil {
			return nil, err
		}
		nextStartingEvent := eventsStartingAfterNow[0]
		nextEventStartTime, err := nextStartingEvent.GetStartAt()
		if err != nil {
			return nil, err
		}

		for _, event := range eventsStartingAfterNow {
			start, err := event.GetStartAt()
			if err != nil {
				return nil, err
			}

			if start.Equal(nextEventStartTime) {
				schedule.NextCalendarEvents = append(schedule.NextCalendarEvents, event)
			}
		}

		timeUntilNextEvent := nextEventStartTime.Sub(now)
		schedule.TimeUntilNextCalendarEvent = &timeUntilNextEvent

		achievableTasksWithinDuration := *tasksForConsideration.Filter(func(ti *todo.TodoItem) bool {
			return ti.Duration != nil && *ti.Duration <= timeUntilNextEvent
		})

		schedule.AchievableTasks = *achievableTasksWithinDuration.Append(&tasksWithoutDurationSet)
	} else {
		schedule.AchievableTasks = *tasksForConsideration
	}

	return schedule, nil
}

func sortEventsByStartTimeAsc(events []*ical.VEvent) error {
	var err error = nil
	slices.SortFunc(events, func(a *ical.VEvent, b *ical.VEvent) bool {
		aStart, e := a.GetStartAt()
		if e != nil {
			err = e
			return false
		}

		bStart, e := b.GetStartAt()
		if e != nil {
			err = e
			return false
		}

		return aStart.Before(bStart)
	})

	return err
}
