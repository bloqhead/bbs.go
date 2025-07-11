# Enhanced BBS (Bulletin Board System)

A modern, full-featured BBS written in Go with chat rooms, user authentication, message persistence, and real-time messaging capabilities.

## Features

- **Multi-User Chat Rooms**: Multiple themed chat rooms (General, Tech, Gaming, Random)
- **User Authentication**: Secure user registration and login with bcrypt password hashing
- **Message Persistence**: SQLite database storage for users, messages, and chat history
- **Message of the Day (MOTD)**: Customizable welcome message for all users
- **Real-Time Messaging**: Live chat with other online users
- **ANSI Color Support**: Rich text formatting and colors in terminal
- **User Management**: Track online users, join/leave notifications
- **Message History**: View recent message history in chat rooms

## Requirements

- Go 1.21 or later
- SQLite3 (automatically handled by Go driver)

## Installation & Setup

1. **Clone or download the project:**
   ```bash
   git clone <repository-url>
   cd bbs
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Build the BBS:**
   ```bash
   go build -o bbs
   ```

4. **Run the BBS server:**
   ```bash
   ./bbs
   ```

The server will start on port 3003 and automatically create the SQLite database (`bbs.db`) with default chat rooms and MOTD.

## Connecting to the BBS

Use any telnet client or netcat to connect to the BBS:

```bash
# Using telnet (recommended)
telnet localhost 3003

# Using netcat (nc) - if telnet is not available
nc localhost 3003
```

### Client Options:
- **telnet** (recommended): Best experience with proper line handling
- **nc (netcat)**: Works well for systems without telnet
- **Windows**: Use built-in telnet or PuTTY
- **macOS/Linux**: Built-in telnet command or nc
- **Modern terminals**: Most terminal emulators support telnet

### Troubleshooting Connection Issues:
If you experience input issues with nc, try:
```bash
# For debug mode to see raw input
./bbs debug
# Then connect with: nc localhost 3004
```

## Usage

### First Time Setup
1. Connect to the BBS via telnet
2. Choose "R" to register a new account
3. Create a username (minimum 3 characters)
4. Create a password (minimum 4 characters)
5. You'll be automatically logged in and joined to the General chat room

### Commands

Once logged in, you can use these commands:

- `help` - Show available commands
- `rooms` - List all available chat rooms
- `join <room>` - Join a specific chat room (e.g., `join Tech`)
- `msg <message>` - Send a message to current room
- `users` - List currently online users
- `history` - Show recent message history for current room
- `motd` - Display the message of the day
- `quit` or `exit` - Leave the BBS

### Quick Messaging
You can also send messages directly without the `msg` command:
```
[General]> Hello everyone!
```

## Default Chat Rooms

The BBS comes with four default chat rooms:

- **General** - General discussion for all users
- **Tech** - Technology and programming discussions  
- **Gaming** - Video games and gaming culture
- **Random** - Random topics and casual chat

## Database Structure

The BBS uses SQLite with the following tables:

- **users** - User accounts with encrypted passwords
- **chat_rooms** - Available chat rooms
- **messages** - Chat message history
- **motd** - Message of the day entries

## Architecture

The BBS is built with a modular architecture:

- `main.go` - Entry point and server initialization
- `server.go` - Multi-client server management and broadcasting
- `client.go` - Individual client session handling
- `database.go` - Database operations and schema management

## Advanced Features

### Multi-User Support
- Concurrent user connections
- Real-time message broadcasting
- User presence notifications (join/leave)
- Thread-safe client management

### Message Broadcasting
- Messages are broadcast in real-time to all users in the same chat room
- Users see join/leave notifications
- Timestamp display for all messages

### Security
- Passwords are hashed using bcrypt
- Input validation and sanitization
- SQL injection protection with prepared statements

## Customization

### Adding New Chat Rooms
You can add new chat rooms by modifying the `createDefaultData()` function in `database.go` or by directly inserting into the database:

```sql
INSERT INTO chat_rooms (name, description) VALUES ('NewRoom', 'Description here');
```

### Updating MOTD
The Message of the Day can be updated by inserting a new record into the `motd` table:

```sql
INSERT INTO motd (content, updated_by) VALUES ('New MOTD content here', 'Admin');
```

## Development

### Building for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bbs-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o bbs-windows.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o bbs-macos
```

### Running in Development Mode

```bash
go run *.go
```

## Troubleshooting

### Common Issues

1. **Port already in use**: Change the port in `main.go` if 3003 is occupied
2. **Database locked**: Ensure no other instances are running
3. **Connection refused**: Check if the server is running and port is accessible

### Logs
The server outputs connection logs and user activity to stdout. Monitor these for debugging.

## Contributing

This is a learning project demonstrating:
- Go networking with goroutines
- SQLite database integration
- Real-time message broadcasting
- Terminal-based user interfaces
- Multi-user system design

Feel free to extend with additional features like:
- Private messaging
- File uploads/downloads
- User roles and permissions
- Message search functionality
- Web interface

## License

This project is for educational purposes. Use and modify as needed for learning Go networking and database programming.