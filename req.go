package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/RomainMichau/cloudscraper_go/cloudscraper"
)

func request() []Link {
	client, _ := cloudscraper.Init(false, false)
	res, err := client.Get("https://eksisozluk.com", make(map[string]string), "")
	if err != nil {
		panic(err)
	}

	return scraper(res.Body)
}

func scraper(respBody string) []Link {
	strings.NewReader(respBody)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(respBody))
	if err != nil {
		panic(err)
	}

	doc.Find("div#index-section ul.topic-list li").Each(func(i int, s *goquery.Selection) {
		if len(links) >= 10 || s.AttrOr("id", "") != "" {
			return
		}

		link := Link{
			Name: strings.TrimSpace(s.Text()),
			URL:  "https://eksisozluk.com" + strings.TrimSpace(s.Find("a").AttrOr("href", "")),
		}
		links = append(links, link)
	})

	return links
}
