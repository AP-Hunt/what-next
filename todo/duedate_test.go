package todo_test

import (
	"time"

	. "github.com/AP-Hunt/what-next/m/todo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/golang-module/carbon"
)

var _ = Describe("Duedate", func() {
	DescribeTable("ParseDueDate",
		func(input string, expected time.Time, shouldError bool) {
			actual, err := ParseDueDate(input)

			if shouldError {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).ToNot(HaveOccurred())
			}

			Expect(actual).To(Equal(expected))
		},
		Entry("using @today", "@today", carbon.Now().EndOfDay().Carbon2Time(), false),
		Entry("using @tod", "@tod", carbon.Now().EndOfDay().Carbon2Time(), false),
		Entry("using @tomorrow", "@tomorrow", carbon.Tomorrow().EndOfDay().Carbon2Time(), false),
		Entry("using @tom", "@tom", carbon.Tomorrow().EndOfDay().Carbon2Time(), false),
		Entry("using @tmrw", "@tmrw", carbon.Tomorrow().EndOfDay().Carbon2Time(), false),

		// Deliberately not testing a variety of date strings.
		// If any of the above shortcuts aren't used, it falls back to using
		// github.com/araddon/dateparse for date parsing
		Entry("using a date string", "2022-01-01T14:12:11", time.Date(2022, 1, 1, 14, 12, 11, 0, time.Local), false),
	)
})
