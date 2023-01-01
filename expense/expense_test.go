package expense

import (
	"testing"
)

var (
	addExpenseJson = `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`
)

func TestAddExpense(t *testing.T) {

	t.Run("Call API: Add expense", func(t *testing.T) {
		// e := echo.New()
		// req := httptest.NewRequest(http.MethodPost, "/expense", strings.NewReader(addExpenseJson))
		// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		// rec := httptest.NewRecorder()
		// c := e.NewContext(req, rec)
		// h := &handler{}
	})
}
