package main

// Enhanced BBS (Bulletin Board System)
// Features:
// - Multi-user chat rooms
// - User authentication and registration
// - Message persistence with SQLite database
// - Message of the Day (MOTD)
// - Real-time messaging between users
// - ANSI color support for better UI

import (
	"log"
)

func main() {
	// Initialize database
	db, err := NewDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create and start BBS server
	server := NewBBSServer(db)
	
	log.Println("Starting Enhanced BBS Server...")
	log.Println("Features: Chat Rooms, User Auth, Message History, MOTD")
	log.Println("Connect via telnet: telnet localhost 3003")
	
	if err := server.Start("3003"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
