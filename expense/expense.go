package expense

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Err struct {
	Message string
}

type Expense struct {
	ID     int64    `json:"id"`
	Title  string   `json:"title"`
	Amount int      `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type handler struct {
	db *DB
}

func NewHandler(db *DB) *handler {
	return &handler{
		db: db,
	}
}

func (h *handler) AddExpenseHandler(c echo.Context) error {

	var e Expense
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	sql := `INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3, $4) RETURNING id`
	row := db.QueryRow(sql, e.Title, e.Amount, e.Note, pq.Array(&e.Tags))

	err = row.Scan(&e.ID)
	if err != nil {
		log.Fatal("Error adding new expense", err.Error())
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)

}

func (h *handler) GetExpenseHandler(c echo.Context) error {
	fmt.Printf("Raw param: %v\n", c.Param("id"))
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	fmt.Printf("Searching for expense ID: %v\n", id)

	// TODO: use db.Prepare
	sql := `SELECT id, title, amount, note, tags FROM expenses WHERE id = $1`

	row := db.QueryRow(sql, id)

	var e Expense
	err := row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, (pq.Array)(&e.Tags))
	if err != nil {
		log.Fatal("Unable to find expense with sepcified id", err.Error())
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, e)
}

func (h *handler) UpdateExpenseHandler(c echo.Context) error {
	id := c.Param("id")

	var e Expense
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	sql := `UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id, title, amount, note, tags`
	row := db.QueryRow(sql, e.Title, e.Amount, e.Note, pq.Array(&e.Tags), id)

	var result Expense
	err = row.Scan(&result.ID, &result.Title, &result.Amount, &result.Note, (pq.Array)(&result.Tags))
	if err != nil {
		log.Fatal("Error updating expense", err.Error())
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, result)

}

func (h *handler) GetAllExpensesHandler(c echo.Context) error {
	// TODO: use db.Prepare
	sql := `SELECT id, title, amount, note, tags FROM expenses`

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal("Unable to find expenses", err.Error())
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	var expenses []Expense

	for rows.Next() {
		var e Expense
		err := rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, (pq.Array)(&e.Tags))
		if err != nil {
			log.Fatal("Unable to find expenses", err.Error())
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}

		expenses = append(expenses, e)
	}

	return c.JSON(http.StatusOK, expenses)
}
