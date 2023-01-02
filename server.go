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

	// handler := expense.NewHandler()
	e.POST("/expenses", expense.AddExpenseHandler)
	e.GET("/expenses/:id", expense.GetExpenseHandler)

	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
}
