package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Database connection
var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=postgres dbname=chatdb sslmode=disable")
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
}

// Struct untuk pesan chat
type Message struct {
	Username  string `json:"username"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp,omitempty"`
}

// List koneksi aktif
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var mutex = &sync.Mutex{}

// Handle koneksi WebSocket
func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer ws.Close()

	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Error membaca pesan:", err)
			mutex.Lock()
			delete(clients, ws)
			mutex.Unlock()
			break
		}

		// Simpan ke database dengan timestamp otomatis
		_, err = db.Exec("INSERT INTO messages (username, message, timestamp) VALUES ($1, $2, NOW())", msg.Username, msg.Message)
		if err != nil {
			log.Println("Error menyimpan pesan ke database:", err)
			continue
		}

		log.Println("Pesan diterima:", msg.Username, "-", msg.Message)
		broadcast <- msg
	}
}

// Handle pengiriman pesan ke semua client
func handleMessages() {
	for msg := range broadcast {
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Error mengirim pesan:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

// Ambil riwayat chat dari database
func getChatHistory(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT username, message, timestamp FROM messages ORDER BY timestamp DESC LIMIT 50")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.Username, &msg.Message, &msg.Timestamp)
		if err != nil {
			log.Println("Error membaca data dari database:", err)
			continue
		}
		messages = append(messages, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func main() {
	defer db.Close() // Tutup database saat program selesai

	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/history", getChatHistory)

	go handleMessages()

	fmt.Println("Server berjalan di port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil)) // Pakai `log.Fatal()` supaya error langsung keluar
}
