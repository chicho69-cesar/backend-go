package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/chicho69-cesar/backend-go/books/internal/database"
	"github.com/chicho69-cesar/backend-go/books/internal/logger"
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

	apiLogger, err := logger.NewLogger("./api.log")
	if err != nil {
		fmt.Println("Error al inicializar el logger:", err)
		log.Fatal("Error: ", err)
		return
	}
	defer apiLogger.Close()

	libraryStore := store.NewLibraryStore(db)
	libraryService := services.NewLibraryService(libraryStore)
	libraryHandler := transport.NewLibraryHandler(libraryService)

	authorStore := store.NewAuthorStore(db)
	authorService := services.NewAuthorService(authorStore)
	authorHandler := transport.NewAuthorHandler(authorService)

	bookStore := store.NewBookStore(db)
	copyStore := store.NewCopyStore(db)
	reservationStore := store.NewReservationStore(db)
	bookService := services.NewBookService(bookStore, authorStore, copyStore, reservationStore)
	bookHandler := transport.NewBookHandler(bookService)

	categoryStore := store.NewCategoryStore(db)
	categoryService := services.NewCategoryService(categoryStore)
	categoryHandler := transport.NewCategoryHandler(categoryService)

	configStore := store.NewConfigurationStore(db)
	configService := services.NewConfigurationService(configStore)
	configHandler := transport.NewConfigurationHandler(configService)

	loanStore := store.NewLoanStore(db)
	copyService := services.NewCopyService(copyStore, bookStore, loanStore)
	copyHandler := transport.NewCopyHandler(copyService)

	publisherStore := store.NewPublisherStore(db)
	publisherService := services.NewPublisherService(publisherStore)
	publisherHandler := transport.NewPublisherHandler(publisherService)

	zoneStore := store.NewLibraryZoneStore(db)
	zoneService := services.NewLibraryZoneService(zoneStore)
	zoneHandler := transport.NewLibraryZoneHandler(zoneService)

	shelfStore := store.NewShelfStore(db)
	shelfService := services.NewShelfService(shelfStore, zoneStore)
	shelfHandler := transport.NewShelfHandler(shelfService)

	userStore := store.NewUserStore(db)
	fineStore := store.NewFineStore(db)
	userService := services.NewUserService(userStore, loanStore, reservationStore, fineStore)
	userHandler := transport.NewUserHandler(userService)

	loanService := services.NewLoanService(loanStore, userStore, copyStore, fineStore)
	loanHandler := transport.NewLoanHandler(loanService)

	reservationService := services.NewReservationService(reservationStore, userStore, bookStore, copyStore, fineStore)
	reservationHandler := transport.NewReservationHandler(reservationService)

	fineService := services.NewFineService(fineStore, userStore, loanStore)
	fineHandler := transport.NewFineHandler(fineService)

	http.HandleFunc(
		"/authors",
		apiLogger.Middleware(authorHandler.HandleAuthors),
	)
	http.HandleFunc(
		"/authors/",
		apiLogger.Middleware(authorHandler.HandleAuthorByID),
	)
	http.HandleFunc(
		"/books",
		apiLogger.Middleware(bookHandler.HandleBooks),
	)
	http.HandleFunc(
		"/books/",
		apiLogger.Middleware(bookHandler.HandleBookByID),
	)
	http.HandleFunc(
		"/categories",
		apiLogger.Middleware(categoryHandler.HandleCategories),
	)
	http.HandleFunc(
		"/categories/",
		apiLogger.Middleware(categoryHandler.HandleCategoryByID),
	)
	http.HandleFunc(
		"/configuration",
		apiLogger.Middleware(configHandler.HandleConfiguration),
	)
	http.HandleFunc(
		"/copies",
		apiLogger.Middleware(copyHandler.HandleCopies),
	)
	http.HandleFunc(
		"/copies/",
		apiLogger.Middleware(copyHandler.HandleCopyByID),
	)
	http.HandleFunc(
		"/fines",
		apiLogger.Middleware(fineHandler.HandleFines),
	)
	http.HandleFunc(
		"/fines/",
		apiLogger.Middleware(fineHandler.HandleFineByID),
	)
	http.HandleFunc(
		"/fines/pay/",
		apiLogger.Middleware(fineHandler.HandleFinePay),
	)
	http.HandleFunc(
		"/fines/waive/",
		apiLogger.Middleware(fineHandler.HandleFineWaive),
	)
	http.HandleFunc(
		"/libraries",
		apiLogger.Middleware(libraryHandler.HandleLibraries),
	)
	http.HandleFunc(
		"/libraries/",
		apiLogger.Middleware(libraryHandler.HandleLibraryByID),
	)
	http.HandleFunc(
		"/loans",
		apiLogger.Middleware(loanHandler.HandleLoans),
	)
	http.HandleFunc(
		"/loans/",
		apiLogger.Middleware(loanHandler.HandleLoanByID),
	)
	http.HandleFunc(
		"/loans/renew/",
		apiLogger.Middleware(loanHandler.HandleLoanRenew),
	)
	http.HandleFunc(
		"/loans/return/",
		apiLogger.Middleware(loanHandler.HandleLoanReturn),
	)
	http.HandleFunc(
		"/publishers",
		apiLogger.Middleware(publisherHandler.HandlePublishers),
	)
	http.HandleFunc(
		"/publishers/",
		apiLogger.Middleware(publisherHandler.HandlePublisherByID),
	)
	http.HandleFunc(
		"/reservations",
		apiLogger.Middleware(reservationHandler.HandleReservations),
	)
	http.HandleFunc(
		"/reservations/",
		apiLogger.Middleware(reservationHandler.HandleReservationByID),
	)
	http.HandleFunc(
		"/reservations/cancel/",
		apiLogger.Middleware(reservationHandler.HandleReservationCancel),
	)
	http.HandleFunc(
		"/reservations/process/",
		apiLogger.Middleware(reservationHandler.HandleReservationProcess),
	)
	http.HandleFunc(
		"/shelves",
		apiLogger.Middleware(shelfHandler.HandleShelves),
	)
	http.HandleFunc(
		"/shelves/",
		apiLogger.Middleware(shelfHandler.HandleShelfByID),
	)
	http.HandleFunc(
		"/users",
		apiLogger.Middleware(userHandler.HandleUsers),
	)
	http.HandleFunc(
		"/users/",
		apiLogger.Middleware(userHandler.HandleUserByID),
	)
	http.HandleFunc(
		"/zones",
		apiLogger.Middleware(zoneHandler.HandleZones),
	)
	http.HandleFunc(
		"/zones/",
		apiLogger.Middleware(zoneHandler.HandleZoneByID),
	)

	fmt.Println("Servidor escuchando en el puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
