package main

import (
	"fmt"
	"net/http"

	"example.com/jobboard/internal/database"
	"example.com/jobboard/internal/routes"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := database.InitDatabase()

	router := routes.NewRouter(db)

	port := 4200

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server listening on http://localhost%s\n", addr)

	serveErr := http.ListenAndServe(addr, router)

	if serveErr != nil {
		panic(serveErr)
	}
}