package pagesdb

import . "sitemap.mmedic.com/m/v2/src/models/page"

type PageDB struct {
	pages []*Page
}

func CreatePageDb() *PageDB {
	return &PageDB{}
}

func (p *PageDB) PushPage(page *Page) {
	p.pages = append(p.pages, page)
}

func (p *PageDB) PopPage() *Page {
	var page *Page
	if !p.IsEmpty() {
		page = p.pages[len(p.pages)-1]
		p.pages = p.pages[:len(p.pages)-1]
		return page
	}
	return nil
}

func (p *PageDB) IsEmpty() bool {
	if len(p.pages) > 0 {
		return false
	}
	return true
}
