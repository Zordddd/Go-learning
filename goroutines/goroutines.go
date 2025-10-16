package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		log.Println("connection opened")
		go handleconnection(conn)
	}
}

func handleconnection(conn net.Conn) {
	defer log.Println("connection closed")
	defer conn.Close()
	for {
		_, err := io.WriteString(conn, time.Now().Format("\r15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}