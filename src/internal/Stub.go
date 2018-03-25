package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:4731")
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}

	buffer := make([]byte, 1024)

	file, _ := os.Open("f.exe")
	defer file.Close()

	packet := Packet{}
	enc := json.NewEncoder(conn)

	i := 0

	for {
		packet.Type = "FILE"
		packet.FileName = "test.exe"
		packet.Done = false
		packet.BytePos = int64(i)
		_, err := file.Read(buffer)
		if err == io.EOF {
			packet.BytePos = 0
			packet.Done = true
			enc.Encode(packet)
			break
		}
		if i == 10 {
			enc.Encode(Packet{"MESSAGE", "AA", 0, []byte("TEST"), false})
		}
		packet.Data = buffer
		enc.Encode(packet)
		i++
		fmt.Println(i)
	}
}
