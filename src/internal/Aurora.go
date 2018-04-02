package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Incoming socket data struct

type Packet struct {
	Type string

	StringData string

	BytePos  int64
	FileData []byte
	Done     bool
}

// CLI interface struct

type Aurora struct {
	listener net.Listener

	connections List
	decoders    map[net.Conn]*json.Decoder
	encoders    map[net.Conn]*json.Encoder

	workingDirectory  string
	downloadDirectory string
}

func (aurora *Aurora) Init() {
	var err error

	aurora.encoders = map[net.Conn]*json.Encoder{}
	aurora.decoders = map[net.Conn]*json.Decoder{}
	aurora.connections = NewList()

	aurora.workingDirectory, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	aurora.downloadDirectory, err = filepath.Abs(aurora.workingDirectory + "\\Downloads")
	if err != nil {
		fmt.Println(err)
	}
	go aurora.listen()
	for {
		buf := bufio.NewReader(os.Stdin)
		data, _, err := buf.ReadLine()
		if err != nil {
			fmt.Println(err)
		}

		msg := string(data)
		parts := strings.Split(msg, " ")
		index, _ := strconv.Atoi(parts[0])

		for len(parts) < 3 {
			parts = append(parts, "")
		}

		aurora.sendPackets(index, parts[1], parts[2])
	}

}

func (aurora *Aurora) sendPackets(index int, packetType string, msg string) {
	chosenConn := aurora.connections.Get(index)
	if chosenConn == nil {
		index = -1
	}
	switch packetType {
	case "file":
		aurora.uploadFile(aurora.connections.Get(0), msg)
	case "p":
		if index >= 0 {
			aurora.encoders[aurora.connections.Get(index)].Encode(Packet{"P", "", 0, nil, false})
		} else {
			for conn, enc := range aurora.encoders {
				err := enc.Encode(Packet{"P", "", 0, nil, false})
				if err != nil {
					aurora.removeConnection(conn)
					fmt.Println(err)
				}
			}
		}
	case "rp":
		if index >= 0 {
			aurora.encoders[aurora.connections.Get(index)].Encode(Packet{"RP", "", 0, nil, false})
		} else {
			for conn, enc := range aurora.encoders {
				err := enc.Encode(Packet{"RP", "", 0, nil, false})
				if err != nil {
					aurora.removeConnection(conn)
					fmt.Println(err)
				}
			}
		}
	case "msg":
		if index >= 0 {
			aurora.encoders[aurora.connections.Get(index)].Encode(Packet{"MSG", msg, 0, nil, false})
		} else {
			for conn, enc := range aurora.encoders {
				err := enc.Encode(Packet{"MSG", msg, 0, nil, false})
				if err != nil {
					aurora.removeConnection(conn)
					fmt.Println(err)
				}
			}
		}
	case "uninstall":
		if index >= 0 {
			aurora.encoders[aurora.connections.Get(index)].Encode(Packet{"UNINSTALL", "", 0, nil, false})
		} else {
			for conn, enc := range aurora.encoders {
				err := enc.Encode(Packet{"UNINSTALL", "", 0, nil, false})
				if err != nil {
					aurora.removeConnection(conn)
					fmt.Println(err)
				}
			}
		}
	case "live":
		for conn, enc := range aurora.encoders {
			err := enc.Encode(Packet{"PING", "", 0, nil, false})
			if err != nil {
				aurora.removeConnection(conn)
				fmt.Println(err)
			}
		}
		for index, conn := range aurora.connections.All() {
			fmt.Println(index, conn.RemoteAddr())
		}
	case "dc":
		for conn, enc := range aurora.encoders {
			err := enc.Encode(Packet{"DC", "", 0, nil, false})
			if err != nil {
				aurora.removeConnection(conn)
				fmt.Println(err)
			}
		}
	default:
	}
}

func (aurora *Aurora) listen() {
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
		aurora.addConnections(conn)
		go aurora.handlePackets(conn)
	}
}

func (aurora *Aurora) uploadFile(conn net.Conn, fileName string) {
	buffer := make([]byte, 1024)
	file, _ := os.Open(fileName)
	defer file.Close()

	i := 0
	for {
		_, err := file.Read(buffer)
		if err == io.EOF {
			err = aurora.encoders[conn].Encode(Packet{"FILE", fileName, 0, nil, true})
			if err != nil {
				aurora.removeConnection(conn)
				fmt.Println(err)
			}
			break
		}
		err = aurora.encoders[conn].Encode(Packet{"FILE", fileName, int64(i), buffer, false})
		if err != nil {
			aurora.removeConnection(conn)
			fmt.Println(err)
		}
		i++
	}
}

func (aurora *Aurora) handlePackets(conn net.Conn) {
	files := map[string]*os.File{}
	for {
		packet := Packet{}
		err := aurora.decoders[conn].Decode(&packet)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			if err != nil {
				aurora.removeConnection(conn)
				fmt.Println(err)
			}
			return
		} else {
			switch packet.Type {
			case "FILE":
				if _, err := os.Stat(aurora.downloadDirectory); os.IsNotExist(err) {
					err := os.MkdirAll(aurora.downloadDirectory, os.ModeDir)
					if err != nil {
						fmt.Println(err)
					}
				}
				fileName := aurora.downloadDirectory + "\\" + packet.StringData
				if packet.Done && files[fileName] != nil {
					files[fileName].Close()
					fmt.Println("Aurora: Finished downloading", packet.StringData)
					delete(files, fileName)
				} else if packet.Done && files[fileName] == nil {
					continue
				} else {
					if files[fileName] == nil {
						fmt.Println("Aurora: Started downloading", packet.StringData)
						if _, err := os.Stat(fileName); os.IsNotExist(err) {
							files[fileName], _ = os.Create(fileName)
						} else {
							files[fileName], _ = os.Open(fileName)
						}
						defer files[fileName].Close()
					}
					files[fileName].WriteAt(packet.FileData, packet.BytePos*1024)
				}
			case "MESSAGE":
				fmt.Println("Aurora: incoming message \"" + string(packet.FileData) + "\"")
			}
		}
	}
}

// Add or remove a connection.

func (aurora *Aurora) addConnections(conn net.Conn) {
	aurora.connections.Add(conn)
	aurora.encoders[conn] = json.NewEncoder(conn)
	aurora.decoders[conn] = json.NewDecoder(conn)
}

func (aurora *Aurora) removeConnection(conn net.Conn) {
	aurora.connections.Remove(conn)
	delete(aurora.encoders, conn)
	delete(aurora.decoders, conn)
}
