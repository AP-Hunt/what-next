package views

import (
	"fmt"
	"math"
	"strings"
)

var errColSpansInvalid = fmt.Errorf("Total col spans must be between 1 and 12")

func colSpansAreValid(colSpans []int) bool {
	totalSpans := 0
	for _, span := range colSpans {
		totalSpans = totalSpans + int(math.Abs(float64(span)))
	}

	return totalSpans <= 12 && totalSpans > 0
}

func oneColRowFormatter(colSpans [1]int) (func(colVals [1]string) string, error) {
	if !colSpansAreValid(colSpans[:]) {
		return nil, errColSpansInvalid
	}

	return func(colVals [1]string) string {
		return formatTableRow(colSpans[:], colVals[:])
	}, nil
}

func twoColRowFormatter(colSpans [2]int) (func(colVals [2]string) string, error) {
	if !colSpansAreValid(colSpans[:]) {
		return nil, errColSpansInvalid
	}

	return func(colVals [2]string) string {
		return formatTableRow(colSpans[:], colVals[:])
	}, nil
}

func threeColRowFormatter(colSpans [3]int) (func(colVals [3]string) string, error) {
	if !colSpansAreValid(colSpans[:]) {
		return nil, errColSpansInvalid
	}

	return func(colVals [3]string) string {
		return formatTableRow(colSpans[:], colVals[:])
	}, nil
}

func fourColRowFormatter(colSpans [4]int) (func(colVals [4]string) string, error) {
	if !colSpansAreValid(colSpans[:]) {
		return nil, errColSpansInvalid
	}

	return func(colVals [4]string) string {
		return formatTableRow(colSpans[:], colVals[:])
	}, nil
}

func fiveColRowFormatter(colSpans [5]int) (func(colVals [5]string) string, error) {
	if !colSpansAreValid(colSpans[:]) {
		return nil, errColSpansInvalid
	}

	return func(colVals [5]string) string {
		return formatTableRow(colSpans[:], colVals[:])
	}, nil
}

func sixColRowFormatter(colSpans [6]int) (func(colVals [6]string) string, error) {
	if !colSpansAreValid(colSpans[:]) {
		return nil, errColSpansInvalid
	}

	return func(colVals [6]string) string {
		return formatTableRow(colSpans[:], colVals[:])
	}, nil
}

func sevenColRowFormatter(colSpans [7]int) (func(colVals [7]string) string, error) {
	if !colSpansAreValid(colSpans[:]) {
		return nil, errColSpansInvalid
	}

	return func(colVals [7]string) string {
		return formatTableRow(colSpans[:], colVals[:])
	}, nil
}

func eightColRowFormatter(colSpans [8]int) (func(colVals [8]string) string, error) {
	if !colSpansAreValid(colSpans[:]) {
		return nil, errColSpansInvalid
	}

	return func(colVals [8]string) string {
		return formatTableRow(colSpans[:], colVals[:])
	}, nil
}

func nineColRowFormatter(colSpans [9]int) (func(colVals [9]string) string, error) {
	if !colSpansAreValid(colSpans[:]) {
		return nil, errColSpansInvalid
	}

	return func(colVals [9]string) string {
		return formatTableRow(colSpans[:], colVals[:])
	}, nil
}

func tenColRowFormatter(colSpans [10]int) (func(colVals [10]string) string, error) {
	if !colSpansAreValid(colSpans[:]) {
		return nil, errColSpansInvalid
	}

	return func(colVals [10]string) string {
		return formatTableRow(colSpans[:], colVals[:])
	}, nil
}

func elevenColRowFormatter(colSpans [11]int) (func(colVals [11]string) string, error) {
	if !colSpansAreValid(colSpans[:]) {
		return nil, errColSpansInvalid
	}

	return func(colVals [11]string) string {
		return formatTableRow(colSpans[:], colVals[:])
	}, nil
}

func twelveColRowFormatter(colSpans [12]int) (func(colVals [12]string) string, error) {
	if !colSpansAreValid(colSpans[:]) {
		return nil, errColSpansInvalid
	}

	return func(colVals [12]string) string {
		return formatTableRow(colSpans[:], colVals[:])
	}, nil
}

// formatTableRow returns a string representing the row of a table
//
// colSpans and colValues must be the same length, and the sum of
// colSpans must be no greater than 12 and at least 1. However,
// formatTableRow is an internal function and these constraints are
// not enforced.
//
// Use one of the *ColRowFormatter functions instead
//
// Each column will be space-padded to the length of colSpans[i]
// and contain the content from colValues[i]
//
// To left align columns, provide negative colspans
func formatTableRow(colSpans []int, colValues []string) string {
	formatString := strings.Repeat("%*s", len(colSpans)) + "\n"

	zip := []interface{}{}
	for i := 0; i < len(colSpans); i++ {
		span := colSpans[i]
		val := colValues[i]

		zip = append(zip, layoutColCharWidth(span))
		zip = append(zip, val)
	}

	return fmt.Sprintf(formatString, zip...)
}
