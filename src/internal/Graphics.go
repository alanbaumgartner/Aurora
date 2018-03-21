package internal

import (
	"fmt"
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

	wx.Bind(self.settingsDialog, wx.EVT_BUTTON, self.saveConfig, self.settingsDialog.saveBtn.GetId())
	wx.Bind(self.settingsDialog, wx.EVT_BUTTON, self.loadConfig, self.settingsDialog.loadBtn.GetId())
	wx.Bind(self.settingsDialog, wx.EVT_BUTTON, self.buildClient, self.settingsDialog.buildBtn.GetId())
	wx.Bind(self.settingsDialog, wx.EVT_BUTTON, self.generateKey, self.settingsDialog.generateKeyBtn.GetId())

	self.Layout()

	self.AddChild(self.settingsDialog)

	return self
}

func (graphics *Graphics) settings(_ wx.Event) {
	graphics.settingsDialog.ShowModal()
}

func (graphics *Graphics) saveConfig(_ wx.Event) {
	graphics.aurora.config.setHost(graphics.settingsDialog.hostCtrl.GetValue())
	graphics.aurora.config.setPort(graphics.settingsDialog.portCtrl.GetValue())
	graphics.aurora.config.setKey(graphics.settingsDialog.keyCtrl.GetValue())
	graphics.aurora.saveConfig()
	graphics.settingsDialog.Close()
	graphics.settingsDialog.hostCtrl.SetValue("")
	graphics.settingsDialog.portCtrl.SetValue("")
	graphics.settingsDialog.keyCtrl.SetValue("")
}

func (graphics *Graphics) loadConfig(_ wx.Event) {
	graphics.aurora.loadConfig()
	graphics.settingsDialog.hostCtrl.SetValue(graphics.aurora.config.getHost())
	graphics.settingsDialog.portCtrl.SetValue(graphics.aurora.config.getPort())
	keyString := string(graphics.aurora.config.Key[:])
	graphics.settingsDialog.keyCtrl.SetValue(keyString)
}

func (graphics *Graphics) buildClient(_ wx.Event) {
	fmt.Println(graphics.aurora.getLiveConnections())
}

func (graphics *Graphics) generateKey(_ wx.Event) {
	key := NewEncryptionKey()
	keyString := string(key[:])
	graphics.settingsDialog.keyCtrl.SetValue(keyString)
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
	keyCtrl  wx.TextEntry

	buildBtn       wx.Button
	saveBtn        wx.Button
	loadBtn        wx.Button
	generateKeyBtn wx.Button
}

func NewSettingsDialog() *SettingsDialog {
	self := &SettingsDialog{}
	self.Dialog = wx.NewDialog(wx.NullWindow, -1, "Settings")

	columnOne := wx.NewBoxSizer(wx.VERTICAL)
	columnTwo := wx.NewBoxSizer(wx.VERTICAL)
	columnThree := wx.NewBoxSizer(wx.VERTICAL)
	row := wx.NewBoxSizer(wx.HORIZONTAL)

	row.Add(columnOne)
	row.Add(columnTwo)
	row.Add(columnThree)

	// Labels

	hostLbl := wx.NewStaticText(self, wx.ID_ANY, "Host", wx.DefaultPosition, wx.DefaultSize, 0)
	columnOne.Add(hostLbl, 0, wx.CENTER, 5)

	portLbl := wx.NewStaticText(self, wx.ID_ANY, "Port", wx.DefaultPosition, wx.DefaultSize, 0)
	columnTwo.Add(portLbl, 0, wx.CENTER, 5)

	passLbl := wx.NewStaticText(self, wx.ID_ANY, "Password", wx.DefaultPosition, wx.DefaultSize, 0)
	columnThree.Add(passLbl, 0, wx.CENTER, 5)

	// Text Entries

	self.hostCtrl = wx.NewTextCtrl(self, wx.ID_ANY, "", wx.DefaultPosition, wx.DefaultSize, 0)
	columnOne.Add(self.hostCtrl, 0, wx.EXPAND|wx.ALL, 5)

	self.portCtrl = wx.NewTextCtrl(self, wx.ID_ANY, "", wx.DefaultPosition, wx.DefaultSize, 0)
	columnTwo.Add(self.portCtrl, 0, wx.EXPAND|wx.ALL, 5)

	self.keyCtrl = wx.NewTextCtrl(self, wx.ID_ANY, "", wx.DefaultPosition, wx.DefaultSize, 0)
	columnThree.Add(self.keyCtrl, 0, wx.EXPAND|wx.ALL, 5)

	// Buttons

	self.generateKeyBtn = wx.NewButton(self, wx.ID_ANY, "Generate Key", wx.DefaultPosition, wx.DefaultSize, 0)
	columnThree.Add(self.generateKeyBtn, 0, wx.EXPAND|wx.ALL, 5)

	self.buildBtn = wx.NewButton(self, wx.ID_ANY, "Build", wx.DefaultPosition, wx.DefaultSize, 0)
	columnTwo.Add(self.buildBtn, 0, wx.EXPAND|wx.ALL, 5)

	self.saveBtn = wx.NewButton(self, wx.ID_ANY, "Save", wx.DefaultPosition, wx.DefaultSize, 0)
	columnTwo.Add(self.saveBtn, 0, wx.EXPAND|wx.ALL, 5)

	self.loadBtn = wx.NewButton(self, wx.ID_ANY, "Load", wx.DefaultPosition, wx.DefaultSize, 0)
	columnThree.Add(self.loadBtn, 0, wx.EXPAND|wx.ALL, 5)

	self.SetSizer(row)

	self.Layout()

	return self
}
