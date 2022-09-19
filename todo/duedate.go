package todo

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/golang-module/carbon/v2"
)

func ParseDueDate(input string) (time.Time, error) {
	switch input {
	case "@tod", "@today":
		return carbon.Now().EndOfDay().Carbon2Time(), nil

	case "@tom", "@tmrw", "@tomorrow":
		return carbon.Tomorrow().EndOfDay().Carbon2Time(), nil

	default:
		date, err := dateparse.ParseLocal(input)
		if err != nil {
			return time.Unix(0, 0), fmt.Errorf("invalid due date format: %s", err)
		}

		return date, nil
	}
}
