package server

import (
	"log"
	"net/http"
	"os"

	"Graduate-work/pkg/api"

	"github.com/go-chi/chi/v5"
)

/*
StartServer initializes and starts the HTTP server with routing and static file serving.

Behavior:
  - Determines the port to bind to: uses the TODO_PORT environment variable if set;
    otherwise, defaults to port 7540.
  - Creates a new Chi router and initializes the API routes via api.Init(r).
  - Locates the directory for static web files:
  - Starts by checking "./web".
  - If it doesn’t exist, checks "../web" as a fallback.
  - Logs the chosen directory for static files.
  - Sets up a file server to serve static assets under all paths (/*).
  - Starts the HTTP server on the determined port.
  - Exits with a fatal log message if the server fails to start.
*/
func StartServer(cfgPort string, apiHandler *api.Handler) error {
	port := "7540"
	if cfgPort != "" {
		port = cfgPort
	}

	r := chi.NewRouter()

	apiHandler.RegisterRoutes(r)

	webDir := "./web"
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		if _, err := os.Stat("../web"); err == nil {
			webDir = "../web"
		}
	}

	log.Printf("Serving static files from: %s", webDir)
	fs := http.FileServer(http.Dir(webDir))

	r.Handle("/*", fs)

	log.Printf("Starting server on port %s", port)
	return http.ListenAndServe(":"+port, r)
}
