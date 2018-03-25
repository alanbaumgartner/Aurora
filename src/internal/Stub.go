package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:4731")
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}

	enc := json.NewEncoder(conn)

	wg.Add(2)

	go sendFile(enc, "test1.exe")
	go sendFile(enc, "test2.exe")

	wg.Wait()

}

func sendFile(enc *json.Encoder, fileName string) {
	buffer := make([]byte, 1024)

	file, _ := os.Open("f.exe")
	defer file.Close()

	packet := Packet{}
	i := 0
	for {
		packet.Type = "FILE"
		packet.FileName = fileName
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
	}
	wg.Done()
}
