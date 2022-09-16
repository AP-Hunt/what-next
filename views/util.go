package views

import (
	"fmt"
	"math"

	"golang.org/x/crypto/ssh/terminal"
)

var termWidth int

func layoutColCharWidth(n int) int {
	oneColWidth := math.Floor(float64(termWidth / 12))
	return int(oneColWidth) * n
}

func init() {
	width, _, err := terminal.GetSize(0)
	if err != nil {
		panic(fmt.Sprintf("getting terminal size: %s", err))
	}
	termWidth = width
}
