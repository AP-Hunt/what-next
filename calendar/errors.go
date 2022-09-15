package calendar

import "fmt"

type ErrNotFound struct {
	msg string
}

func NewErrNotFound(msg string) *ErrNotFound {
	return &ErrNotFound{msg}
}

func (err ErrNotFound) Error() string {
	return fmt.Sprintf("item not found: %s", err.msg)
}

type ErrDuplicateCalendarDisplayName struct {
	displyName string
}

func NewErrDuplicateCalendarDisplayName(displayName string) *ErrDuplicateCalendarDisplayName {
	return &ErrDuplicateCalendarDisplayName{displayName}
}

func (err ErrDuplicateCalendarDisplayName) Error() string {
	return fmt.Sprintf("duplicate calendar display name: %s", err.displyName)
}
