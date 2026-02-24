package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"cloud.google.com/go/firestore"
)

var firestoreClient *firestore.Client

func main() {
	// Get PORT from environment (required by Cloud Run)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Firestore client
	projectID := os.Getenv("GCP_PROJECT")
	databaseID := os.Getenv("FIRESTORE_DATABASE")
	if projectID != "" && databaseID != "" {
		ctx := context.Background()
		var err error
		firestoreClient, err = firestore.NewClientWithDatabase(ctx, projectID, databaseID)
		if err != nil {
			log.Printf("Warning: Failed to initialize Firestore client: %v", err)
		} else {
			defer firestoreClient.Close()
			log.Printf("Connected to Firestore database: %s", databaseID)
		}
	}

	// Set up HTTP routes
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/write", handleWrite)
	http.HandleFunc("/read", handleRead)

	// Start server
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Handle graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Cloud Run! 🚀\n")
	fmt.Fprintf(w, "Version: 1.0.0\n")
	fmt.Fprintf(w, "Firestore: %s\n", getFirestoreStatus())
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func handleWrite(w http.ResponseWriter, r *http.Request) {
	if firestoreClient == nil {
		http.Error(w, "Firestore client not initialized", http.StatusServiceUnavailable)
		return
	}

	ctx := context.Background()
	docRef := firestoreClient.Collection("test").Doc("sample")
	_, err := docRef.Set(ctx, map[string]interface{}{
		"message":   "Hello from Cloud Run",
		"timestamp": time.Now(),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to write to Firestore: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully wrote to Firestore!\n")
}

func handleRead(w http.ResponseWriter, r *http.Request) {
	if firestoreClient == nil {
		http.Error(w, "Firestore client not initialized", http.StatusServiceUnavailable)
		return
	}

	ctx := context.Background()
	docRef := firestoreClient.Collection("test").Doc("sample")
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read from Firestore: %v", err), http.StatusInternalServerError)
		return
	}

	data := docSnap.Data()
	fmt.Fprintf(w, "Data from Firestore: %v\n", data)
}

func getFirestoreStatus() string {
	if firestoreClient != nil {
		return "Connected"
	}
	return "Not configured"
}
