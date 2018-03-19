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

	settingsDialog *SettingsDialog
}

func NewGraphics(aurora *Aurora) *Graphics {
	self := &Graphics{}

	self.settingsDialog = NewSettingsDialog()

	self.aurora = aurora

	self.Frame = wx.NewFrame(wx.NullWindow, -1, "Aurora RAT", wx.DefaultPosition, wx.NewSizeT(800, 600))

	self.statusBar = self.CreateStatusBar()
	self.statusBar.SetStatusText("Not Currently Listening.")

	self.menuBar = wx.NewMenuBar()
	menuFile := wx.NewMenu()

	self.menuItemListen = wx.NewMenuItem(menuFile, wx.ID_ANY, "&Start Listening", "Start Listening", wx.ITEM_NORMAL)
	menuFile.Append(self.menuItemListen)

	menuItemSettings := wx.NewMenuItem(menuFile, wx.ID_ANY, "&Settings", "Settings", wx.ITEM_NORMAL)
	menuFile.Append(menuItemSettings)

	wx.Bind(self, wx.EVT_MENU, func(e wx.Event) {
		self.Close(true)
	}, wx.ID_EXIT)
	menuFile.Append(wx.ID_EXIT)

	self.menuBar.Append(menuFile, "&File")

	self.SetMenuBar(self.menuBar)

	wx.Bind(self, wx.EVT_MENU, self.toggleListening, self.menuItemListen.GetId())
	wx.Bind(self, wx.EVT_MENU, self.settings, menuItemSettings.GetId())

	wx.Bind(self.settingsDialog, wx.EVT_BUTTON, self.save, self.settingsDialog.saveBtn.GetId())
	wx.Bind(self.settingsDialog, wx.EVT_BUTTON, self.load, self.settingsDialog.loadBtn.GetId())

	self.Layout()

	self.AddChild(self.settingsDialog)

	return self
}

func (graphics *Graphics) settings(_ wx.Event) {
	graphics.settingsDialog.ShowModal()
}

func (graphics *Graphics) save(_ wx.Event) {
	graphics.aurora.config.setHost(graphics.settingsDialog.hostCtrl.GetValue())
	graphics.aurora.config.setPort(graphics.settingsDialog.portCtrl.GetValue())
	graphics.aurora.saveConfig()
	graphics.settingsDialog.Close()
	graphics.settingsDialog.hostCtrl.SetValue("")
	graphics.settingsDialog.portCtrl.SetValue("")
}

func (graphics *Graphics) load(_ wx.Event) {
	graphics.aurora.loadConfig()
	graphics.settingsDialog.hostCtrl.SetValue(graphics.aurora.config.getHost())
	graphics.settingsDialog.portCtrl.SetValue(graphics.aurora.config.getPort())
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

type SettingsDialog struct {
	wx.Dialog
	hostCtrl wx.TextEntry
	portCtrl wx.TextEntry
	saveBtn  wx.Button
	loadBtn  wx.Button
}

func NewSettingsDialog() *SettingsDialog {
	self := &SettingsDialog{}
	self.Dialog = wx.NewDialog(wx.NullWindow, -1, "Settings")

	columnOne := wx.NewBoxSizer(wx.VERTICAL)
	rowOneTwo := wx.NewBoxSizer(wx.VERTICAL)
	row := wx.NewBoxSizer(wx.HORIZONTAL)

	row.Add(columnOne)
	row.Add(rowOneTwo)

	self.hostCtrl = wx.NewTextCtrl(self, wx.ID_ANY, "", wx.DefaultPosition, wx.DefaultSize, 0)
	columnOne.Add(self.hostCtrl, 0, wx.ALL|wx.EXPAND, 5)

	self.portCtrl = wx.NewTextCtrl(self, wx.ID_ANY, "", wx.DefaultPosition, wx.DefaultSize, 0)
	rowOneTwo.Add(self.portCtrl, 0, wx.ALL|wx.EXPAND, 5)

	self.saveBtn = wx.NewButton(self, wx.ID_ANY, "Save", wx.DefaultPosition, wx.DefaultSize, 0)
	columnOne.Add(self.saveBtn, 0, wx.ALL|wx.EXPAND, 5)

	self.loadBtn = wx.NewButton(self, wx.ID_ANY, "Load", wx.DefaultPosition, wx.DefaultSize, 0)
	rowOneTwo.Add(self.loadBtn, 0, wx.ALL|wx.EXPAND, 5)

	self.SetSizer(row)

	self.Layout()

	return self
}
