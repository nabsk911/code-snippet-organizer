package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/nabsk911/code-snippet-organizer/internal/handlers"
	"github.com/nabsk911/code-snippet-organizer/internal/store"
	"github.com/nabsk911/code-snippet-organizer/migrations"
)

type Application struct {
	DB             *sql.DB
	Logger         *log.Logger
	UserHandler    *handlers.UserHandler
	SnippetHandler *handlers.SnippetHandler
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()

	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.MigrationsFS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Our stores go here
	userStore := store.NewPostgresUserStore(pgDB)
	snippetStore := store.NewPostgresSnippetStore(pgDB)

	// Our handlers go here
	userHandler := handlers.NewUserHandler(userStore, logger)
	snippetHandler := handlers.NewSnippetHandler(logger, snippetStore)

	return &Application{
		DB:             pgDB,
		Logger:         logger,
		UserHandler:    userHandler,
		SnippetHandler: snippetHandler,
	}, nil
}
