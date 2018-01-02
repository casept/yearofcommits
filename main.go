package main

import (
	"log"

	"github.com/getlantern/systray"
	"github.com/plutov/yearofcommits/icon"
)

func main() {
	// Should be called at the very beginning of main().
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("100")
	mQuit := systray.AddMenuItem("Quit", "Quit program")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
		log.Println("Quit!")
	}()
}

func onExit() {
	// clean up here
}
