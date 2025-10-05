package routes

import (
	"net/http"

	"github.com/nabsk911/code-snippet-organizer/internal/app"
	"github.com/nabsk911/code-snippet-organizer/internal/middleware"
)

func SetupRoutes(app *app.Application) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /register", app.UserHandler.HandleRegister)
	router.HandleFunc("POST /login", app.UserHandler.HandleLogin)

	router.HandleFunc("POST /snippets", middleware.Authentication(app.SnippetHandler.HandleCreateSnippet))
	router.HandleFunc("GET /snippets", middleware.Authentication(app.SnippetHandler.HandleGetSnippetsByUserID))
	router.HandleFunc("DELETE /snippets/{id}", middleware.Authentication(app.SnippetHandler.HandleDeleteSnippet))
	router.HandleFunc("PUT /snippets/{id}", middleware.Authentication(app.SnippetHandler.HandleUpdateSnippet))

	return router
}
