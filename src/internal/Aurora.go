package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
)

// Incoming socket data struct

type Packet struct {
	Type     string
	FileName string
	BytePos  int64
	Data     []byte
	Done     bool
}

// CLI interface struct

type Aurora struct {
	listener net.Listener
}

func (aurora *Aurora) Listen() {
	var err error
	aurora.listener, err = net.Listen("tcp", ":4731")
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := aurora.listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go aurora.handlePackets(conn)
	}
}

func (aurora *Aurora) handlePackets(conn net.Conn) {
	var newFile *os.File
	decoder := json.NewDecoder(conn)
	for {
		packet := Packet{}
		err := decoder.Decode(&packet)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		} else {
			switch packet.Type {
			case "FILE":
				if packet.Done && newFile != nil {
					newFile.Close()
					fmt.Println("Finished downloading file named " + packet.FileName)
				} else if packet.Done && newFile == nil {
					continue
				} else {
					if newFile == nil {
						fmt.Println("Started downloading file named " + packet.FileName)
						if _, err := os.Stat(packet.FileName); os.IsNotExist(err) {
							newFile, _ = os.Create(packet.FileName)
						} else {
							newFile, _ = os.Open(packet.FileName)
						}
					}
					newFile.WriteAt(packet.Data, packet.BytePos*1024)
				}
			case "MESSAGE":
				fmt.Println("Message: " + string(packet.Data))
			}
		}
	}
}
