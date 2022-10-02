package integration_test_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("todo add", func() {
	Describe("--duration", func() {
		DescribeTable("accepts a variety of duration specifiers",
			func(specifier string) {
				RunIntegrationTest(func(exec Executor, cfg *testConfig) {
					err := exec([]string{"todo", "add", "test action", "--duration", specifier})
					Expect(err).ToNot(HaveOccurred())
				})
			},
			Entry("seconds", "30s"),
			Entry("minutes", "30m"),
			Entry("hours", "30h"),
			Entry("minutes and seconds", "1m30s"),
			Entry("hours and minutes", "1h30m"),
		)

		DescribeTable("does not support duration specifiers greater than hours",
			func(specifier string) {
				RunIntegrationTest(func(exec Executor, cfg *testConfig) {
					err := exec([]string{"todo", "add", "test action", "--duration", specifier})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unknown unit"))
				})
			},
			Entry("days", "1d"),
			Entry("weeks", "1w"),
			Entry("months", "1M"),
			Entry("years", "1y"),
		)
	})

	Describe("--due", func() {
		DescribeTable("accepts dates in a variety of formats",
			func(date string) {
				RunIntegrationTest(func(exec Executor, cfg *testConfig) {
					err := exec([]string{"todo", "add", "test action", "--due", date})
					Expect(err).ToNot(HaveOccurred())
				})
			},
			Entry("ISO 8601 date", "2020-01-01"),
			Entry("ISO 8601 datetime", "2020-01-01T15:30:22+01:00"),
			Entry("RFC1123", time.Now().Format(time.RFC1123)),
			Entry("RFC1123Z", time.Now().Format(time.RFC1123Z)),
			Entry("RFC3339", time.Now().Format(time.RFC3339)),
			Entry("RFC3339Nano", time.Now().Format(time.RFC3339Nano)),
			Entry("RFC822", time.Now().Format(time.RFC822)),
			Entry("RFC822Z", time.Now().Format(time.RFC822Z)),
			Entry("RFC850", time.Now().Format(time.RFC850)),
			Entry("ANSIC", time.Now().Format(time.ANSIC)),
			Entry("@today", "@today"),
			Entry("@tod", "@tod"),
			Entry("@tom", "@tom"),
			Entry("@tmrw", "@tmrw"),
			Entry("@tomorrow", "@tomorrow"),
		)
	})
})
