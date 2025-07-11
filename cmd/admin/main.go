package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Database types and methods (copied from main package for admin tool)
type Database struct {
	db *sql.DB
}

type MOTD struct {
	ID        int
	Content   string
	UpdatedAt time.Time
	UpdatedBy string
}

type ChatRoom struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("sqlite3", "../../bbs.db")
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) GetMOTD() (*MOTD, error) {
	var motd MOTD
	err := d.db.QueryRow("SELECT id, content, updated_at, updated_by FROM motd ORDER BY updated_at DESC LIMIT 1").
		Scan(&motd.ID, &motd.Content, &motd.UpdatedAt, &motd.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &motd, nil
}

func (d *Database) SetMOTD(content, updatedBy string) error {
	_, err := d.db.Exec("INSERT INTO motd (content, updated_by) VALUES (?, ?)", content, updatedBy)
	return err
}

func (d *Database) GetChatRooms() ([]ChatRoom, error) {
	rows, err := d.db.Query("SELECT id, name, description, created_at FROM chat_rooms ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []ChatRoom
	for rows.Next() {
		var room ChatRoom
		if err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.CreatedAt); err != nil {
			continue
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (d *Database) CreateChatRoom(name, description string) error {
	_, err := d.db.Exec("INSERT OR IGNORE INTO chat_rooms (name, description) VALUES (?, ?)", name, description)
	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func main() {
	// Initialize database
	db, err := NewDatabase()
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		fmt.Println("Make sure you run this from the BBS root directory and the database exists.")
		return
	}
	defer db.Close()

	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Println("=== BBS Admin Tool ===")
	fmt.Println("Commands: motd, room, users, help, quit")
	
	for {
		fmt.Print("admin> ")
		if !scanner.Scan() {
			break
		}
		
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		
		parts := strings.Fields(input)
		command := strings.ToLower(parts[0])
		
		switch command {
		case "help":
			showAdminHelp()
		case "motd":
			handleMOTD(db, scanner)
		case "room":
			handleRoom(db, scanner, parts[1:])
		case "users":
			handleUsers(db)
		case "quit", "exit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func showAdminHelp() {
	help := `
Available Admin Commands:
  motd        - Update the Message of the Day
  room        - Manage chat rooms (room list, room create)
  users       - List all registered users
  help        - Show this help message
  quit/exit   - Exit admin tool

Examples:
  motd                    - Update MOTD interactively
  room list               - List all chat rooms
  room create             - Create a new chat room
  users                   - Show all registered users
`
	fmt.Println(help)
}

func handleMOTD(db *Database, scanner *bufio.Scanner) {
	// Show current MOTD
	if motd, err := db.GetMOTD(); err == nil {
		fmt.Println("\nCurrent MOTD:")
		fmt.Println("=" + strings.Repeat("=", 50))
		fmt.Println(motd.Content)
		fmt.Println("=" + strings.Repeat("=", 50))
		fmt.Printf("Last updated: %s by %s\n\n", motd.UpdatedAt.Format("2006-01-02 15:04:05"), motd.UpdatedBy)
	}
	
	fmt.Println("Enter new MOTD (type 'END' on a line by itself to finish):")
	
	var lines []string
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			return
		}
		
		line := scanner.Text()
		if line == "END" {
			break
		}
		
		lines = append(lines, line)
	}
	
	if len(lines) == 0 {
		fmt.Println("MOTD not updated (empty content).")
		return
	}
	
	newMOTD := strings.Join(lines, "\n")
	
	if err := db.SetMOTD(newMOTD, "Admin"); err != nil {
		fmt.Printf("Failed to update MOTD: %v\n", err)
		return
	}
	
	fmt.Println("MOTD updated successfully!")
}

func handleRoom(db *Database, scanner *bufio.Scanner, args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: room <list|create>")
		return
	}
	
	subcommand := strings.ToLower(args[0])
	
	switch subcommand {
	case "list":
		rooms, err := db.GetChatRooms()
		if err != nil {
			fmt.Printf("Failed to get chat rooms: %v\n", err)
			return
		}
		
		fmt.Println("\nChat Rooms:")
		fmt.Println("=" + strings.Repeat("=", 60))
		for _, room := range rooms {
			fmt.Printf("ID: %d | Name: %s | Description: %s\n", room.ID, room.Name, room.Description)
			fmt.Printf("Created: %s\n", room.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println(strings.Repeat("-", 60))
		}
		
	case "create":
		fmt.Print("Room name: ")
		if !scanner.Scan() {
			return
		}
		name := strings.TrimSpace(scanner.Text())
		
		if name == "" {
			fmt.Println("Room name cannot be empty.")
			return
		}
		
		fmt.Print("Room description: ")
		if !scanner.Scan() {
			return
		}
		description := strings.TrimSpace(scanner.Text())
		
		if description == "" {
			description = "No description provided"
		}
		
		if err := db.CreateChatRoom(name, description); err != nil {
			fmt.Printf("Failed to create room: %v\n", err)
			return
		}
		
		fmt.Printf("Chat room '%s' created successfully!\n", name)
		
	default:
		fmt.Println("Usage: room <list|create>")
	}
}

func handleUsers(db *Database) {
	rows, err := db.db.Query("SELECT id, username, joined_at, last_seen FROM users ORDER BY username")
	if err != nil {
		fmt.Printf("Failed to get users: %v\n", err)
		return
	}
	defer rows.Close()
	
	fmt.Println("\nRegistered Users:")
	fmt.Println("=" + strings.Repeat("=", 80))
	fmt.Printf("%-5s | %-20s | %-20s | %-20s\n", "ID", "Username", "Joined", "Last Seen")
	fmt.Println(strings.Repeat("-", 80))
	
	for rows.Next() {
		var id int
		var username string
		var joinedAt, lastSeen string
		
		if err := rows.Scan(&id, &username, &joinedAt, &lastSeen); err != nil {
			continue
		}
		
		fmt.Printf("%-5d | %-20s | %-20s | %-20s\n", id, username, joinedAt[:19], lastSeen[:19])
	}
}