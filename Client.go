package main

import (
	"net"
	"os"
	"time"
	"bufio"
	"fmt"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:25565")
	if err != nil {
		fmt.Println(err)
	}
	exit := make(chan string)

	go send(conn)

	for {
		select {
		case <- exit: {
			os.Exit(0)
		}
		}
	}
}

func send(conn net.Conn) {
	for {
		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		_, err := rw.WriteString("Hello.\\")
		if err != nil {
			fmt.Println(err)
			conn.Close()
			os.Exit(0)
		}
		rw.Flush()
		fmt.Println("Sent")
		time.Sleep(2000000000)
	}
}