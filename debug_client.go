package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// Debug version of client for testing nc vs telnet
type DebugClient struct {
	conn   net.Conn
	reader *bufio.Reader
}

func NewDebugClient(conn net.Conn) *DebugClient {
	return &DebugClient{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}
}

func (c *DebugClient) Handle() {
	defer c.conn.Close()

	c.write("=== BBS Debug Mode ===\n")
	c.write("This will show exactly what input is being received.\n")
	c.write("Type 'quit' to exit.\n\n")

	for {
		c.write("debug> ")
		
		// Set timeout
		c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		
		// Read raw input
		line, err := c.reader.ReadString('\n')
		if err != nil {
			c.write(fmt.Sprintf("Error reading: %v\n", err))
			break
		}
		
		// Show what we received
		c.write(fmt.Sprintf("Raw input: %q (len=%d)\n", line, len(line)))
		
		// Clean and process
		clean := strings.TrimRight(line, "\r\n")
		c.write(fmt.Sprintf("Cleaned: %q (len=%d)\n", clean, len(clean)))
		
		if clean == "quit" {
			c.write("Goodbye!\n")
			break
		}
		
		c.write(fmt.Sprintf("You typed: %s\n\n", clean))
	}
}

func (c *DebugClient) write(message string) {
	c.conn.Write([]byte(message))
}

// Add debug mode to main server
func (s *BBSServer) StartDebug(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to start debug server: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Debug server started on port %s\n", port)
	fmt.Println("Connect with: nc localhost", port)
	fmt.Println("Or: telnet localhost", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go func() {
			client := NewDebugClient(conn)
			client.Handle()
		}()
	}
}