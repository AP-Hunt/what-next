package calendar

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	ical "github.com/arran4/golang-ical"
)

//counterfeiter:generate -o fakes/ . CalendarServiceInterface
type CalendarServiceInterface interface {
	OpenCalendar(url string) (*ical.Calendar, error)
}

type CalendarService struct {
	httpClient *http.Client
}

func NewCalendarService() *CalendarService {
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	client := &http.Client{Transport: transport}

	return &CalendarService{
		httpClient: client,
	}
}

func (c *CalendarService) OpenCalendar(url string) (*ical.Calendar, error) {
	resp, err := c.httpClient.Get(url)

	if err != nil {
		return nil, fmt.Errorf("fetch calendar '%s': %s", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("non-success status code: %d", resp.StatusCode)
	}

	buf := bytes.Buffer{}
	bytesCopied, err := io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("copy body stream: %s", err)
	}

	if bytesCopied == 0 {
		return nil, fmt.Errorf("empty response from url '%s'", url)
	}

	cal, err := ical.ParseCalendar(&buf)
	if err != nil {
		return nil, fmt.Errorf("parse calendar: %s", err)
	}

	return cal, nil
}
