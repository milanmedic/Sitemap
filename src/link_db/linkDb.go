package linkdb

import (
	. "sitemap.mmedic.com/m/v2/src/models/link"
)

type LinkDB struct {
	links map[string]*Link
}

func CreateLinkDb() *LinkDB {
	return &LinkDB{}
}

func (ld *LinkDB) AddLink(link *Link) {
	ld.links[link.GetHref()] = link
}

func (ld *LinkDB) GetLink(url string) *Link {
	return ld.links[url]
}

func (ld *LinkDB) DeleteLink(url string) *Link {
	if link, ok := ld.links[url]; ok {
		delete(ld.links, url)
		return link
	}
	return nil
}
