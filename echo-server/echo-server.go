package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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

	input := bufio.NewScanner(conn)
	for input.Scan() {
		go echo(conn, input.Text(), 1 * time.Second)
	}
	if err := input.Err(); err != nil {
		log.Println("ошибка чтения:", err)
	}

}

func echo(conn net.Conn, str string, delay time.Duration) {
	fmt.Fprintln(conn, strings.ToUpper(str))
	time.Sleep(delay)
	fmt.Fprintln(conn, str)
	time.Sleep(delay)
	fmt.Fprintln(conn, strings.ToLower(str))
}