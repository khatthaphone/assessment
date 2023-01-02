package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/khatthaphone/assesment/expense"
	"github.com/labstack/echo/v4"

	_ "github.com/lib/pq"
)

func main() {
	closeDB := expense.InitDB()
	defer closeDB()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/expenses", expense.AddExpenseHandler)

	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
}
