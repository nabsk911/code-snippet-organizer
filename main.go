package main

import (
	"net/http"

	"github.com/nabsk911/code-snippet-organizer/internal/app"
	"github.com/nabsk911/code-snippet-organizer/internal/routes"
)

func main() {
	app, err := app.NewApplication()

	if err != nil {
		panic(err)
	}

	defer app.DB.Close()
	app.Logger.Println("Starting server on port 8080")

	r := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	server.ListenAndServe()
}
