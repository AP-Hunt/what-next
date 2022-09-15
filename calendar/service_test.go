package calendar_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"path/filepath"

	. "github.com/AP-Hunt/what-next/m/calendar"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	Describe("OpenCalendar", func() {
		It("can open file URLs", func() {
			Expect("../fake.ical").To(BeAnExistingFile(), "the fake calendar file 'fake.ical' does not exist. Have you run `make fake-calendar`?")

			pwd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			calPath, err := filepath.Abs(path.Join(pwd, "..", "fake.ical"))
			Expect(err).ToNot(HaveOccurred())

			calendarSvc := NewCalendarService()
			cal, err := calendarSvc.OpenCalendar("file://" + calPath)

			Expect(err).ToNot(HaveOccurred())
			Expect(len(cal.Events())).To(BeNumerically(">=", 1))
		})

		It("can open HTTP urls", func() {
			Expect("../fake.ical").To(BeAnExistingFile(), "the fake calendar file 'fake.ical' does not exist. Have you run `make fake-calendar`?")

			localHttpSrv := httptest.NewServer(http.FileServer(http.Dir("..")))
			defer localHttpSrv.Close()

			url, err := url.JoinPath(localHttpSrv.URL, "fake.ical")
			Expect(err).ToNot(HaveOccurred())

			calendarSvc := NewCalendarService()
			cal, err := calendarSvc.OpenCalendar(url)

			Expect(err).ToNot(HaveOccurred())
			Expect(len(cal.Events())).To(BeNumerically(">=", 1))
		})
	})
})
