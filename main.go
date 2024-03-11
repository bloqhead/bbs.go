package main

// learned:
// 1. goroutines and how they work
// 2. got a throwback lesson on ANSI
// 3. channels and how they can communicate across goroutines

import (
	"fmt"
	"os"
	"net"
	"strings"
	"time"
	"bufio"
	"context"
	"os/signal"
)

// utility function for handling errors
func check(err error, message string) {
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", message)
}

type ClientJob struct {
	cmd string
	conn net.Conn
}

type User struct {
	username string
	password string
}

func createAccount(username string, password string) {
}

func login(username string, password string) {
}

func logout() {
}

func generateResponses(clientJobs chan ClientJob) {
	for {
		clientJob := <-clientJobs

		// add a delay between jobs
		time.Sleep(time.Second)

		// if the user logs out, we give them a response
		if strings.Compare("logout", clientJob.cmd) == 0 {
			// say goodbye!
			clientJob.conn.Write([]byte("Bye!\n"))

			// and close the connection
			clientJob.conn.Close()
			break
		} else {
			// otherwise we pass the command along and output it for testing purposes
			fmt.Printf("%s\n", clientJob.cmd)
			clientJob.conn.Write([]byte(clientJob.cmd))
		}
	}
}

func main() {
	logo, err := os.ReadFile("logo.ans")

	// channels help communicate between goroutines
	// they can send and receive communications
	clientJobs := make(chan ClientJob)

	// this is a goroutine
	// a goroutine is a lightweight thread of execution
	// goroutines run concurrently
	go generateResponses(clientJobs)

	// start listening. this assumes nginx
	ln, err := net.Listen("tcp", ":3003")
	check(err, "Server ready!")

	// this is an anonymouse goroutine
	conn, err := ln.Accept()

	check(err, "Accepted connection!\n")

	go func() {
		buf := bufio.NewReader(conn)

		// we can pass our ANSI logo straight in here
		conn.Write([]byte(logo))

		// this is using old school ANSI escape sequences
		conn.Write([]byte("\033[H\033[2J"))
		conn.Write([]byte("\033[33m"))
		conn.Write([]byte("welcome to the BBS!\n"))

		for {
			conn.Write([]byte(">"))

			cmd, err := buf.ReadString('\n')
			cmd = strings.Replace(cmd, "\r\n", "", -1)

			if err != nil {
				fmt.Printf("Client disconnected.\n")
				break
			}

			// send our jobs
			clientJobs <- ClientJob{cmd, conn}
		}
	}()

	// close the connection when we're done
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
}
