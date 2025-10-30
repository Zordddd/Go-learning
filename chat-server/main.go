package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on " + listener.Addr().String())

	defer func() {
		err := listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go broadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(err, " handleConn")
		}
	}()
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "Вы " + who + " успешно подключились\n"
	messages <- who + " подключился\n"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	leaving <- ch
}

func clientWriter(conn net.Conn, ch chan string) {
	for msg := range ch {
		_, err := conn.Write([]byte(msg + "\r\n"))
		if err != nil {
			log.Fatal(err, " clientWriter ", msg)
		}
	}
}

func broadcast() {
	clients := make(map[client]bool)

	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}
