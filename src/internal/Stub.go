package internal

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
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
		case <-exit:
			{
				os.Exit(0)
			}
		}
	}
}

func send(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	rw.WriteString("Hello.\\")
	rw.Flush()
	time.Sleep(2000000000)
	rw.WriteString("Hello.\\")
	rw.Flush()
	time.Sleep(2000000000)
	rw.WriteString("Hello.\\")
	rw.Flush()
	time.Sleep(2000000000)
	rw.WriteString("Hello.\\")
	rw.Flush()
	time.Sleep(2000000000)
	rw.WriteString("Done\\")
	rw.Flush()
}
