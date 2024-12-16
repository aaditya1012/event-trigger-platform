package main

import (
	"log"
	"net/http"

	"event-trigger-platform/db"
	"event-trigger-platform/routes"

	_ "event-trigger-platform/docs" // Import generated Swagger docs
)

// @title Event Trigger Platform API
// @version 1.0
// @description API documentation for managing event triggers and logs.
// @termsOfService http://example.com/terms/

// @contact.name Developer Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @openapi 3.0.0
func main() {
	// Initialize database connection
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// Define HTTP routes for triggers
	http.HandleFunc("/triggers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			routes.CreateTrigger(w, r)
		} else if r.Method == http.MethodGet {
			routes.ListTriggers(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/triggers/delete", routes.DeleteTrigger)

	// Define HTTP route for logs
	http.HandleFunc("/logs", routes.ListLogs)

	// Serve Swagger UI static files
	fs := http.FileServer(http.Dir("./swagger-ui/dist"))
	http.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	// Serve swagger.json from the /docs directory
	http.Handle("/swagger.json", http.FileServer(http.Dir("./docs")))

	// Start event processing in the background
	go routes.StartEventProcessing()

	// Start the HTTP server
	log.Println("Server started on :8080")
	log.Println("Swagger UI available at http://localhost:8080/swagger/index.html")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
