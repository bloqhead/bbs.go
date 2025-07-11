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
	"os"
)

func main() {
	// Check for debug mode
	if len(os.Args) > 1 && os.Args[1] == "debug" {
		// Initialize minimal database for debug mode
		db, err := NewDatabase()
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
		defer db.Close()

		server := NewBBSServer(db)
		log.Println("Starting BBS Debug Server...")
		log.Println("This will show raw input to help debug nc vs telnet issues")
		log.Println("Connect with: nc localhost 3004")
		log.Println("Or: telnet localhost 3004")
		
		if err := server.StartDebug("3004"); err != nil {
			log.Fatalf("Debug server error: %v", err)
		}
		return
	}

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
	log.Println("Connect via nc: nc localhost 3003")
	log.Println("For debug mode: ./bbs debug")
	
	if err := server.Start("3003"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
