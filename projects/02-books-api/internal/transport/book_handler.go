package transport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
	"github.com/chicho69-cesar/backend-go/books/internal/services"
	"github.com/chicho69-cesar/backend-go/books/internal/store"
)

type BookHandler struct {
	bookService *services.BookService
}

func NewBookHandler(bookService *services.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

// GET /books - Obtener todos los libros o con filtros
// POST /books - Crear un nuevo libro
func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			var hasFilters bool
			filter := store.BookFilter{}

			isbn := r.URL.Query().Get("isbn")
			if isbn != "" {
				filter.ISBN = isbn
				hasFilters = true
			}

			shelfIDStr := r.URL.Query().Get("shelf_id")
			if shelfIDStr != "" {
				shelfID, err := strconv.ParseInt(shelfIDStr, 10, 64)
				if err != nil || shelfID <= 0 {
					http.Error(w, "El ID del estante es inválido", http.StatusBadRequest)
					return
				}

				filter.ShelfID = &shelfID
				hasFilters = true
			}

			authorIDStr := r.URL.Query().Get("author_id")
			if authorIDStr != "" {
				authorID, err := strconv.ParseInt(authorIDStr, 10, 64)
				if err != nil || authorID <= 0 {
					http.Error(w, "El ID del autor es inválido", http.StatusBadRequest)
					return
				}

				filter.AuthorID = &authorID
				hasFilters = true
			}

			categoryIDStr := r.URL.Query().Get("category_id")
			if categoryIDStr != "" {
				categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
				if err != nil || categoryID <= 0 {
					http.Error(w, "El ID de la categoría es inválido", http.StatusBadRequest)
					return
				}

				filter.CategoryID = &categoryID
				hasFilters = true
			}

			var books []*models.Book
			var err error

			if hasFilters {
				books, err = h.bookService.GetBooksFiltered(filter)
			} else {
				books, err = h.bookService.GetAllBooks()
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(books)

		case http.MethodPost:
			var book models.Book
			err := json.NewDecoder(r.Body).Decode(&book)
			if err != nil {
				http.Error(w, "Datos de libro inválidos", http.StatusBadRequest)
				return
			}

			createdBook, err := h.bookService.CreateBook(&book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createdBook)

		default:
			http.Error(w, "Unavailable method", http.StatusMethodNotAllowed)
	}
}

// GET /books/{id} - Obtener un libro por ID
// PUT /books/{id} - Actualizar un libro por ID
// DELETE /books/{id} - Eliminar un libro por ID
func (h *BookHandler) HandleBookByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/books/")
	parts := strings.Split(path, "/")
	
	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "El parámetro ID es requerido", http.StatusBadRequest)
		return
	}

	idParam := parts[0]
	readId, err := strconv.Atoi(idParam)
	if err != nil || readId <= 0 {
		http.Error(w, "El ID es inválido", http.StatusBadRequest)
		return
	}

	id := int64(readId)

	if len(parts) > 1 {
		switch parts[1] {
			case "authors":
				h.handleBookAuthors(w, r, id, parts[2:])
				return
			case "categories":
				h.handleBookCategories(w, r, id, parts[2:])
				return
			default:
				http.Error(w, "Ruta no encontrada", http.StatusNotFound)
				return
		}
	}

	switch r.Method {
		case http.MethodGet:
			book, err := h.bookService.GetBookByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)

		case http.MethodPut:
			var book models.Book
			err := json.NewDecoder(r.Body).Decode(&book)
			if err != nil {
				http.Error(w, "Datos de libro inválidos", http.StatusBadRequest)
				return
			}

			updatedBook, err := h.bookService.UpdateBook(id, &book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedBook)

		case http.MethodDelete:
			err := h.bookService.DeleteBook(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Unavailable method", http.StatusMethodNotAllowed)
	}
}

// GET /books/{id}/authors - Obtener autores del libro
// POST /books/{id}/authors - Agregar autor al libro
// DELETE /books/{id}/authors/{authorId} - Eliminar autor del libro
// PUT /books/{id}/authors/{authorId} - Actualizar posición del autor
func (h *BookHandler) handleBookAuthors(w http.ResponseWriter, r *http.Request, bookID int64, parts []string) {
	switch r.Method {
		case http.MethodGet:
			authors, err := h.bookService.GetBookAuthors(bookID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(authors)

		case http.MethodPost:
			var bookAuthor models.BookAuthor
			err := json.NewDecoder(r.Body).Decode(&bookAuthor)
			if err != nil {
				http.Error(w, "Datos inválidos", http.StatusBadRequest)
				return
			}

			bookAuthor.BookID = bookID

			err = h.bookService.AddAuthorToBook(&bookAuthor)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Autor agregado exitosamente",
			})

		case http.MethodDelete:
			if len(parts) == 0 || parts[0] == "" {
				http.Error(w, "El ID del autor es requerido", http.StatusBadRequest)
				return
			}

			authorID, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil || authorID <= 0 {
				http.Error(w, "El ID del autor es inválido", http.StatusBadRequest)
				return
			}

			err = h.bookService.RemoveAuthorFromBook(bookID, authorID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		case http.MethodPut:
			if len(parts) == 0 || parts[0] == "" {
				http.Error(w, "El ID del autor es requerido", http.StatusBadRequest)
				return
			}

			authorID, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil || authorID <= 0 {
				http.Error(w, "El ID del autor es inválido", http.StatusBadRequest)
				return
			}

			var data struct {
				Position int `json:"position"`
			}

			err = json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, "Datos inválidos", http.StatusBadRequest)
				return
			}

			err = h.bookService.UpdateAuthorPosition(bookID, authorID, data.Position)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Posición del autor actualizada exitosamente",
			})

		default:
			http.Error(w, "Unavailable method", http.StatusMethodNotAllowed)
	}
}

// GET /books/{id}/categories - Obtener categorías del libro
// POST /books/{id}/categories - Agregar categoría al libro
// DELETE /books/{id}/categories/{categoryId} - Eliminar categoría del libro
func (h *BookHandler) handleBookCategories(w http.ResponseWriter, r *http.Request, bookID int64, parts []string) {
	switch r.Method {
		case http.MethodGet:
			categories, err := h.bookService.GetBookCategories(bookID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)

		case http.MethodPost:
			var bookCategory models.BookCategory
			err := json.NewDecoder(r.Body).Decode(&bookCategory)
			if err != nil {
				http.Error(w, "Datos inválidos", http.StatusBadRequest)
				return
			}

			bookCategory.BookID = bookID

			err = h.bookService.AddCategoryToBook(&bookCategory)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Categoría agregada exitosamente",
			})

		case http.MethodDelete:
			if len(parts) == 0 || parts[0] == "" {
				http.Error(w, "El ID de la categoría es requerido", http.StatusBadRequest)
				return
			}

			categoryID, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil || categoryID <= 0 {
				http.Error(w, "El ID de la categoría es inválido", http.StatusBadRequest)
				return
			}

			err = h.bookService.RemoveCategoryFromBook(bookID, categoryID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Unavailable method", http.StatusMethodNotAllowed)
	}
}
