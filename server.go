package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type BBSServer struct {
	db      *Database
	clients map[*Client]bool
	mutex   sync.RWMutex
}

func NewBBSServer(db *Database) *BBSServer {
	return &BBSServer{
		db:      db,
		clients: make(map[*Client]bool),
	}
}

func (s *BBSServer) Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("BBS Server started on port %s", port)

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Shutdown signal received, closing server...")
		cancel()
		listener.Close()
	}()

	// Accept connections
	for {
		select {
		case <-ctx.Done():
			log.Println("Server shutdown completed")
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				if ctx.Err() != nil {
					return nil // Server is shutting down
				}
				log.Printf("Error accepting connection: %v", err)
				continue
			}

			// Handle client in a new goroutine
			go func() {
				client := NewClient(conn, s.db, s)
				client.Handle()
			}()
		}
	}
}

func (s *BBSServer) AddClient(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	s.clients[client] = true
	log.Printf("User %s connected. Total online: %d", client.user.Username, len(s.clients))
	
	// Notify other users in the same room
	if client.currentRoom != nil {
		message := fmt.Sprintf("\033[90m*** %s joined the room ***\033[0m\n", client.user.Username)
		s.broadcastToRoomExcluding(client.currentRoom.ID, message, client)
	}
}

func (s *BBSServer) RemoveClient(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, exists := s.clients[client]; exists {
		delete(s.clients, client)
		log.Printf("User %s disconnected. Total online: %d", client.user.Username, len(s.clients))
		
		// Notify other users in the same room
		if client.currentRoom != nil {
			message := fmt.Sprintf("\033[90m*** %s left the room ***\033[0m\n", client.user.Username)
			s.broadcastToRoomExcluding(client.currentRoom.ID, message, client)
		}
	}
}

func (s *BBSServer) BroadcastToRoom(roomID int, message string, sender *Client) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	for client := range s.clients {
		if client != sender && client.currentRoom != nil && client.currentRoom.ID == roomID {
			client.write(message)
		}
	}
}

func (s *BBSServer) broadcastToRoomExcluding(roomID int, message string, excludeClient *Client) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	for client := range s.clients {
		if client != excludeClient && client.currentRoom != nil && client.currentRoom.ID == roomID {
			client.write(message)
		}
	}
}

func (s *BBSServer) GetOnlineUsers() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	var users []string
	for client := range s.clients {
		if client.user != nil {
			users = append(users, client.user.Username)
		}
	}
	
	return users
}

func (s *BBSServer) GetClientCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.clients)
}

func (s *BBSServer) BroadcastGlobal(message string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	for client := range s.clients {
		client.write(message)
	}
}

func (s *BBSServer) GetClientsInRoom(roomID int) []*Client {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	var clients []*Client
	for client := range s.clients {
		if client.currentRoom != nil && client.currentRoom.ID == roomID {
			clients = append(clients, client)
		}
	}
	
	return clients
}