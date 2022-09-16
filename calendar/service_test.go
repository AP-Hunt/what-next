package calendar_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"path/filepath"

	. "github.com/AP-Hunt/what-next/m/calendar"
	"github.com/AP-Hunt/what-next/m/db"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pressly/goose/v3"
)

func fakeCalFilePath() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	calPath, err := filepath.Abs(path.Join(pwd, "..", "fake.ical"))
	if err != nil {
		return "", err
	}

	return calPath, nil
}

var _ = Describe("Service", func() {
	goose.SetLogger(goose.NopLogger())
	var (
		inMemoryConn *sqlx.DB
		calendarSvc  *CalendarService
	)

	BeforeEach(func() {
		Expect("../fake.ical").To(BeAnExistingFile(), "the fake calendar file 'fake.ical' does not exist. Have you run `make fake-calendar`?")

		conn, err := db.Connect(":memory:")
		Expect(err).ToNot(HaveOccurred())

		err = db.Migrate(conn.DB)
		Expect(err).ToNot(HaveOccurred())

		inMemoryConn = conn

		calendarSvc = NewCalendarService(inMemoryConn, context.Background())

	})

	AfterEach(func() {
		inMemoryConn.Close()
	})

	Describe("OpenCalendar", func() {
		It("can open file URLs", func() {
			calPath, err := fakeCalFilePath()
			Expect(err).ToNot(HaveOccurred())

			cal, err := calendarSvc.OpenCalendar("file://" + calPath)

			Expect(err).ToNot(HaveOccurred())
			Expect(len(cal.Events())).To(BeNumerically(">=", 1))
		})

		It("can open HTTP urls", func() {
			localHttpSrv := httptest.NewServer(http.FileServer(http.Dir("..")))
			defer localHttpSrv.Close()

			url, err := url.JoinPath(localHttpSrv.URL, "fake.ical")
			Expect(err).ToNot(HaveOccurred())

			cal, err := calendarSvc.OpenCalendar(url)

			Expect(err).ToNot(HaveOccurred())
			Expect(len(cal.Events())).To(BeNumerically(">=", 1))
		})
	})

	Describe("AddCalendar", func() {
		It("will throw an error if the calendar URL can't be reached", func() {
			_, err := calendarSvc.AddCalendar("file://not.a.thing", "display")

			Expect(err).To(HaveOccurred())
		})

		It("will throw an error if the calendar URL doesn't provide a valid calendar", func() {
			notACalFilePath, err := filepath.Abs(path.Join("..", "go.mod"))
			Expect(err).ToNot(HaveOccurred())

			_, err = calendarSvc.AddCalendar("file://"+notACalFilePath, "display")
			Expect(err).To(HaveOccurred())
		})

		It("will insert a new calendar record if the URL provides a valid calendar", func() {
			calPath, err := fakeCalFilePath()
			Expect(err).ToNot(HaveOccurred())

			entry, err := calendarSvc.AddCalendar("file://"+calPath, "display")

			Expect(err).ToNot(HaveOccurred())
			Expect(entry.URL).To(Equal("file://" + calPath))
			Expect(entry.Id).To(Equal(1))
		})

		It("will prevent duplicates by display name", func() {
			calPath, err := fakeCalFilePath()
			Expect(err).ToNot(HaveOccurred())

			_, err = calendarSvc.AddCalendar("file://"+calPath, "display")
			Expect(err).ToNot(HaveOccurred())

			_, err = calendarSvc.AddCalendar("file://"+calPath, "display")

			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(new(ErrDuplicateCalendarDisplayName)))
		})
	})

	Describe("GetCalendarByDisplayName", func() {
		It("will return an error when no such calendar exists", func() {
			_, err := calendarSvc.GetCalendarByDisplayName("foo")
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(new(ErrNotFound)))
		})

		It("will return a valid record when it finds such a calendar", func() {
			calPath, err := fakeCalFilePath()
			Expect(err).ToNot(HaveOccurred())

			addedRecord, err := calendarSvc.AddCalendar("file://"+calPath, "test name")
			Expect(err).ToNot(HaveOccurred())

			fetchedRecord, err := calendarSvc.GetCalendarByDisplayName("test name")
			Expect(err).ToNot(HaveOccurred())
			Expect(fetchedRecord.Id).To(Equal(addedRecord.Id))
			Expect(fetchedRecord.DisplayName).To(Equal("test name"))
		})
	})
})