package main

import (
	"embed"
	"os/exec"
	"runtime"
	"strings"

	"fyne.io/systray"
	"github.com/gocolly/colly/v2"
)

//go:embed icons
var icon embed.FS

type Link struct {
	Name string
	URL  string
}

var links []Link

func main() {
	links, _ = getTitles()
	links, _ = getTitles()

	systray.Run(start, nil)
}

func getTitles() ([]Link, error) {
	var url = "https://eksisozluk.com"
	var links []Link
	c := colly.NewCollector()

	//FAKE HEADERS
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		r.Headers.Set("Accept-Charset", "UTF-8,*;q=0.5")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.8")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36")
	})

	/*c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response Size:",len(r.Body))
		fmt.Println("Response Received:", string(r.Body))
	})*/

	c.OnHTML("div#index-section ul.topic-list li", func(e *colly.HTMLElement) {
		if len(links) >= 10 || e.Attr("id") != "" {
			return
		}
		link := Link{
			Name: strings.TrimSpace(e.Text),
			URL:  url + strings.TrimSpace(e.ChildAttr("a", "href")),
		}
		links = append(links, link)
	})

	if err := c.Visit(url); err != nil {
		return nil, err
	}
	return links, nil
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
	links, _ = getTitles()
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
