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
})
