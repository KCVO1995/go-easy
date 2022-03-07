package main

import (
	"github.com/KCVO1995/go-easy/server"
	"github.com/webview/webview"
)

func main() {
	go server.Run()
	openWebview()
}

func openWebview () {
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("GoEasy")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://127.0.0.1:"+server.Port+"/static")
	w.Run()
}
