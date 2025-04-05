// @title DevDistillery API
// @version 1.0
// @description
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"func/cmd/server"
	"log"
	"net/http"
)

// Middleware para manejar CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" || origin == "http://localhost:8080" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	application := server.New()

	handler := corsMiddleware(application)

	log.Println("Servidor iniciado en http://localhost:8081")
	if err := http.ListenAndServe(":8081", handler); err != nil {
		log.Panic("Error al iniciar el servidor: ", err)
	}
}
