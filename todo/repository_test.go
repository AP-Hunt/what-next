package todo_test

import (
	"context"
	"time"

	"github.com/AP-Hunt/what-next/m/db"
	"github.com/AP-Hunt/what-next/m/todo"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository", func() {
	var (
		inMemoryConn *sqlx.DB
		repo         *todo.TodoSQLRepository
	)

	BeforeEach(func() {
		conn, err := db.Connect(":memory:")
		Expect(err).ToNot(HaveOccurred())

		err = db.Migrate(conn.DB)
		Expect(err).ToNot(HaveOccurred())

		inMemoryConn = conn
		repo = todo.NewTodoSQLRepository(inMemoryConn, context.Background())
	})

	AfterEach(func() {
		inMemoryConn.Close()
	})

	It("can add a new item", func() {
		By("ignoring the id field")

		item := todo.TodoItem{
			Id:        999,
			Action:    "doing some stuff",
			DueDate:   time.Now(),
			Completed: false,
		}

		addedItem, err := repo.Add(item)
		Expect(err).ToNot(HaveOccurred())

		Expect(addedItem).ToNot(BeIdenticalTo(item))
		Expect(addedItem.Id).To(Equal(1))
	})

	It("can fetch an item that was previously inserted", func() {
		item := todo.TodoItem{
			Action:    "new item",
			DueDate:   time.Now(),
			Completed: false,
		}

		addedItem, err := repo.Add(item)
		Expect(err).ToNot(HaveOccurred())

		retrievedItem, err := repo.Get(addedItem.Id)
		Expect(err).ToNot(HaveOccurred())

		Expect(retrievedItem).To(Equal(addedItem))
		Expect(retrievedItem).ToNot(BeIdenticalTo(addedItem))
	})
})
