package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Expense struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Amount int      `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

func AddExpenseHandler(c echo.Context) error {

	var e Expense
	err := c.Bind(&e)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	return c.JSON(http.StatusOK, e)

}
