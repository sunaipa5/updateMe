package main

import (
	"embed"
	"os/exec"
	"runtime"

	"fyne.io/systray"
)

//go:embed icons
var icon embed.FS

type Link struct {
	Name string
	URL  string
}

var links []Link

func main() {
	links = request()

	systray.Run(start, nil)
}

func openBrowser(data string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", data).Start()
	case "windows":
		exec.Command("cmd", "/c", "start", data).Start()
	case "darwin":
		exec.Command("open", data).Start()
	default:
	}
}

func start() {
	//Select tray icon based on OS type
	var trayIcon []byte
	switch runtime.GOOS {
	case "linux", "darwin":
		trayIcon, _ = icon.ReadFile("icons/icon.png")
	case "windows":
		trayIcon, _ = icon.ReadFile("icons/icon.ico")
	}

	systray.SetIcon(trayIcon)
	systray.SetTitle("UpdateMe")
	systray.SetTooltip("UpdateMe")
	runMenu()
}

func restart() {
	links = request()
	systray.ResetMenu()
	runMenu()
}

func runMenu() {
	systray.AddMenuItem("UpdateMe", "").Disable()
	systray.AddSeparator()

	//Link items
	for _, link := range links {
		menuItem := systray.AddMenuItem(link.Name, link.URL)
		go func(url string) {
			<-menuItem.ClickedCh
			openBrowser(url)
		}(link.URL)
	}

	systray.AddSeparator()

	//Refresh Item
	refresh := systray.AddMenuItem("Refresh", "Refresh")
	go func() {
		<-refresh.ClickedCh
		restart()
	}()

	//Quit Item
	quit := systray.AddMenuItem("Quit", "Quit")
	<-quit.ClickedCh
	systray.Quit()
}
