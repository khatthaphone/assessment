package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	_ "github.com/lib/pq"

	"github.com/khatthaphone/assesment/expense"
)

func main() {
	db := expense.InitDB()
	defer db.Close()

	e := echo.New()
	e.GET("/", helloWorldHandler)

	e.Use(loginMiddleware)

	handler := expense.NewHandler(db)
	e.POST("/expenses", handler.AddExpenseHandler)
	e.GET("/expenses/:id", handler.GetExpenseHandler)
	e.PUT("/expenses/:id", handler.UpdateExpenseHandler)
	e.GET("/expenses", handler.GetAllExpensesHandler)

	port := os.Getenv("PORT")

	go func() {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
	}()

	fmt.Println("Please use server.go for main file")
	fmt.Println("Start at port:", os.Getenv("PORT"))

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	fmt.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		fmt.Println("Shutdown err: ", err)
	}
	fmt.Println("Bye bye!")
}

func helloWorldHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
