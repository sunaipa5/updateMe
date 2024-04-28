package main

import (
	"embed"
	"fyne.io/systray"
	"github.com/gocolly/colly/v2"
	"os/exec"
	"runtime"
	"strings"
)

//go:embed icon.png
var icon embed.FS

type Link struct {
	Name string
	URL  string
}

var links []Link

func main() {
	links, _ = getTitles()
	systray.Run(start, nil)
}

func getTitles() ([]Link, error) {
	var url = "https://eksisozluk.com"
	var links []Link
	c := colly.NewCollector()
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
	trayIcon, _ := icon.ReadFile("icon.png")
	systray.SetIcon(trayIcon)
	runMenu()
}

func restart() {
	links, _ = getTitles()
	systray.ResetMenu()
	runMenu()
}

func runMenu() {
	systray.AddMenuItem("UpdateMe","").Disable()
	systray.AddSeparator()

	for _, link := range links {
		menuItem := systray.AddMenuItem(link.Name, link.URL)
		go func(url string) {
			for {
				select {
				case <-menuItem.ClickedCh:
					openBrowser(url)
				}
			}
		}(link.URL)
	}
	systray.AddSeparator()
	refresh := systray.AddMenuItem("Refresh", "Refresh")
	go func() {
		for {
			select {
			case <-refresh.ClickedCh:
				restart()

			}
		}
	}()
	quit := systray.AddMenuItem("Quit", "Quit")
	<-quit.ClickedCh
	systray.Quit()
}
