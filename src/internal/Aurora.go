package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var err error

// Aurora Struct and Functions

type Aurora struct {
	settingsDir  string
	settingsFile string

	listener  net.Listener
	listening bool

	liveConnections ItemSet

	config *Config
}

func NewAurora() *Aurora {
	self := &Aurora{}
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}

	self.settingsDir, _ = filepath.Abs(usr.HomeDir + "\\Aurora")
	self.settingsFile, _ = filepath.Abs(self.settingsDir + "\\settings.json")

	self.loadConfig()

	self.setListening(false)

	return self
}

func (aurora *Aurora) getListening() bool {
	return aurora.listening
}

func (aurora *Aurora) setListening(listening bool) {
	aurora.listening = listening
}

// Loads the currently saved configuration. This is normally only called on startup
// and should not need to be called again because the struct and file are always synced.

func (aurora *Aurora) loadConfig() {
	_, err := os.Stat(aurora.settingsFile)
	if os.IsNotExist(err) {
		aurora.defaultConfig()
	} else if err != nil {
		fmt.Println(err, 1)
	} else {
		configFile, err := ioutil.ReadFile(aurora.settingsFile)
		if err != nil {
			fmt.Println(err, 2)
		}
		err = json.Unmarshal(configFile, &aurora.config)
		if err != nil {
			fmt.Println(err, 3)
		}
		if aurora.config == nil {
			aurora.config = newConfig("127.0.0.1", "4731", "UXpQV01GT2hXb1JBMHZ1QU0vdG5BMlZrZElRMHV2YmxTNXVoKytDUVBRZz0=")
			aurora.saveConfig()
		}
	}
}

// Saves the current configuration.

func (aurora *Aurora) saveConfig() {
	settings, err := json.Marshal(aurora.config)
	if err != nil {
		fmt.Println(err, 4)
	}
	err = ioutil.WriteFile(aurora.settingsFile, settings, 0644)
	if err != nil {
		fmt.Println(err, 5)
	}
}

// Creates a default config and saves it.

func (aurora *Aurora) defaultConfig() {
	err := os.MkdirAll(aurora.settingsDir, os.ModeDir)
	if err != nil {
		fmt.Println(err, 6)
	}
	_, err = os.Create(aurora.settingsFile)
	if err != nil {
		fmt.Println(err, 7)
	}
	aurora.config = newConfig("127.0.0.1", "4731", "UXpQV01GT2hXb1JBMHZ1QU0vdG5BMlZrZElRMHV2YmxTNXVoKytDUVBRZz0=")
	settings, err := json.Marshal(aurora.config)
	if err != nil {
		fmt.Println(err, 8)
	}
	err = ioutil.WriteFile(aurora.settingsFile, settings, 0644)
	if err != nil {
		fmt.Println(err, 9)
	}
}

// Establishes connection with a client.

func (aurora *Aurora) startListening() {
	aurora.setListening(true)
	aurora.listener, err = net.Listen("tcp", ":"+aurora.config.getPort())
	if err != nil {
		fmt.Println(err, 10)
	}
	for {
		if !aurora.getListening() {
			runtime.Goexit()
		}
		conn, err := aurora.listener.Accept()
		if err != nil {
			fmt.Println(err, 11)
			aurora.setListening(false)
			aurora.stopListening()
			runtime.Goexit()
		}
		aurora.liveConnections.Add(conn)
	}
}

// Disconnects from all clients.

func (aurora *Aurora) stopListening() {
	aurora.setListening(false)
	aurora.listener.Close()
}

// Send commands to the clients.

func (aurora *Aurora) sendCommand(conn net.Conn, cmd string) {
	writer := bufio.NewWriter(conn)
	switch cmd {
	case "PING\\":
		_, err := writer.WriteString(cmd)
		if err != nil {
			aurora.liveConnections.Delete(conn)
		}
		go aurora.receive(conn, cmd)
	}
}

// Amount of currently live connections.

func (aurora *Aurora) getLiveConnections() int {
	for _, client := range aurora.liveConnections.Items() {
		aurora.sendCommand(client, "PING\\")
	}
	return aurora.liveConnections.Size()
}

// Receives responses from the clients based on each command.

func (aurora *Aurora) receive(conn net.Conn, cmd string) {
	reader := bufio.NewReader(conn)
	if !aurora.getListening() {
		aurora.liveConnections.Delete(conn)
		runtime.Goexit()
	}
	message, err := reader.ReadString('\\')

	if err != nil {
		fmt.Println(err, 12)
		aurora.liveConnections.Delete(conn)
		runtime.Goexit()
	}

	switch cmd {
	case "PING\\":
		if message != "PONG\\" {
			aurora.liveConnections.Delete(conn)
		}
		runtime.Goexit()
	case "SCREENSHOT\\":
		runtime.Goexit()
	}
}

// Stub Config Struct and Functions

type Config struct {
	Host string
	Port string
	Key  [32]byte
}

func newConfig(newHost string, newPort string, newKey string) *Config {
	self := &Config{}
	self.Host = newHost
	self.Port = newPort
	var key [32]byte
	copy(key[:], newKey)
	self.Key = key
	return self
}

func (config *Config) setHost(newHost string) {
	config.Host = newHost
}

func (config *Config) getHost() string {
	return config.Host
}

func (config *Config) setPort(newPort string) {
	config.Port = newPort
}

func (config *Config) getPort() string {
	return config.Port
}

func (config *Config) setKey(newKey string) {
	var key [32]byte
	copy(key[:], newKey)
	config.Key = key
}

func (config *Config) getKey() [32]byte {
	return config.Key
}
