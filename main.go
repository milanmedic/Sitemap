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

	startPage := flag.String("start_page", "https://dave.cheney.net/", "Starting URL")
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

	var xmlOut []byte
	ioutil.WriteFile("sitemap.xml", []byte(xml.Header), 0)

	for !pageDb.IsEmpty() && cnt <= 3 {
		p := pageDb.PopPage()
		pageContents := string(p.GetContents())

		links, err := htmlParser.GetLinks(pageContents)
		if err != nil {
			panic(err)
		}

		// ADDING DOMAINS TO LINKS
		AddDomainToLinks(links, p.GetURL())

		// HERE WE SHOULD CREATE A SITEMAP (BUGGY!)

		for _, l := range links {
			u := &Url{Location: l.GetHref()}
			out, err := xml.MarshalIndent(&u, " ", "  ")
			xmlOut = append(xmlOut, out...)

			if err != nil {
				panic(err)
			}
		}

		ioutil.WriteFile("sitemap.xml", xmlOut, 0)

		for k, l := range links {
			AddNewPages(pg, pageDb, l)
			delete(links, k)
		}

		cnt++
	}
}

func AddNewPages(pg *PageGetter, pageDb *PageDB, l *Link) {
	contents, err := pg.GetPage(l.GetHref())
	if err != nil {
		panic(err)
	}
	p := page.CreatePage(l.GetHref(), contents)
	pageDb.PushPage(p)
}
