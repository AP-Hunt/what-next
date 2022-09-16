package cmd_test

import (
	"context"
	"time"

	"github.com/AP-Hunt/what-next/m/calendar"
	. "github.com/AP-Hunt/what-next/m/calendar/fakes"
	"github.com/AP-Hunt/what-next/m/cmd"
	commandContext "github.com/AP-Hunt/what-next/m/context"
	"github.com/AP-Hunt/what-next/m/views"
	. "github.com/AP-Hunt/what-next/m/views/fakes"
	ical "github.com/arran4/golang-ical"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Calendar", func() {
	Describe("View", func() {
		var (
			viewEngine      *FakeViewEngineInterface
			calendarService *FakeCalendarServiceInterface
			cmdContext      commandContext.CommandContext
		)

		BeforeEach(func() {
			viewEngine = &FakeViewEngineInterface{}
			calendarService = &FakeCalendarServiceInterface{}

			cmdContext = commandContext.NewCommandContext(context.Background()).
				WithCalendarService(calendarService).
				WithViewEngine(viewEngine)

		})

		It("restricts the calendar entries to those occurring today", func() {
			midnight := time.Now().Truncate(24 * time.Hour)

			cal := ical.NewCalendar()
			evtInsideToday := cal.AddEvent("1")
			evtInsideToday.SetStartAt(midnight.Add(10 * time.Hour))
			evtInsideToday.SetEndAt(midnight.Add(11 * time.Hour))

			evtStartingYesterdayEndingToday := cal.AddEvent("2")
			evtStartingYesterdayEndingToday.SetStartAt(midnight.Add(-2 * time.Hour))
			evtStartingYesterdayEndingToday.SetEndAt(midnight.Add(2 * time.Hour))

			evtStartingTodayEndingTomorrow := cal.AddEvent("3")
			evtStartingTodayEndingTomorrow.SetStartAt(midnight.Add(23 * time.Hour))
			evtStartingTodayEndingTomorrow.SetEndAt(midnight.Add(26 * time.Hour))

			evtEvtStartingAndEndingYesterday := cal.AddEvent("4")
			evtEvtStartingAndEndingYesterday.SetStartAt(midnight.Add(-12 * time.Hour))
			evtEvtStartingAndEndingYesterday.SetEndAt(midnight.Add(-10 * time.Hour))

			evtEvtStartingAndEndingTomorrow := cal.AddEvent("5")
			evtEvtStartingAndEndingTomorrow.SetStartAt(midnight.Add(26 * time.Hour))
			evtEvtStartingAndEndingTomorrow.SetEndAt(midnight.Add(28 * time.Hour))

			calendarService.GetCalendarByDisplayNameReturns(&calendar.CalendarRecord{
				Id:          1,
				DisplayName: "foo",
				URL:         "file://an.ical",
			}, nil)

			calendarService.OpenCalendarReturns(cal, nil)

			PrepareCommandForTest(cmd.CalendarViewCmd, []string{"foo"})

			err := cmd.CalendarViewCmd.ExecuteContext(cmdContext)
			Expect(err).ToNot(HaveOccurred())

			Expect(viewEngine.DrawCallCount()).To(Equal(1))
			view := viewEngine.DrawArgsForCall(0)
			viewData := view.Data().(*views.CalendarViewData)

			passedEvents := viewData.Calendar.Events()
			Expect(passedEvents).To(HaveLen(3))

			selectedIds := []string{}
			for _, evt := range passedEvents {
				selectedIds = append(selectedIds, evt.Id())
			}

			Expect(selectedIds).To(Equal([]string{"1", "2", "3"}))
		})
	})
})
