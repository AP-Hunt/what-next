package todo_test

import (
	"context"
	"time"

	"github.com/AP-Hunt/what-next/m/db"
	"github.com/AP-Hunt/what-next/m/todo"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pressly/goose/v3"
)

var _ = Describe("Repository", func() {
	goose.SetLogger(goose.NopLogger())
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
		now := time.Now()
		duration := 60 * time.Second
		item := todo.TodoItem{
			Id:        999,
			Action:    "doing some stuff",
			DueDate:   &now,
			Duration:  &duration,
			Completed: false,
		}

		addedItem, err := repo.Add(item)
		Expect(err).ToNot(HaveOccurred())

		Expect(addedItem).ToNot(BeIdenticalTo(item))
		Expect(addedItem.Id).To(Equal(1))
		Expect(addedItem.Action).To(Equal(item.Action))
		Expect(*addedItem.DueDate).To(BeTemporally("==", *item.DueDate))
		Expect(*addedItem.Duration).To(Equal(*item.Duration))
	})

	It("can fetch an item that was previously inserted", func() {
		now := time.Now()
		duration := 60 * time.Second
		item := todo.TodoItem{
			Action:    "new item",
			DueDate:   &now,
			Duration:  &duration,
			Completed: false,
		}

		addedItem, err := repo.Add(item)
		Expect(err).ToNot(HaveOccurred())

		retrievedItem, err := repo.Get(addedItem.Id)
		Expect(err).ToNot(HaveOccurred())

		Expect(retrievedItem).To(Equal(addedItem))
		Expect(retrievedItem).ToNot(BeIdenticalTo(addedItem))
	})

	It("can list all existing items", func() {
		_, err := repo.Add(todo.TodoItem{Action: "Item 1"})
		Expect(err).ToNot(HaveOccurred())

		_, err = repo.Add(todo.TodoItem{Action: "Item 2"})
		Expect(err).ToNot(HaveOccurred())

		collection, err := repo.List()
		Expect(err).ToNot(HaveOccurred())
		Expect(collection.Len()).To(Equal(2))
	})

	It("update an item", func() {
		item := todo.TodoItem{
			Id:        0,
			Action:    "initial",
			DueDate:   nil,
			Duration:  nil,
			Completed: false,
		}

		addedItem, err := repo.Add(item)
		Expect(err).ToNot(HaveOccurred())

		dueDate := time.Date(2020, 01, 01, 00, 00, 00, 00, time.Local)
		addedItem.DueDate = &dueDate

		duration := 30 * time.Minute
		addedItem.Duration = &duration

		addedItem.Action = "updated"
		addedItem.Complete()

		updatedItem, err := repo.Update(addedItem)

		Expect(err).ToNot(HaveOccurred())
		Expect(updatedItem).ToNot(BeNil())
		Expect(updatedItem).ToNot(BeIdenticalTo(addedItem), "the repo update method should not be returning the exact same struct")

		Expect(*updatedItem.DueDate).To(BeTemporally("==", dueDate))
		Expect(int(*updatedItem.Duration)).To(BeNumerically("==", int(duration)))
		Expect(updatedItem.Completed).To(BeTrue())
        Expect(updatedItem.CompletedAt).ToNot(BeNil())
        Expect(*updatedItem.CompletedAt).To(BeTemporally("~", time.Now(), 5 * time.Second))
		Expect(updatedItem.Action).To(Equal("updated"))
	})
})
