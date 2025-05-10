// @title DevDistillery API
// @version 1.0
// @description API para la gestión de proyectos de desarrollo
// @host localhost:8081
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"func/cmd/server"
	"log"
	"net/http"
	"time"
)

func main() {
	server := &http.Server{
		Addr:         ":8081",
		Handler:      server.New(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Servidor iniciado en http://localhost:8081")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
