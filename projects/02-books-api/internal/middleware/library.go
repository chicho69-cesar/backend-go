package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type contextKey string

const LibraryIDKey contextKey = "library_id"

func ExtractLibraryID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		segments := strings.Split(path, "/")

		if len(segments) == 0 || segments[0] == "" {
			http.Error(w, "URL inválida: falta library_id", http.StatusBadRequest)
			return
		}

		libraryIDStr := segments[0]
		libraryID, err := strconv.ParseInt(libraryIDStr, 10, 64)
		if err != nil || libraryID <= 0 {
			http.Error(w, fmt.Sprintf("library_id inválido: %s", libraryIDStr), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), LibraryIDKey, libraryID)

		newPath := "/" + strings.Join(segments[1:], "/")
		r.URL.Path = newPath

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetLibraryID(r *http.Request) (int64, error) {
	libraryID, ok := r.Context().Value(LibraryIDKey).(int64)
	if !ok {
		return 0, fmt.Errorf("library_id no encontrado en el contexto")
	}

	return libraryID, nil
}
