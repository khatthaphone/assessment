package expense

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setup() (*echo.Echo, *sql.DB, func()) {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to DB", err)
	}

	dropTableSql := `DROP TABLE IF EXISTS expenses`
	_, err = conn.Exec(dropTableSql)
	if err != nil {
		log.Fatal("can't drop existing table", err)
	}

	createTableSql := `
		CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			title TEXT,
			amount FLOAT,
			note TEXT,
			tags TEXT[]
		);
	`
	_, err = conn.Exec(createTableSql)
	if err != nil {
		log.Fatal("can't create table", err)
	}

	close := func() {
		conn.Exec(`TRUNCATE expenses`)

		conn.Close()
	}

	e := echo.New()

	return e, conn, close
}

func TestAddExpense(t *testing.T) {

	// TODO: db init, migrate, seedz
	e, db, close := setup()
	defer close()

	expense := &Expense{
		Title:  "apple smotthie",
		Amount: 89,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}

	expenseJson, err := json.Marshal(expense)
	if err != nil {
		t.Fatal("Failed to contruct json req body", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(string(expenseJson)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHandler(db)

	if assert.NoError(t, h.AddExpenseHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var res Expense
		json.Unmarshal(rec.Body.Bytes(), &res)

		assert.IsType(t, int64(1), res.ID)
		assert.Equal(t, expense.Title, res.Title)
		assert.Equal(t, expense.Amount, res.Amount)
		assert.Equal(t, expense.Note, res.Note)
		assert.Equal(t, expense.Tags, res.Tags)
	}
}

func TestAddExpenseNoJsonShouldFail(t *testing.T) {

	// TODO: db init, migrate, seedz
	e, db, close := setup()
	defer close()

	req := httptest.NewRequest(http.MethodPost, "/expenses", http.NoBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHandler(db)

	if assert.NoError(t, h.AddExpenseHandler(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestAddExpenseInvalidJsonShouldFail(t *testing.T) {

	// TODO: db init, migrate, seedz
	e, db, close := setup()
	defer close()

	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(`{"hello": "world!"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHandler(db)

	if assert.NoError(t, h.AddExpenseHandler(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestGetExpenseById(t *testing.T) {
	e, db, close := setup()
	defer close()

	expense, _ := migrateExpenseForTest(db)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/expenses/%v", 1), strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Manually set path
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%v", 1))
	h := NewHandler(db)

	if assert.NoError(t, h.GetExpenseHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var res Expense
		json.Unmarshal(rec.Body.Bytes(), &res)

		assert.IsType(t, int64(1), res.ID)
		assert.Equal(t, expense.Title, res.Title)
		assert.Equal(t, expense.Amount, res.Amount)
		assert.Equal(t, expense.Note, res.Note)
		assert.Equal(t, expense.Tags, res.Tags)
	}
}

func TestUpdateExepense(t *testing.T) {
	e, db, close := setup()
	defer close()

	expense, expenseJson := migrateExpenseForTest(db)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/expenses/%d", 1), strings.NewReader(string(expenseJson)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Manually set path
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", 1))

	h := NewHandler(db)

	if assert.NoError(t, h.UpdateExpenseHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var res Expense
		json.Unmarshal(rec.Body.Bytes(), &res)

		// assert.Equal(t, expense.ID, res.ID)
		assert.Equal(t, expense.Title, res.Title)
		assert.Equal(t, expense.Amount, res.Amount)
		assert.Equal(t, expense.Note, res.Note)
		assert.Equal(t, expense.Tags, res.Tags)

	}
}

func TestGetAllExpenses(t *testing.T) {
	e, db, close := setup()
	defer close()

	// Add 2 expenses
	migrateExpenseForTest(db)
	migrateExpenseForTest(db)

	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHandler(db)

	if assert.NoError(t, h.GetAllExpensesHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var res []Expense
		json.Unmarshal(rec.Body.Bytes(), &res)

		expenses := []Expense{
			{
				Title:  "apple smotthie",
				Amount: 89,
				Note:   "no discount",
				Tags:   []string{"beverage"},
			},
		}

		assert.IsType(t, expenses, res)
		assert.Equal(t, 2, len(res))
	}
}

func migrateExpenseForTest(db *sql.DB) (*Expense, []byte) {

	expense := &Expense{
		Title:  "apple smotthie",
		Amount: 89,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}

	expenseJson, err := json.Marshal(expense)
	if err != nil {
		log.Fatal("Failed contructing update req body")
		return nil, nil
	}

	res := db.QueryRow(`INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3, $4) RETURNING id`, expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags))
	var insertId int
	err = res.Scan(&insertId)
	if err != nil {
		log.Fatalf("Error migrating for update exepense: %v", err.Error())
		return nil, nil
	}

	return expense, expenseJson
}
