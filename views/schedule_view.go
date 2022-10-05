package views

import (
	"fmt"
	"io"

	"github.com/AP-Hunt/what-next/m/calendar"
	"github.com/AP-Hunt/what-next/m/scheduler"
	ical "github.com/arran4/golang-ical"
	"github.com/fatih/color"
	"github.com/hako/durafmt"
)

type ScheduleView struct {
	schedule *scheduler.Schedule
}

func (s *ScheduleView) Draw(out io.Writer) error {
	err := s.drawCurrentMeeting(out)
	if err != nil {
		return err
	}

	err = s.drawNextMeeting(out)
	if err != nil {
		return err
	}

	err = s.drawAchievableTasks(out)
	if err != nil {
		return err
	}

	return nil
}

func (s *ScheduleView) drawCurrentMeeting(out io.Writer) error {
	boldWhite := color.New(color.FgWhite, color.Bold)

	if len(s.schedule.CurrentCalendarEvents) == 0 {
		boldWhite.Fprintln(out, "You are not supposed to be in any meetings right now 🎉️")
	} else {
		boldWhite.Fprintf(out, "You have %d meetings happening now\n", len(s.schedule.CurrentCalendarEvents))
		fmt.Fprintln(out)
		for _, evt := range s.schedule.CurrentCalendarEvents {
			start, end, err := calendar.EventStartAndEnd(evt)
			if err != nil {
				return err
			}

			start = start.Local()
			end = end.Local()

			fmt.Fprintf(
				out,
				"\t%02d%02d-%02d%02d %s\n",
				start.Hour(),
				start.Minute(),
				end.Hour(),
				end.Minute(),
				evt.GetProperty(ical.ComponentProperty(ical.PropertySummary)).Value,
			)
		}
	}
	fmt.Fprintln(out)

	return nil
}

func (s *ScheduleView) drawNextMeeting(out io.Writer) error {
	boldWhite := color.New(color.FgWhite, color.Bold)

	if len(s.schedule.NextCalendarEvents) == 0 {
		fmt.Fprintln(out, boldWhite.Sprintf("You do not have any more meetings today 🎉️"))
		fmt.Fprintln(out)
	} else {
		if len(s.schedule.NextCalendarEvents) == 1 {
			boldWhite.Fprintln(out, boldWhite.Sprint("You have one meeting coming up"))
		} else {
			boldWhite.Fprintf(out, boldWhite.Sprintf("You have %d conflicting meetings starting at the same time coming up\n", len(s.schedule.NextCalendarEvents)))
		}

		fmt.Fprintln(out)

		for _, evt := range s.schedule.NextCalendarEvents {
			start, end, err := calendar.EventStartAndEnd(evt)
			if err != nil {
				return err
			}

			start = start.Local()
			end = end.Local()

			fmt.Fprintf(
				out,
				"\t%02d%02d-%02d%02d %s\n",
				start.Hour(),
				start.Minute(),
				end.Hour(),
				end.Minute(),
				evt.GetProperty(ical.ComponentProperty(ical.PropertySummary)).Value,
			)
		}
	}

	fmt.Fprintln(out)

	return nil
}

func (s *ScheduleView) drawAchievableTasks(out io.Writer) error {
	boldWhite := color.New(color.FgWhite, color.Bold)

	anyTasks := s.schedule.AchievableTasks.Len() > 0

	if s.schedule.TimeUntilNextCalendarEvent == nil {
		if anyTasks {
			boldWhite.Fprintln(out, "In the rest of your day, these are the tasks you could try to complete")
		} else {
			boldWhite.Fprintln(out, "You don't have anything on your todo list for the rest of the day")
			return nil
		}
	} else {

		durationStr := durafmt.Parse(*s.schedule.TimeUntilNextCalendarEvent).LimitFirstN(2)
		if anyTasks {
			boldWhite.Fprintf(out, "In the %s until your next meeting, these are the tasks you could try to complete\n", durationStr)
		} else {
			boldWhite.Fprintf(out, "In the %s until your next meeting, there are no achievable things on you todo list\n", durationStr)
		}
	}
	fmt.Fprintln(out)

	todoListView := TodoListView{}
	todoListView.SetData(&s.schedule.AchievableTasks)

	err := todoListView.Draw(out)
	if err != nil {
		return err
	}

	return nil
}

func (s *ScheduleView) SetData(data interface{}) {
	s.schedule = data.(*scheduler.Schedule)
}

func (s *ScheduleView) Data() interface{} {
	return s.schedule
}
