package expense

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func GetExpenseHandler(c echo.Context) error {
	id := c.Param("id")

	fmt.Sprintf("Searching for expense ID: %s", id)

	// TODO: use db.Prepare
	sql := `SELECT id, title, amount, note, tags FROM expenses WHERE id = $1 LIMIT 1`

	row := db.QueryRow(sql, id)

	var e Expense
	err := row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, (pq.Array)(&e.Tags))
	if err != nil {
		log.Fatal("Unable to find expense with sepcified id", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, e)
}
