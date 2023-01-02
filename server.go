package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	_ "github.com/lib/pq"

	"github.com/khatthaphone/assesment/expense"
)

func main() {
	db := expense.InitDB()
	defer db.Close()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	handler := expense.NewHandler(db)
	e.POST("/expenses", handler.AddExpenseHandler)
	e.GET("/expenses/:id", handler.GetExpenseHandler)
	e.PUT("/expenses/:id", handler.UpdateExpenseHandler)
	e.GET("/expenses", handler.GetAllExpensesHandler)

	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
}
