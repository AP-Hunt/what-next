package calendar

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/AP-Hunt/what-next/m/db"
	ical "github.com/arran4/golang-ical"
	"github.com/jmoiron/sqlx"
)

//counterfeiter:generate -o fakes/ . CalendarServiceInterface
type CalendarServiceInterface interface {
	OpenCalendar(url string) (*ical.Calendar, error)
	AddCalendar(url string, displayName string) (*CalendarRecord, error)
	GetCalendarByDisplayName(displayName string) (*CalendarRecord, error)
	GetAllCalendars() ([]CalendarRecord, error)
	RemoveById(id int) error
}

type CalendarService struct {
	httpClient *http.Client
	db         *sqlx.DB
	ctx        context.Context
	cache      CalendarCacheInterface
}

func NewCalendarService(dbConection *sqlx.DB, cache CalendarCacheInterface, ctx context.Context) *CalendarService {
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	client := &http.Client{Transport: transport}

	return &CalendarService{
		httpClient: client,
		db:         dbConection,
		ctx:        ctx,
		cache:      cache,
	}
}

func (c *CalendarService) OpenCalendar(url string) (*ical.Calendar, error) {
	qualifiedUrl, err := c.resolveUrl(url)
	if err != nil {
		return nil, err
	}
	url = qualifiedUrl.String()

	cachedCal, err := c.cache.Get(url)
	cacheMiss := false
	if err == nil {
		return c.parseCalFromBytes(cachedCal)
	} else {
		if _, ok := err.(*ErrCacheMiss); ok {
			cacheMiss = true
		} else {
			return nil, err
		}
	}

	calBytes, err := c.fetchCalendarOverNetwork(url)
	if err != nil {
		return nil, err
	}

	if cacheMiss {
		err := c.cache.Put(url, calBytes)
		if err != nil {
			return nil, err
		}
	}

	return c.parseCalFromBytes(calBytes)
}

func (c *CalendarService) resolveUrl(maybeUrl string) (*url.URL, error) {
	u, err := url.Parse(maybeUrl)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		absPath, err := filepath.Abs(path.Join(wd, u.Path))
		if err != nil {
			return nil, err
		}

		u.Scheme = "file"
		u.Path = absPath
	}

	return u, nil
}

func (c *CalendarService) fetchCalendarOverNetwork(url string) ([]byte, error) {
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

	return buf.Bytes(), nil
}

func (c *CalendarService) parseCalFromBytes(bs []byte) (*ical.Calendar, error) {
	byteBuf := bytes.NewBuffer(bs)
	cal, err := ical.ParseCalendar(byteBuf)
	if err != nil {
		return nil, fmt.Errorf("parse calendar: %s", err)
	}

	return cal, nil
}

func (c *CalendarService) AddCalendar(url string, displayName string) (*CalendarRecord, error) {
	qualifiedUrl, err := c.resolveUrl(url)
	if err != nil {
		return nil, err
	}
	url = qualifiedUrl.String()

	calBytes, err := c.fetchCalendarOverNetwork(url)
	if err != nil {
		return nil, fmt.Errorf("open calendar '%s': %s", url, err)
	}

	_, err = c.parseCalFromBytes(calBytes)
	if err != nil {
		return nil, err
	}

	_, err = c.GetCalendarByDisplayName(displayName)

	if err == nil {
		return nil, NewErrDuplicateCalendarDisplayName(displayName)
	} else {
		if _, ok := err.(*ErrNotFound); !ok {
			return nil, fmt.Errorf("checking for duplicates: %s", err)
		}
	}

	return db.InTransaction(
		func(tx *sqlx.Tx) (*CalendarRecord, error) {
			row := tx.QueryRowx(
				`
				INSERT INTO calendars 
					(display_name, calendar_url)
				VALUES 
					(?, ?)
				RETURNING *
				`,
				displayName,
				url,
			)

			newRecord := CalendarRecord{}
			err = row.StructScan(&newRecord)
			return &newRecord, err
		},
		c.db,
		c.ctx,
	)
}

func (c *CalendarService) GetCalendarByDisplayName(displayName string) (*CalendarRecord, error) {
	record := CalendarRecord{}
	err := c.db.GetContext(c.ctx, &record, "SELECT * FROM calendars WHERE display_name = ?", displayName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewErrNotFound(err.Error())
		}
		return nil, err
	}

	return &record, nil
}

func (c *CalendarService) GetAllCalendars() ([]CalendarRecord, error) {
	records := []CalendarRecord{}
	err := c.db.SelectContext(c.ctx, &records, "SELECT * FROM calendars")

	if err != nil {
		return []CalendarRecord{}, err
	}

	return records, nil
}

func (c *CalendarService) RemoveById(id int) error {
	_, err := db.InTransaction(
		func(tx *sqlx.Tx) (*int, error) {
			result, err := tx.Exec(
				`
				DELETE FROM calendars
				WHERE id = ?
				`,
				id,
			)

			rowsAffected, err := result.RowsAffected()

			if err != nil {
				errVal := -1
				return &errVal, err
			}

			intRowsAffected := int(rowsAffected)
			return &intRowsAffected, err
		},
		c.db,
		c.ctx,
	)

	return err
}
