package expense

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Err struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type Expense struct {
	ID     int64    `json:"id"`
	Title  string   `json:"title"`
	Amount int      `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *handler {
	return &handler{
		db: db,
	}
}

func (h *handler) AddExpenseHandler(c echo.Context) error {

	if c.Request().Body == http.NoBody {
		return c.JSON(http.StatusBadRequest, Err{Message: "No request body provided", Error: "NoBody"})
	}

	var e Expense
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "Error occurred when binding json from request body", Error: err.Error()})
	}
	if e.Title == "" || e.Note == "" || len(e.Tags) == 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Invalid JSON structure", Error: "InvalidJson"})
	}

	sql := `INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3, $4) RETURNING id`
	row := h.db.QueryRow(sql, e.Title, e.Amount, e.Note, pq.Array(&e.Tags))

	err = row.Scan(&e.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Error occured when binding sql result", Error: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)

}

func (h *handler) GetExpenseHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "No epense id provided", Error: err.Error()})
	}

	// TODO: use db.Prepare
	sql := `SELECT id, title, amount, note, tags FROM expenses WHERE id = $1`

	row := h.db.QueryRow(sql, id)

	var e Expense
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, (pq.Array)(&e.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Error occured when binding sql result", Error: err.Error()})
	}

	return c.JSON(http.StatusOK, e)
}

func (h *handler) UpdateExpenseHandler(c echo.Context) error {

	if c.Request().Body == http.NoBody {
		return c.JSON(http.StatusBadRequest, Err{Message: "No request body provided", Error: "NoBody"})
	}

	id := c.Param("id")

	var e Expense
	err := c.Bind(&e)
	if err != nil || e.Title == "" || e.Note == "" || len(e.Tags) == 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	sql := `UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id, title, amount, note, tags`
	row := h.db.QueryRow(sql, e.Title, e.Amount, e.Note, pq.Array(&e.Tags), id)

	var result Expense
	err = row.Scan(&result.ID, &result.Title, &result.Amount, &result.Note, (pq.Array)(&result.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Error occured when binding sql result", Error: err.Error()})
	}

	return c.JSON(http.StatusOK, result)

}

func (h *handler) GetAllExpensesHandler(c echo.Context) error {
	// TODO: use db.Prepare
	sql := `SELECT id, title, amount, note, tags FROM expenses`

	rows, err := h.db.Query(sql)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Error occured when querying expenses", Error: err.Error()})
	}

	var expenses []Expense

	for rows.Next() {
		var e Expense
		err := rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, (pq.Array)(&e.Tags))
		if err != nil {

			return c.JSON(http.StatusInternalServerError, Err{Message: "Error occured when binding sql result", Error: err.Error()})
		}

		expenses = append(expenses, e)
	}

	return c.JSON(http.StatusOK, expenses)
}
