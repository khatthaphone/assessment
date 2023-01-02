package expense

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func AddExpenseHandler(c echo.Context) error {

	var e Expense
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	sql := `INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3, $4) RETURNING id`
	row := db.QueryRow(sql, e.Title, e.Amount, e.Note, pq.Array(e.Tags))

	err = row.Scan(&e.ID)
	if err != nil {
		log.Fatal("Error adding new expense", err.Error())
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusMultipleChoices, e)

}
