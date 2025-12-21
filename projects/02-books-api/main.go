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
	fmt.Println(schema)
	_, err = db.Exec(schema)
	if err != nil {
		fmt.Println("Error al ejecutar las migraciones:", err)
		log.Fatal("Error: ", err)
		return
	}

	authorStore := store.NewAuthorStore(db)
	authorService := services.NewAuthorService(authorStore)
	authorHandler := transport.NewAuthorHandler(authorService)

	http.HandleFunc("/author", authorHandler.HandleAuthors)
	http.HandleFunc("/author/", authorHandler.HandleAuthorByID)

	fmt.Println("Servidor escuchando en el puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
