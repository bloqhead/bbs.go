package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	conn         net.Conn
	user         *User
	currentRoom  *ChatRoom
	db           *Database
	server       *BBSServer
	authenticated bool
	scanner      *bufio.Scanner
}

func NewClient(conn net.Conn, db *Database, server *BBSServer) *Client {
	return &Client{
		conn:    conn,
		db:      db,
		server:  server,
		scanner: bufio.NewScanner(conn),
	}
}

func (c *Client) Handle() {
	defer func() {
		if c.user != nil {
			c.server.RemoveClient(c)
		}
		c.conn.Close()
	}()

	// Display logo and welcome message
	c.displayWelcome()

	// Authentication loop
	if !c.authenticate() {
		return
	}

	// Display MOTD
	c.displayMOTD()

	// Join default room
	if room, err := c.db.GetChatRoom("General"); err == nil {
		c.currentRoom = room
		c.write(fmt.Sprintf("\n\033[32mJoined chat room: %s\033[0m\n", room.Name))
		c.displayRecentMessages()
	}

	// Add client to server
	c.server.AddClient(c)

	// Main command loop
	c.commandLoop()
}

func (c *Client) displayWelcome() {
	// Display logo
	if logo, err := readLogo(); err == nil {
		c.write(string(logo))
	}

	// Clear screen and set color
	c.write("\033[H\033[2J")
	c.write("\033[33mWelcome to the Enhanced BBS!\033[0m\n\n")
}

func (c *Client) authenticate() bool {
	for {
		c.write("Do you want to (L)ogin or (R)egister? ")
		if !c.scanner.Scan() {
			return false
		}

		choice := strings.ToLower(strings.TrimSpace(c.scanner.Text()))
		switch choice {
		case "l", "login":
			if c.login() {
				return true
			}
		case "r", "register":
			if c.register() {
				return true
			}
		case "q", "quit":
			return false
		default:
			c.write("Please enter 'L' for login, 'R' for register, or 'Q' to quit.\n")
		}
	}
}

func (c *Client) login() bool {
	c.write("Username: ")
	if !c.scanner.Scan() {
		return false
	}
	username := strings.TrimSpace(c.scanner.Text())

	c.write("Password: ")
	if !c.scanner.Scan() {
		return false
	}
	password := strings.TrimSpace(c.scanner.Text())

	user, err := c.db.AuthenticateUser(username, password)
	if err != nil {
		c.write("\033[31mInvalid username or password.\033[0m\n\n")
		return false
	}

	c.user = user
	c.authenticated = true
	c.write(fmt.Sprintf("\033[32mWelcome back, %s!\033[0m\n\n", user.Username))
	return true
}

func (c *Client) register() bool {
	c.write("Choose a username: ")
	if !c.scanner.Scan() {
		return false
	}
	username := strings.TrimSpace(c.scanner.Text())

	if len(username) < 3 {
		c.write("\033[31mUsername must be at least 3 characters long.\033[0m\n\n")
		return false
	}

	c.write("Choose a password: ")
	if !c.scanner.Scan() {
		return false
	}
	password := strings.TrimSpace(c.scanner.Text())

	if len(password) < 4 {
		c.write("\033[31mPassword must be at least 4 characters long.\033[0m\n\n")
		return false
	}

	if err := c.db.CreateUser(username, password); err != nil {
		c.write("\033[31mUsername already exists or registration failed.\033[0m\n\n")
		return false
	}

	// Auto-login after registration
	user, err := c.db.AuthenticateUser(username, password)
	if err != nil {
		c.write("\033[31mRegistration succeeded but login failed.\033[0m\n\n")
		return false
	}

	c.user = user
	c.authenticated = true
	c.write(fmt.Sprintf("\033[32mAccount created! Welcome, %s!\033[0m\n\n", user.Username))
	return true
}

func (c *Client) displayMOTD() {
	if motd, err := c.db.GetMOTD(); err == nil {
		c.write("\033[36m" + strings.Repeat("=", 60) + "\033[0m\n")
		c.write("\033[36mMESSAGE OF THE DAY\033[0m\n")
		c.write("\033[36m" + strings.Repeat("=", 60) + "\033[0m\n")
		c.write(motd.Content + "\n")
		c.write("\033[36m" + strings.Repeat("=", 60) + "\033[0m\n\n")
	}
}

func (c *Client) displayRecentMessages() {
	if c.currentRoom == nil {
		return
	}

	messages, err := c.db.GetRecentMessages(c.currentRoom.ID, 10)
	if err != nil || len(messages) == 0 {
		c.write("No recent messages in this room.\n\n")
		return
	}

	c.write(fmt.Sprintf("\033[35mRecent messages in %s:\033[0m\n", c.currentRoom.Name))
	c.write(strings.Repeat("-", 40) + "\n")
	
	for _, msg := range messages {
		timestamp := msg.Timestamp.Format("15:04")
		c.write(fmt.Sprintf("\033[90m[%s]\033[0m \033[33m%s:\033[0m %s\n", timestamp, msg.Username, msg.Content))
	}
	c.write(strings.Repeat("-", 40) + "\n\n")
}

func (c *Client) commandLoop() {
	c.write(fmt.Sprintf("\033[32mType 'help' for commands. Current room: %s\033[0m\n", c.currentRoom.Name))
	
	for {
		c.write(fmt.Sprintf("\033[34m[%s]>\033[0m ", c.currentRoom.Name))
		
		if !c.scanner.Scan() {
			break
		}

		input := strings.TrimSpace(c.scanner.Text())
		if input == "" {
			continue
		}

		if !c.handleCommand(input) {
			break
		}
	}
}

func (c *Client) handleCommand(input string) bool {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return true
	}

	command := strings.ToLower(parts[0])
	args := parts[1:]

	switch command {
	case "help":
		c.showHelp()
	case "rooms":
		c.listRooms()
	case "join":
		if len(args) > 0 {
			c.joinRoom(strings.Join(args, " "))
		} else {
			c.write("Usage: join <room_name>\n")
		}
	case "users":
		c.listUsers()
	case "msg":
		if len(args) > 0 {
			c.sendMessage(strings.Join(args, " "))
		} else {
			c.write("Usage: msg <your_message>\n")
		}
	case "history":
		c.showHistory()
	case "motd":
		c.displayMOTD()
	case "quit", "exit":
		c.write("Goodbye!\n")
		return false
	default:
		// If it doesn't start with a command, treat it as a message
		c.sendMessage(input)
	}

	return true
}

func (c *Client) showHelp() {
	help := `
\033[36mAvailable Commands:\033[0m
  help                 - Show this help message
  rooms                - List all available chat rooms
  join <room>          - Join a specific chat room
  msg <message>        - Send a message to current room
  users                - List users currently online
  history              - Show recent message history
  motd                 - Display message of the day
  quit/exit            - Leave the BBS

\033[36mQuick messaging:\033[0m
  You can also just type your message directly without 'msg'

\033[36mNavigation:\033[0m
  Your current room is shown in the prompt: [RoomName]>
`
	c.write(help + "\n")
}

func (c *Client) listRooms() {
	rooms, err := c.db.GetChatRooms()
	if err != nil {
		c.write("Error loading chat rooms.\n")
		return
	}

	c.write("\033[36mAvailable Chat Rooms:\033[0m\n")
	c.write(strings.Repeat("-", 50) + "\n")
	
	for _, room := range rooms {
		currentMarker := ""
		if c.currentRoom != nil && room.ID == c.currentRoom.ID {
			currentMarker = " \033[32m(current)\033[0m"
		}
		c.write(fmt.Sprintf("\033[33m%s\033[0m - %s%s\n", room.Name, room.Description, currentMarker))
	}
	c.write(strings.Repeat("-", 50) + "\n\n")
}

func (c *Client) joinRoom(roomName string) {
	room, err := c.db.GetChatRoom(roomName)
	if err != nil {
		c.write(fmt.Sprintf("Room '%s' not found.\n", roomName))
		return
	}

	if c.currentRoom != nil && room.ID == c.currentRoom.ID {
		c.write("You are already in this room.\n")
		return
	}

	c.currentRoom = room
	c.write(fmt.Sprintf("\033[32mJoined room: %s\033[0m\n", room.Name))
	c.displayRecentMessages()
}

func (c *Client) listUsers() {
	users := c.server.GetOnlineUsers()
	c.write(fmt.Sprintf("\033[36mOnline Users (%d):\033[0m\n", len(users)))
	c.write(strings.Repeat("-", 30) + "\n")
	
	for _, user := range users {
		currentMarker := ""
		if user == c.user.Username {
			currentMarker = " \033[32m(you)\033[0m"
		}
		c.write(fmt.Sprintf("\033[33m%s\033[0m%s\n", user, currentMarker))
	}
	c.write(strings.Repeat("-", 30) + "\n\n")
}

func (c *Client) sendMessage(content string) {
	if c.currentRoom == nil {
		c.write("You are not in a chat room.\n")
		return
	}

	if err := c.db.AddMessage(c.currentRoom.ID, c.user.ID, c.user.Username, content); err != nil {
		c.write("Failed to send message.\n")
		return
	}

	// Broadcast message to all clients in the same room
	timestamp := time.Now().Format("15:04")
	message := fmt.Sprintf("\033[90m[%s]\033[0m \033[33m%s:\033[0m %s\n", timestamp, c.user.Username, content)
	c.server.BroadcastToRoom(c.currentRoom.ID, message, c)
}

func (c *Client) showHistory() {
	if c.currentRoom == nil {
		c.write("You are not in a chat room.\n")
		return
	}
	c.displayRecentMessages()
}

func (c *Client) write(message string) {
	c.conn.Write([]byte(message))
}

func (c *Client) GetUser() *User {
	return c.user
}

func (c *Client) GetCurrentRoom() *ChatRoom {
	return c.currentRoom
}

func readLogo() ([]byte, error) {
	// Try to read the logo file, fallback to simple text if not available
	return []byte(`
     ____  ____  ____    ____            __                
    / __ )/ __ )/ __ \  / __ \          / /_____ _________
   / __  / __  / / / / / / / /_____   / __/ __ '/ ___/ __ \
  / /_/ / /_/ / /_/ / / /_/ /_____/  / /_/ /_/ / /  / /_/ /
 /_____/_____/\____/  \____/        \__/\__,_/_/   \____/ 

        Enhanced Bulletin Board System v2.0
`), nil
}