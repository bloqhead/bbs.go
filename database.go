package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	db *sql.DB
}

type User struct {
	ID       int
	Username string
	Password string
	JoinedAt time.Time
	LastSeen time.Time
}

type ChatRoom struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
}

type Message struct {
	ID        int
	RoomID    int
	UserID    int
	Username  string
	Content   string
	Timestamp time.Time
}

type MOTD struct {
	ID        int
	Content   string
	UpdatedAt time.Time
	UpdatedBy string
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("sqlite3", "bbs.db")
	if err != nil {
		return nil, err
	}

	database := &Database{db: db}
	if err := database.createTables(); err != nil {
		return nil, err
	}

	// Create default chat rooms and MOTD
	database.createDefaultData()

	return database, nil
}

func (d *Database) createTables() error {
	// Users table
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_seen DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Chat rooms table
	roomTable := `
	CREATE TABLE IF NOT EXISTS chat_rooms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Messages table
	messageTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		room_id INTEGER,
		user_id INTEGER,
		username TEXT,
		content TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (room_id) REFERENCES chat_rooms(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	// MOTD table
	motdTable := `
	CREATE TABLE IF NOT EXISTS motd (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT
	);`

	tables := []string{userTable, roomTable, messageTable, motdTable}
	for _, table := range tables {
		if _, err := d.db.Exec(table); err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) createDefaultData() {
	// Create default chat rooms
	defaultRooms := []struct {
		name, description string
	}{
		{"General", "General discussion for all users"},
		{"Tech", "Technology and programming discussions"},
		{"Gaming", "Video games and gaming culture"},
		{"Random", "Random topics and casual chat"},
	}

	for _, room := range defaultRooms {
		d.CreateChatRoom(room.name, room.description)
	}

	// Create default MOTD
	defaultMOTD := `Welcome to the Enhanced BBS!

Features:
- Multiple chat rooms for different topics
- Real-time messaging with other users
- User registration and authentication
- Message history and user tracking

Commands:
- 'help' - Show available commands
- 'rooms' - List all chat rooms
- 'join <room>' - Join a chat room
- 'msg <message>' - Send message to current room
- 'users' - List online users
- 'quit' - Exit the BBS

Enjoy your stay!`

	d.SetMOTD(defaultMOTD, "System")
}

func (d *Database) CreateUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = d.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, string(hashedPassword))
	return err
}

func (d *Database) AuthenticateUser(username, password string) (*User, error) {
	var user User
	var hashedPassword string

	err := d.db.QueryRow("SELECT id, username, password, joined_at, last_seen FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &hashedPassword, &user.JoinedAt, &user.LastSeen)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, err
	}

	// Update last seen
	d.db.Exec("UPDATE users SET last_seen = CURRENT_TIMESTAMP WHERE id = ?", user.ID)

	return &user, nil
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

func (d *Database) GetChatRoom(name string) (*ChatRoom, error) {
	var room ChatRoom
	err := d.db.QueryRow("SELECT id, name, description, created_at FROM chat_rooms WHERE name = ?", name).
		Scan(&room.ID, &room.Name, &room.Description, &room.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (d *Database) AddMessage(roomID, userID int, username, content string) error {
	_, err := d.db.Exec("INSERT INTO messages (room_id, user_id, username, content) VALUES (?, ?, ?, ?)",
		roomID, userID, username, content)
	return err
}

func (d *Database) GetRecentMessages(roomID int, limit int) ([]Message, error) {
	rows, err := d.db.Query(`
		SELECT id, room_id, user_id, username, content, timestamp 
		FROM messages 
		WHERE room_id = ? 
		ORDER BY timestamp DESC 
		LIMIT ?`, roomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.RoomID, &msg.UserID, &msg.Username, &msg.Content, &msg.Timestamp); err != nil {
			continue
		}
		messages = append([]Message{msg}, messages...) // Reverse order for chronological display
	}

	return messages, nil
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

func (d *Database) Close() error {
	return d.db.Close()
}