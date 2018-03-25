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
	fmt.Println("Aurora: Now accepting connections.")
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := aurora.listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Aurora: New connection from", conn.RemoteAddr())
		go aurora.handlePackets(conn)
	}
}

func (aurora *Aurora) handlePackets(conn net.Conn) {
	files := map[string]*os.File{}
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
				if packet.Done && files[packet.FileName] != nil {
					files[packet.FileName].Close()
					fmt.Println("Aurora: Finished downloading", packet.FileName)
				} else if packet.Done && files[packet.FileName] == nil {
					continue
				} else {
					if files[packet.FileName] == nil {
						fmt.Println("Aurora: Started downloading", packet.FileName)
						if _, err := os.Stat(packet.FileName); os.IsNotExist(err) {
							files[packet.FileName], _ = os.Create(packet.FileName)
						} else {
							files[packet.FileName], _ = os.Open(packet.FileName)
						}
					}
					files[packet.FileName].WriteAt(packet.Data, packet.BytePos*1024)
				}
			case "MESSAGE":
				fmt.Println("Aurora: incoming message \"" + string(packet.Data) + "\"")
			}
		}
	}
}
