package internal

import (
	"bufio"
	"fmt"
	"net"
	"runtime"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:4731")
	if err != nil {
		fmt.Println(err)
	}
	receiveCommand(conn)
}

func sendMessage(conn net.Conn, cmd string) {
	writer := bufio.NewWriter(conn)
	switch cmd {
	case "PING":
		writer.WriteString("PONG\\")
	}
	writer.Flush()
}

func receiveCommand(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		cmd, err := reader.ReadString('\\')

		if err != nil {
			fmt.Println(err)
			runtime.Goexit()
		}

		cmd = strings.TrimRight(cmd, "\\")
		sendMessage(conn, cmd)
	}
}
