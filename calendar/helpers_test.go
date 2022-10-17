package calendar_test

import (
	"fmt"
	"time"

	"github.com/AP-Hunt/what-next/m/calendar"
	ical "github.com/arran4/golang-ical"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func fmtDATETIME(t time.Time) string {
	return fmt.Sprintf(
		"%d%02d%02dT%02d%02d%02d%s",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Format("Z"),
	)
}

func fmtDATE(t time.Time) string {
	return fmt.Sprintf(
		"%0d%02d%02d",
		t.Year(),
		t.Month(),
		t.Day(),
	)
}

var _ = Describe("Helpers", func() {
	midnightToday := time.Now().Truncate(24 * time.Hour)

	Describe("EventStartsToday", func() {
		It("returns false if event starts before midnight today", func() {
			evt := ical.NewEvent("evt")
			evt.SetStartAt(midnightToday.Add(-4 * time.Hour))

			actual, err := calendar.EventStartsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeFalse())
		})

		It("returns false if event starts after midnight tomorrow", func() {
			evt := ical.NewEvent("evt")
			evt.SetStartAt(midnightToday.Add(36 * time.Hour))

			actual, err := calendar.EventStartsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeFalse())
		})

		It("returns false if event starts at midnight tomorrow", func() {
			evt := ical.NewEvent("evt")
			evt.SetStartAt(midnightToday.Add(24 * time.Hour))

			actual, err := calendar.EventStartsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeFalse())
		})

		It("returns true if event starts between midnight today and midnight tomorrow", func() {
			evt := ical.NewEvent("evt")
			evt.SetStartAt(midnightToday.Add(6 * time.Hour))

			actual, err := calendar.EventStartsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeTrue())
		})

		It("returns true if event starts at midnight today", func() {
			evt := ical.NewEvent("evt")
			evt.SetStartAt(midnightToday)

			actual, err := calendar.EventStartsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeTrue())
		})
	})

	Describe("EventEndsToday", func() {
		It("returns false if event ends before midnight today", func() {
			evt := ical.NewEvent("evt")
			evt.SetStartAt(midnightToday.Add(-5 * time.Hour))
			evt.SetEndAt(midnightToday.Add(-4 * time.Hour))

			actual, err := calendar.EventEndsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeFalse())
		})

		It("returns false if event ends after midnight tomorrow", func() {
			evt := ical.NewEvent("evt")
			evt.SetStartAt(midnightToday.Add(1 * time.Hour))
			evt.SetEndAt(midnightToday.Add(36 * time.Hour))

			actual, err := calendar.EventEndsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeFalse())
		})

		It("returns false if event ends at midnight tomorrow", func() {
			evt := ical.NewEvent("evt")
			evt.SetStartAt(midnightToday.Add(1 * time.Hour))
			evt.SetEndAt(midnightToday.Add(24 * time.Hour))

			actual, err := calendar.EventEndsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeFalse())
		})

		It("returns true if event ends between midnight today and midnight tomorrow", func() {
			evt := ical.NewEvent("evt")
			evt.SetStartAt(midnightToday.Add(4 * time.Hour))
			evt.SetEndAt(midnightToday.Add(6 * time.Hour))

			actual, err := calendar.EventEndsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeTrue())
		})

		It("returns true if event ends at midnight today", func() {
			evt := ical.NewEvent("evt")
			evt.SetStartAt(midnightToday.Add(-2 * time.Hour))
			evt.SetEndAt(midnightToday)

			actual, err := calendar.EventEndsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeTrue())
		})
	})

	Describe("EventStartAndEnd", func() {

		Context("when an event has a proper start and end time of type DATETIME", func() {
			It("will return those times", func() {
				evt := ical.NewEvent("evt")
				start := time.Date(2020, 01, 30, 12, 20, 00, 00, time.UTC)
				evt.SetProperty(ical.ComponentPropertyDtStart, fmtDATETIME(start))

				end := time.Date(2020, 01, 30, 12, 20, 00, 00, time.UTC)
				evt.SetProperty(ical.ComponentPropertyDtEnd, fmtDATETIME(end))

				actualStart, actualEnd, err := calendar.EventStartAndEnd(evt)
				Expect(err).ToNot(HaveOccurred())
				Expect(actualStart).To(BeTemporally("~", start, 1*time.Minute))
				Expect(actualEnd).To(BeTemporally("~", end, 1*time.Minute))
			})
		})

		Context("when an event has a start time of type DATE but no end time or duration", func() {
			It("will return the start as midnight of the date, and end as the midnight of the next day", func() {
				evt := ical.NewEvent("evt")

				start := time.Date(2020, 01, 30, 12, 20, 00, 00, time.UTC)
				evt.SetProperty(ical.ComponentPropertyDtStart, fmtDATE(start))

				expectedStart := time.Date(2020, 01, 30, 00, 00, 00, 00, time.UTC)
				expectedEnd := time.Date(2020, 01, 31, 00, 00, 00, 0, time.UTC)

				actualStart, actualEnd, err := calendar.EventStartAndEnd(evt)
				Expect(err).ToNot(HaveOccurred())
				Expect(actualStart).To(BeTemporally("~", expectedStart, 1*time.Minute))
				Expect(actualEnd).To(BeTemporally("~", expectedEnd, 1*time.Minute))
			})
		})

		Context("when an event has a start time of type DATETIME and a duration", func() {
			It("will return the defined start time, and the end time as start + duration", func() {
				evt := ical.NewEvent("evt")

				start := time.Date(2020, 01, 30, 12, 20, 00, 00, time.UTC)
				evt.SetProperty(ical.ComponentPropertyDtStart, fmtDATETIME(start))

				evt.SetProperty("DURATION", "30M")

				expectedStart := time.Date(2020, 01, 30, 12, 20, 00, 00, time.UTC)
				expectedEnd := time.Date(2020, 01, 30, 12, 50, 00, 0, time.UTC)

				actualStart, actualEnd, err := calendar.EventStartAndEnd(evt)
				Expect(err).ToNot(HaveOccurred())
				Expect(actualStart).To(BeTemporally("~", expectedStart, 1*time.Minute))
				Expect(actualEnd).To(BeTemporally("~", expectedEnd, 1*time.Minute))
			})
		})

		Context("when an event has a start time of type DATETIME but no end time or duration", func() {
			It("will return the start and end time as the same time", func() {
				evt := ical.NewEvent("evt")

				start := time.Date(2020, 01, 30, 12, 20, 00, 00, time.UTC)
				evt.SetProperty(ical.ComponentPropertyDtStart, fmtDATETIME(start))

				expectedStart := time.Date(2020, 01, 30, 12, 20, 00, 00, time.UTC)

				actualStart, actualEnd, err := calendar.EventStartAndEnd(evt)
				Expect(err).ToNot(HaveOccurred())
				Expect(actualStart).To(BeTemporally("~", expectedStart, 1*time.Minute))
				Expect(actualEnd).To(BeTemporally("~", expectedStart, 1*time.Minute))
			})
		})
	})

	Describe("IsAllDayEvent", func() {
		It("returns true when start and end values are of DATE type", func() {
			evt := ical.NewEvent("evt")
			evt.SetProperty(ical.ComponentPropertyDtStart, "20200130")
			evt.SetProperty(ical.ComponentPropertyDtEnd, "20200130")

			actual, err := calendar.IsAllDayEvent(evt)

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeTrue())
		})

		It("returns false when start is of type DATE and end is of type DATETIME", func() {
			evt := ical.NewEvent("evt")
			evt.SetProperty(ical.ComponentPropertyDtStart, "20200130")
			evt.SetProperty(ical.ComponentPropertyDtEnd, "20200130T045200Z")

			actual, err := calendar.IsAllDayEvent(evt)

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeFalse())
		})

		It("returns false when start is of type DATETIME and end is of type DATE", func() {
			evt := ical.NewEvent("evt")
			evt.SetProperty(ical.ComponentPropertyDtStart, "20200130T013055Z")
			evt.SetProperty(ical.ComponentPropertyDtEnd, "20200130")

			actual, err := calendar.IsAllDayEvent(evt)

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeFalse())
		})

		It("returns true when start is present but has no end time or duration", func() {
			evt := ical.NewEvent("evt")
			evt.SetProperty(ical.ComponentPropertyDtStart, "20200130")

			actual, err := calendar.IsAllDayEvent(evt)

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeTrue())
		})

		It("returns false when the event has start and duration properties", func() {
			evt := ical.NewEvent("evt")
			evt.SetProperty(ical.ComponentPropertyDtStart, "20200130")
			evt.SetProperty("DURATION", "30M")

			actual, err := calendar.IsAllDayEvent(evt)

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeFalse())
		})
	})

	Describe("IsRecurringEventDefinition", func() {
		It("will return false when the event doesn't have an RRULE property", func() {
			evt := ical.NewEvent("evt")

			actual := calendar.IsRecurringEventDefinition(evt)

			Expect(actual).To(BeFalse())
		})

		It("will return true when the event has an RRULE property", func() {
			evt := ical.NewEvent("evt")
			evt.AddRrule("FREQ=WEEKLY;BYDAY=TU")

			actual := calendar.IsRecurringEventDefinition(evt)

			Expect(actual).To(BeTrue())
		})
	})
})
