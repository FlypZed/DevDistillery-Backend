package main

import (
	"func/cmd/server"
	"log"
	"net/http"
)

func main() {
	application := server.New()

	log.Println("Servidor iniciado en http://localhost:8081")
	if err := http.ListenAndServe(":8081", application); err != nil {
		log.Panic("Error al iniciar el servidor: ", err)
	}
}
