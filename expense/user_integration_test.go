package expense

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	addExpenseJson = `{"title": "strawberry smoothie","amount": 79,"note": "night market promotion discount 10 bath","tags": ["food", "beverage"]}`
)

func setup() *sql.DB {
	db := InitDB()
	return db
}

func TestAddExpense(t *testing.T) {

	t.Run("Call API: Add expense", func(t *testing.T) {
		setup()
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(addExpenseJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, AddExpenseHandler(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
		// TODO: Test res body match req

	})
}