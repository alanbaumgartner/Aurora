package internal

import (
	"github.com/dontpanic92/wxGo/wx"
)

type Graphics struct {
	wx.Frame
	statusBar wx.StatusBar
	menuBar   wx.MenuBar

	menuItemListen wx.MenuItem

	aurora *Aurora
}

func NewGraphics(aurora *Aurora) *Graphics {
	self := &Graphics{}

	self.aurora = aurora

	self.Frame = wx.NewFrame(wx.NullWindow, -1, "Aurora RAT", wx.DefaultPosition, wx.NewSizeT(800, 600))

	self.statusBar = self.CreateStatusBar()
	self.statusBar.SetStatusText("Not Currently Listening.")

	self.menuBar = wx.NewMenuBar()
	menuFile := wx.NewMenu()
	menuSettings := wx.NewMenu()

	self.menuItemListen = wx.NewMenuItem(menuFile, wx.ID_ANY, "&Start Listening", "Start Listening", wx.ITEM_NORMAL)
	menuFile.Append(self.menuItemListen)

	wx.Bind(self, wx.EVT_MENU, func(e wx.Event) {
		self.Close(true)
	}, wx.ID_EXIT)
	menuFile.Append(wx.ID_EXIT)

	menuItemSave := wx.NewMenuItem(menuSettings, wx.ID_ANY, "&Save Settings", "Save Settings", wx.ITEM_NORMAL)
	menuSettings.Append(menuItemSave)

	self.menuBar.Append(menuFile, "&File")
	self.menuBar.Append(menuSettings, "&Settings")

	self.SetMenuBar(self.menuBar)

	wx.Bind(self, wx.EVT_MENU, self.toggleListening, self.menuItemListen.GetId())

	self.Layout()

	return self
}

func (graphics *Graphics) toggleListening(_ wx.Event) {
	if graphics.aurora.getListening() {
		graphics.statusBar.SetStatusText("Not Currently Listening.")
		graphics.menuItemListen.SetHelp("Start Listening")
		graphics.menuItemListen.SetItemLabel("Start Listening")
		go graphics.aurora.stopListening()
	} else {
		graphics.statusBar.SetStatusText("Currently Listening.")
		graphics.menuItemListen.SetHelp("Stop Listening")
		graphics.menuItemListen.SetItemLabel("Stop Listening")
		go graphics.aurora.startListening()
	}
}
