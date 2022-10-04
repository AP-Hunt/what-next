package scheduler_test

import (
	"time"

	"github.com/AP-Hunt/what-next/m/calendar"
	"github.com/AP-Hunt/what-next/m/scheduler"
	"github.com/AP-Hunt/what-next/m/todo"
	ical "github.com/arran4/golang-ical"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-module/carbon/v2"
	"github.com/hako/durafmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func newEvent(now time.Time, relativeStartTime string, duration string) *ical.VEvent {
	startTimeDurationOffset, err := durafmt.ParseString(relativeStartTime)
	if err != nil {
		panic(err)
	}

	eventDuration, err := durafmt.ParseString(duration)
	if err != nil {
		panic(err)
	}

	event := ical.NewEvent(gofakeit.UUID())
	event.SetProperty(ical.ComponentProperty(ical.PropertySummary), gofakeit.Phrase())

	event.SetStartAt(now.Add(startTimeDurationOffset.Duration()))
	event.SetDuration(eventDuration.Duration())

	return event
}

func randomEventNotHappeningNow(now time.Time) *ical.VEvent {
	midnightOfNow := now.Truncate(24 * time.Hour)
	endOfDay := midnightOfNow.Add(24 * time.Hour)

	var event *ical.VEvent
	happeningNow := true

	// Keep looping until you get a time that's not now
	for happeningNow {
		start := gofakeit.DateRange(
			midnightOfNow,
			endOfDay,
		)

		end := gofakeit.DateRange(
			start,
			endOfDay,
		)

		event = ical.NewEvent(gofakeit.UUID())
		event.SetProperty(ical.ComponentProperty(ical.PropertySummary), gofakeit.Phrase())
		event.SetStartAt(start)
		event.SetEndAt(end)

		happening, err := calendar.EventIsCurrentlyHappening(event, now)
		if err != nil {
			panic(err)
		}

		happeningNow = happening
	}

	return event
}

func generateCalendar(events ...*ical.VEvent) *ical.Calendar {
	cal := ical.NewCalendar()
	for _, event := range events {
		cal.AddVEvent(event)
	}

	return cal
}

func taskWithDuration(duration time.Duration) *todo.TodoItem {
	return &todo.TodoItem{
		Id:        gofakeit.Number(0, 100),
		Action:    gofakeit.Phrase(),
		DueDate:   nil,
		Duration:  &duration,
		Completed: false,
	}
}

func taskWithoutDuration() *todo.TodoItem {
	return &todo.TodoItem{
		Id:        gofakeit.Number(0, 100),
		Action:    gofakeit.Phrase(),
		DueDate:   nil,
		Duration:  nil,
		Completed: false,
	}
}

var _ = Describe("Scheduler", func() {
	now := time.Now()
	Describe("GenerateSchedule", func() {
		Context("when 'now' falls within a calendar event", func() {
			It("the schedule contains that event in the CurrentCalendarEvents field", func() {
				currentEvent := newEvent(now, "-15m", "30m")
				cal := generateCalendar(
					newEvent(now, "-2h", "1h"),
					newEvent(now, "-1h", "30m"),
					currentEvent,
					newEvent(now, "1h", "15m"),
				)

				todoList := todo.NewTodoItemCollection([]*todo.TodoItem{})

				schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{cal}, todoList)
				Expect(err).ToNot(HaveOccurred())

				Expect(schedule.CurrentCalendarEvents).To(ContainElement(currentEvent))
			})

			It("the schedule contains all events that are currently happneing in the CurrentCalendarEvents field", func() {
				eventOne := newEvent(now, "-40m", "60m")
				eventTwo := newEvent(now, "-5m", "10m")

				cal := generateCalendar(
					randomEventNotHappeningNow(now),
					randomEventNotHappeningNow(now),
					randomEventNotHappeningNow(now),
					eventOne,
					randomEventNotHappeningNow(now),
					randomEventNotHappeningNow(now),
					eventTwo,
					randomEventNotHappeningNow(now),
				)

				todoList := todo.NewTodoItemCollection([]*todo.TodoItem{})

				schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{cal}, todoList)
				Expect(err).ToNot(HaveOccurred())

				Expect(schedule.CurrentCalendarEvents).To(ContainElement(eventOne))
				Expect(schedule.CurrentCalendarEvents).To(ContainElement(eventTwo))
			})
		})

		Context("when there is at least one more meeting starting after now", func() {
			It("the schedule will contain the next event in the NextCalendarEvents field", func() {
				nextEvent := newEvent(now, "1h", "20m")
				cal := generateCalendar(
					newEvent(now, "-1h", "45m"),
					nextEvent,
					newEvent(now, "2h", "30m"),
				)

				todoList := todo.NewTodoItemCollection([]*todo.TodoItem{})

				schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{cal}, todoList)
				Expect(err).ToNot(HaveOccurred())

				Expect(schedule.NextCalendarEvents).To(ContainElement(nextEvent))
			})

			It("the schedule does not contain all events starting after now in the NextCalendarEvents field", func() {
				nextEvent := newEvent(now, "1h", "20m")
				eventAferThat := newEvent(now, "2h", "30m")
				cal := generateCalendar(
					newEvent(now, "-1h", "45m"),
					nextEvent,
					eventAferThat,
				)

				todoList := todo.NewTodoItemCollection([]*todo.TodoItem{})

				schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{cal}, todoList)
				Expect(err).ToNot(HaveOccurred())

				Expect(schedule.NextCalendarEvents).To(ContainElement(nextEvent))
				Expect(schedule.NextCalendarEvents).ToNot(ContainElement(eventAferThat))
			})

			Context("and when two or more events start next and share the same start time", func() {
				It("the schedule will contain both events in the NextCalendarEvents field", func() {
					nextEvent := newEvent(now, "1h", "30m")
					simultaneousEvent := newEvent(now, "1h", "1h")
					cal := generateCalendar(
						newEvent(now, "-1h", "45m"),
						nextEvent,
						simultaneousEvent,
					)

					todoList := todo.NewTodoItemCollection([]*todo.TodoItem{})

					schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{cal}, todoList)
					Expect(err).ToNot(HaveOccurred())

					Expect(schedule.NextCalendarEvents).To(ContainElement(nextEvent))
					Expect(schedule.NextCalendarEvents).To(ContainElement(simultaneousEvent))

				})
			})

			Context("but the next one starts tomorrow", func() {
				It("the schedule will not contain any calendar events in the NextCalendarEvents field and the TimeUntilNextCalendarEvent field will be nil", func() {
					cal := generateCalendar(
						newEvent(now, "36h", "1h"),
					)

					todoList := todo.NewTodoItemCollection([]*todo.TodoItem{})

					schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{cal}, todoList)
					Expect(err).ToNot(HaveOccurred())

					Expect(schedule.NextCalendarEvents).To(BeEmpty())
					Expect(schedule.TimeUntilNextCalendarEvent).To(BeNil())
				})
			})

			It("the schedule will contain the duration of time until the next calendar event in the TimeUntilNextCalendarEvent field", func() {
				nextEvent := newEvent(now, "4h", "20m")
				cal := generateCalendar(
					newEvent(now, "-1h", "45m"),
					nextEvent,
				)

				todoList := todo.NewTodoItemCollection([]*todo.TodoItem{})

				schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{cal}, todoList)
				Expect(err).ToNot(HaveOccurred())

				Expect(schedule.TimeUntilNextCalendarEvent).ToNot(BeNil())

				// Expect the duration to match within 30s
				// because there are always likely to be tiny
				// time differences when measuring in microseconds
				Expect(int(*schedule.TimeUntilNextCalendarEvent)).To(
					BeNumerically("==",
						int(4*time.Hour),
						int(30*time.Second),
					))
			})
		})

		Context("when there are no more calendar events in the day", func() {
			It("the schedule's NextCalendarEvents field will be empty", func() {
				cal := generateCalendar(
					newEvent(now, "-6h", "30m"),
					newEvent(now, "-4h30m", "30m"),
					newEvent(now, "-30m", "10m"),
				)

				todoList := todo.NewTodoItemCollection([]*todo.TodoItem{})

				schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{cal}, todoList)
				Expect(err).ToNot(HaveOccurred())

				Expect(schedule.NextCalendarEvents).To(HaveLen(0))
			})

			It("the schedules TimeUntilNextCalendarEvent field will be nil", func() {
				cal := generateCalendar(
					newEvent(now, "-6h", "30m"),
					newEvent(now, "-4h30m", "30m"),
					newEvent(now, "-30m", "10m"),
				)

				todoList := todo.NewTodoItemCollection([]*todo.TodoItem{})

				schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{cal}, todoList)
				Expect(err).ToNot(HaveOccurred())

				Expect(schedule.TimeUntilNextCalendarEvent).To(BeNil())
			})
		})

		Describe("the AchievableTasks field", func() {
			cal := generateCalendar(
				newEvent(now, "2h", "30m"),
			)

			It("will not contain completed tasks", func() {
				tomorrow := carbon.Tomorrow().Carbon2Time()
				otherwiseAchievableTask := taskWithDuration(30 * time.Minute)
				otherwiseAchievableTask.Completed = true
				otherwiseAchievableTask.DueDate = &tomorrow

				tasks := todo.NewTodoItemCollection([]*todo.TodoItem{
					otherwiseAchievableTask,
				})

				schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{cal}, tasks)
				Expect(err).ToNot(HaveOccurred())

				Expect(schedule.AchievableTasks.Enumerate()).ToNot(ContainElement(otherwiseAchievableTask))
			})

			It("will contain tasks whose expected duration is smaller than the amount of time until the next calendar event", func() {
				taskTooLong := taskWithDuration(4 * time.Hour)
				taskShortEnough := taskWithDuration(40 * time.Minute)

				tasks := todo.NewTodoItemCollection([]*todo.TodoItem{
					taskTooLong,
					taskShortEnough,
				})

				schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{cal}, tasks)
				Expect(err).ToNot(HaveOccurred())

				Expect(schedule.AchievableTasks.Enumerate()).To(ContainElement(taskShortEnough))
				Expect(schedule.AchievableTasks.Enumerate()).ToNot(ContainElement(taskTooLong))
			})

			It("will contain tasks which don't have an expected duration set", func() {
				taskHasNoDuration := taskWithoutDuration()
				taskShortEnough := taskWithDuration(40 * time.Minute)

				tasks := todo.NewTodoItemCollection([]*todo.TodoItem{
					taskHasNoDuration,
					taskShortEnough,
				})

				schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{cal}, tasks)
				Expect(err).ToNot(HaveOccurred())

				Expect(schedule.AchievableTasks.Enumerate()).To(ContainElement(taskHasNoDuration))
			})

			It("will only contain tasks whose due date is in the future, if one is set", func() {
				thePast := time.Date(2020, 01, 01, 00, 00, 00, 00, time.Local)
				theFuture := now.AddDate(0, 0, 1)

				taskWithDueDateInThePast := taskWithoutDuration()
				taskWithDueDateInThePast.DueDate = &thePast

				taskWithDueDateInTheFuture := taskWithoutDuration()
				taskWithDueDateInTheFuture.DueDate = &theFuture

				tasks := todo.NewTodoItemCollection([]*todo.TodoItem{
					taskWithDueDateInThePast,
					taskWithDueDateInTheFuture,
				})

				schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{cal}, tasks)
				Expect(err).ToNot(HaveOccurred())

				Expect(schedule.AchievableTasks.Enumerate()).To(ContainElement(taskWithDueDateInTheFuture))
				Expect(schedule.AchievableTasks.Enumerate()).ToNot(ContainElement(taskWithDueDateInThePast))
			})

			It("will order tasks by due date ascending, with tasks without due dates at the end", func() {
				cal := generateCalendar(
					newEvent(now, "-1h", "2h"),
				)

				aWeekToday := now.Add(7 * 24 * time.Hour)
				taskA := &todo.TodoItem{
					Id:        1,
					Action:    "A",
					DueDate:   &aWeekToday,
					Duration:  nil,
					Completed: false,
				}

				twoWeeksToday := now.Add(14 * 24 * time.Hour)
				taskB := &todo.TodoItem{
					Id:        2,
					Action:    "B",
					DueDate:   &twoWeeksToday,
					Duration:  nil,
					Completed: false,
				}

				threeWeeksToday := now.Add(21 * 24 * time.Hour)
				taskC := &todo.TodoItem{
					Id:        3,
					Action:    "C",
					DueDate:   &threeWeeksToday,
					Duration:  nil,
					Completed: false,
				}

				tasks := todo.NewTodoItemCollection([]*todo.TodoItem{taskB, taskC, taskA})

				schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{cal}, tasks)
				Expect(err).ToNot(HaveOccurred())

				scheduledTasks := schedule.AchievableTasks.Enumerate()
				actions := []string{}
				for _, task := range scheduledTasks {
					actions = append(actions, task.Action)
				}
				Expect(actions).To(Equal([]string{"A", "B", "C"}))
			})

			Context("when there are no more meetings in the day", func() {
				calWithoutFutureEvents := generateCalendar(
					newEvent(now, "-2h", "30m"),
				)

				It("will contain all tasks which have no due date set, or a due date in the future, without consideration of duration", func() {
					thePast := time.Date(2020, 01, 01, 00, 00, 00, 00, time.Local)
					theFuture := now.AddDate(0, 0, 1)

					taskWithDueDateInThePast := taskWithoutDuration()
					taskWithDueDateInThePast.DueDate = &thePast

					taskWithDueDateInTheFuture := taskWithoutDuration()
					taskWithDueDateInTheFuture.DueDate = &theFuture

					taskWithNoDueDate := taskWithoutDuration()

					tasks := todo.NewTodoItemCollection([]*todo.TodoItem{
						taskWithDueDateInThePast,
						taskWithDueDateInTheFuture,
						taskWithNoDueDate,
					})

					schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{calWithoutFutureEvents}, tasks)
					Expect(err).ToNot(HaveOccurred())

					Expect(schedule.AchievableTasks.Enumerate()).To(ContainElement(taskWithDueDateInTheFuture))
					Expect(schedule.AchievableTasks.Enumerate()).ToNot(ContainElement(taskWithDueDateInThePast))
					Expect(schedule.AchievableTasks.Enumerate()).To(ContainElement(taskWithNoDueDate))
				})
			})
		})

		It("CurrentCalendarEvents and NextCalendarEvents fields will contain events from all input calendars", func() {
			currentEventInCalA := newEvent(now, "-30m", "1h")
			currentEventInCalB := newEvent(now, "-60m", "1h30m")
			nextEventInCalA := newEvent(now, "2h", "10m")
			nextEventInCalB := newEvent(now, "2h", "30m")

			calA := generateCalendar(currentEventInCalA, nextEventInCalA)
			calB := generateCalendar(currentEventInCalB, nextEventInCalB)

			tasks := todo.NewTodoItemCollection([]*todo.TodoItem{})

			schedule, err := scheduler.GenerateSchedule(now, []*ical.Calendar{calA, calB}, tasks)
			Expect(err).ToNot(HaveOccurred())

			Expect(schedule.CurrentCalendarEvents).To(ContainElements(currentEventInCalA, currentEventInCalB))
			Expect(schedule.NextCalendarEvents).To(ContainElements(nextEventInCalA, nextEventInCalB))
		})
	})
})
