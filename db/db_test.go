package db_test

import (
	"github.com/AP-Hunt/what-next/m/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Db", func() {
	Describe("Migrate", func() {
		It("applies migrations to a supplied database", func() {
			dbConn, err := db.Connect(":memory:")
			Expect(err).ToNot(HaveOccurred())
			defer dbConn.Close()

			err = db.Migrate(dbConn.DB)
			Expect(err).ToNot(HaveOccurred())

			var count int
			row := dbConn.QueryRow("SELECT count(*) FROM goose_db_version")

			err = row.Scan(&count)
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(BeNumerically(">", 0))
		})

		It("does not apply migrations twice", func() {
			dbConn, err := db.Connect(":memory:")
			Expect(err).ToNot(HaveOccurred())
			defer dbConn.Close()

			// Apply once
			err = db.Migrate(dbConn.DB)
			Expect(err).ToNot(HaveOccurred())

			// Apply twice
			err = db.Migrate(dbConn.DB)
			Expect(err).ToNot(HaveOccurred())

			var count int
			row := dbConn.QueryRow("SELECT count(*) FROM goose_db_version")

			err = row.Scan(&count)
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(BeNumerically(">", 0))
		})
	})
})
