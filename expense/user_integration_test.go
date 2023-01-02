package expense

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	addExpenseJson  = `{"title": "strawberry smoothie","amount": 79,"note": "night market promotion discount 10 bath","tags": ["food", "beverage"]}`
	editExpenseJson = `{"title": "apple smoothie","amount": 89,"note": "no discount","tags": ["beverage"]}`
)

func setup() *sql.DB {
	db := InitDB()
	return db
}

func TestAddExpense(t *testing.T) {

	// TODO: db init, migrate, seed
	setup()
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(addExpenseJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db := NewDB(InitDB())
	h := &handler{db}

	if assert.NoError(t, h.AddExpenseHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
	// TODO: Test res body match req
}

func TestGetExpenseById(t *testing.T) {
	id := 29
	setup()
	e := echo.New()
	url := fmt.Sprintf("/expenses/%d", id)
	fmt.Printf("Testing get ID: %s\n", url)
	req := httptest.NewRequest(http.MethodGet, url, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Manually set path
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("29")
	db := NewDB(InitDB())
	h := &handler{db}

	if assert.NoError(t, h.GetExpenseHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	// TODO: Test res body match req
}

func TestUpdateExepense(t *testing.T) {
	id := 29
	setup()
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/expenses/%d", id), strings.NewReader(editExpenseJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Manually set path
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues("29")
	db := NewDB(InitDB())
	h := &handler{db}

	if assert.NoError(t, h.UpdateExpenseHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	// TODO: Test res body match req
}
