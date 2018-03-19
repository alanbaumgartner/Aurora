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
	"strings"
)

var err error

// Aurora Struct and Functions

type Aurora struct {
	settingsDir  string
	settingsFile string

	listener  net.Listener
	listening bool

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
			aurora.config = newConfig("127.0.0.1", "4731")
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
	aurora.config = newConfig("127.0.0.1", "4731")
	settings, err := json.Marshal(aurora.config)
	if err != nil {
		fmt.Println(err, 8)
	}
	err = ioutil.WriteFile(aurora.settingsFile, settings, 0644)
	if err != nil {
		fmt.Println(err, 9)
	}
}

func (aurora *Aurora) stopListening() {
	aurora.setListening(false)
	aurora.listener.Close()
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
		go aurora.receive(conn)
	}
}

// Listens for incoming messages from connections established in startListening().

func (aurora *Aurora) receive(conn net.Conn) {
	for {
		if !aurora.getListening() {
			runtime.Goexit()
		}
		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		message, err := rw.ReadString('\\')
		message = strings.TrimRight(message, "\\")

		if err != nil {
			fmt.Println(err, 12)
			runtime.Goexit()
		}
		rw.Flush()
		fmt.Println(message)
	}
}

// Stub Config Struct and Functions

type Config struct {
	Host string
	Port string
}

func newConfig(newHost string, newPort string) *Config {
	self := &Config{}
	self.Host = newHost
	self.Port = newPort
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
