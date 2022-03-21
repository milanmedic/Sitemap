package main

import (
	"encoding/xml"
	"flag"
	"io/ioutil"

	. "sitemap.mmedic.com/m/v2/src/models/link"
	"sitemap.mmedic.com/m/v2/src/models/page"
	. "sitemap.mmedic.com/m/v2/src/pages_db"
	. "sitemap.mmedic.com/m/v2/src/utils/html_parser"
	. "sitemap.mmedic.com/m/v2/src/utils/link_formatter"
	. "sitemap.mmedic.com/m/v2/src/utils/page_getter"
)

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Urls    []Url    `xml:"url"`
}

type Url struct {
	Location string `xml:"loc"`
}

func main() {

	startPage := flag.String("start_page", "https://gobyexample.com", "Starting URL")
	flag.Parse()

	pg := CreatePageGetter()
	contents, err := pg.GetPage(*startPage)
	if err != nil {
		panic(err)
	}

	pageDb := CreatePageDb()
	pageDb.PushPage(page.CreatePage(*startPage, contents))

	htmlParser := CreateHTMLParser()

	cnt := 1

	var xmlOut []byte = []byte(xml.Header)

	for !pageDb.IsEmpty() && cnt <= 2 {
		p := pageDb.PopPage()

		links, err := ExtractLinksFromPage(p, htmlParser)
		if err != nil {
			panic(err)
		}

		AddDomainToLinks(links, p.GetURL())

		xmlOut, err = ConvertURLsToXML(xmlOut, links)
		if err != nil {
			panic(err)
		}

		for k, l := range links {
			AddNewPages(pg, pageDb, l)
			delete(links, k)
		}

		cnt++
	}
	err = ioutil.WriteFile("sitemap.xml", []byte(xml.Header), 0644)
	err = ioutil.WriteFile("sitemap.xml", xmlOut, 0644)

}

func ExtractLinksFromPage(p *page.Page, htmlParser *HTMLParser) (map[string]*Link, error) {
	pageContents := string(p.GetContents())

	links, err := htmlParser.GetLinks(pageContents)
	if err != nil {
		return nil, err
	}

	return links, err
}

func ConvertURLsToXML(xmlOut []byte, links map[string]*Link) ([]byte, error) {
	for _, l := range links {
		u := &Url{Location: l.GetHref()}
		out, err := xml.MarshalIndent(&u, " ", "  ")
		xmlOut = append(xmlOut, out...)

		if err != nil {
			return nil, err
		}
	}
	return xmlOut, nil
}

func AddNewPages(pg *PageGetter, pageDb *PageDB, l *Link) {
	contents, err := pg.GetPage(l.GetHref())
	if err != nil {
		panic(err)
	}
	p := page.CreatePage(l.GetHref(), contents)
	pageDb.PushPage(p)
}
