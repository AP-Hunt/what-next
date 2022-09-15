package calendar

type CalendarRecord struct {
	Id          int
	DisplayName string `db:"display_name"`
	URL         string `db:"calendar_url"`
}
