package expense

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setup() (*sql.DB, func()) {
	conn, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost/expenses?sslmode=disable")
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

	return conn, close
}

func TestAddExpense(t *testing.T) {

	// TODO: db init, migrate, seedz
	db, close := setup()
	defer close()

	addExpenseJson := `{"title": "strawberry smoothie","amount": 79,"note": "night market promotion discount 10 bath","tags": ["food", "beverage"]}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(addExpenseJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handler{db}

	if assert.NoError(t, h.AddExpenseHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
	// TODO: Test res body match req
}

func TestGetExpenseById(t *testing.T) {
	db, close := setup()
	defer close()

	migrateExpenseForTest(db)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/expenses/%v", 1), strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Manually set path
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%v", 1))
	h := &handler{db}

	if assert.NoError(t, h.GetExpenseHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	// TODO: Test res body match req
}

func TestUpdateExepense(t *testing.T) {
	db, close := setup()
	defer close()

	migrateExpenseForTest(db)

	editExpenseJson := `{"title": "apple smoothie","amount": 89,"note": "no discount","tags": ["beverage"]}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/expenses/%d", 1), strings.NewReader(editExpenseJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Manually set path
	c.SetPath("/expenses/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", 1))

	h := &handler{db}

	if assert.NoError(t, h.UpdateExpenseHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	// TODO: Test res body match req
}

func TestGetAllExpenses(t *testing.T) {
	db, close := setup()
	defer close()

	migrateExpenseForTest(db)
	migrateExpenseForTest(db)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handler{db}

	if assert.NoError(t, h.GetAllExpensesHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	// TODO: Test res body match req
}

func migrateExpenseForTest(db *sql.DB) {
	res := db.QueryRow(`INSERT INTO expenses(title, amount, note, tags) VALUES($1, $2, $3, $4) RETURNING id`, "apple smoothie", 89, "no discount", pq.Array([]string{"beverage"}))
	var insertId int
	err := res.Scan(&insertId)
	if err != nil {
		log.Fatalf("Error migrating for update exepense: %v", err.Error())
		return
	}
}
