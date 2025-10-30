package main

import (
	"io"
	"net"
	"log"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		mustCopy(os.Stdout, conn)
		done <- struct{}{}
	}()
	mustCopy(conn, os.Stdin)

	<-done
}

func mustCopy(w io.Writer, r io.Reader) {
	_, err := io.Copy(w, r)
	if err != nil {
		log.Fatal(err)
	}
}