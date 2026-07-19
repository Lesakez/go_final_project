package server

import (
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
)

const defaultPort = "7540"

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	api.Init()
	http.Handle("/", http.FileServer(http.Dir("web")))

	log.Printf("Server started on http://localhost:%s", port)
	return http.ListenAndServe(":"+port, nil)
}
