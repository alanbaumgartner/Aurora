package main

import (
	"Aurora/src/internal"
	"github.com/dontpanic92/wxGo/wx"
)

var backend *internal.Aurora
var frontend *internal.Graphics

func main() {
	backend = internal.NewAurora()
	frontend = internal.NewGraphics(backend)

	wx1 := wx.NewApp()
	frontend.Show()
	wx1.MainLoop()
}
