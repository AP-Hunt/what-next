package calendar_test

import (
	"time"

	"github.com/AP-Hunt/what-next/m/calendar"
	ical "github.com/arran4/golang-ical"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
			evt.SetEndAt(midnightToday.Add(-4 * time.Hour))

			actual, err := calendar.EventEndsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeFalse())
		})

		It("returns false if event ends after midnight tomorrow", func() {
			evt := ical.NewEvent("evt")
			evt.SetEndAt(midnightToday.Add(36 * time.Hour))

			actual, err := calendar.EventEndsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeFalse())
		})

		It("returns false if event ends at midnight tomorrow", func() {
			evt := ical.NewEvent("evt")
			evt.SetEndAt(midnightToday.Add(24 * time.Hour))

			actual, err := calendar.EventEndsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeFalse())
		})

		It("returns true if event ends between midnight today and midnight tomorrow", func() {
			evt := ical.NewEvent("evt")
			evt.SetEndAt(midnightToday.Add(6 * time.Hour))

			actual, err := calendar.EventEndsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeTrue())
		})

		It("returns true if event ends at midnight today", func() {
			evt := ical.NewEvent("evt")
			evt.SetEndAt(midnightToday)

			actual, err := calendar.EventEndsToday(evt)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(BeTrue())
		})
	})

	Describe("EventStartAndEnd", func() {
		Context("when an event has a proper start and end time", func() {
			It("will return those", func() {
				evt := ical.NewEvent("evt")
				start := time.Now().Add(-10 * time.Minute).UTC()
				evt.SetStartAt(start)
				end := time.Now().Add(10 * time.Minute).UTC()
				evt.SetEndAt(end)

				actualStart, actualEnd, err := calendar.EventStartAndEnd(evt)
				Expect(err).ToNot(HaveOccurred())
				Expect(actualStart).To(BeTemporally("~", start, 1*time.Minute))
				Expect(actualEnd).To(BeTemporally("~", end, 1*time.Minute))
			})
		})

		Context("when an event has a start time but no end time", func() {
			It("will return the start and end times as the same value, as per the spec", func() {
				evt := ical.NewEvent("evt")
				start := time.Now().Add(-10 * time.Minute).UTC()
				evt.SetStartAt(start)

				actualStart, actualEnd, err := calendar.EventStartAndEnd(evt)
				Expect(err).ToNot(HaveOccurred())
				Expect(actualStart).To(BeTemporally("~", start, 1*time.Minute))
				Expect(actualEnd).To(BeTemporally("~", start, 1*time.Minute))
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
})
