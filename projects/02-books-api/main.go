package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/chicho69-cesar/backend-go/books/internal/database"
	"github.com/chicho69-cesar/backend-go/books/internal/services"
	"github.com/chicho69-cesar/backend-go/books/internal/store"
	"github.com/chicho69-cesar/backend-go/books/internal/transport"
)

func main() {
	db, err := sql.Open("sqlite3", "./books.db")
	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)
		log.Fatal("Error: ", err)
		return
	}
	defer db.Close()

	schema := database.GetMigrationSchema()
	_, err = db.Exec(schema)
	if err != nil {
		fmt.Println("Error al ejecutar las migraciones:", err)
		log.Fatal("Error: ", err)
		return
	}

	authorStore := store.NewAuthorStore(db)
	authorService := services.NewAuthorService(authorStore)
	authorHandler := transport.NewAuthorHandler(authorService)

	categoryStore := store.NewCategoryStore(db)
	categoryService := services.NewCategoryService(categoryStore)
	categoryHandler := transport.NewCategoryHandler(categoryService)

	configStore := store.NewConfigurationStore(db)
	configService := services.NewConfigurationService(configStore)
	configHandler := transport.NewConfigurationHandler(configService)

	publisherStore := store.NewPublisherStore(db)
	publisherService := services.NewPublisherService(publisherStore)
	publisherHandler := transport.NewPublisherHandler(publisherService)

	http.HandleFunc("/authors", authorHandler.HandleAuthors)
	http.HandleFunc("/authors/", authorHandler.HandleAuthorByID)
	http.HandleFunc("/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/categories/", categoryHandler.HandleCategoryByID)
	http.HandleFunc("/configuration", configHandler.HandleConfiguration)
	http.HandleFunc("/publishers", publisherHandler.HandlePublishers)
	http.HandleFunc("/publishers/", publisherHandler.HandlePublisherByID)

	fmt.Println("Servidor escuchando en el puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
