package main

import (
	"log"
	"net/http"
	"printer-service/config"
	"printer-service/websocket"
)

func main() {
	cfg := config.Load()

	// WebSocket endpoint
	http.HandleFunc("/ws", websocket.HandleWebSocket)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Printer service is running"))
	})

	log.Printf("Starting printer service on %s", cfg.WebSocketPort)
	log.Printf("WebSocket available at ws://localhost%s/ws", cfg.WebSocketPort)

	if err := http.ListenAndServe(cfg.WebSocketPort, nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
