package page

type Page struct {
	url      string
	contents []byte
}

func CreatePage(url string, contents []byte) *Page {
	return &Page{url: url, contents: contents}
}

func (p *Page) GetURL() string {
	return p.url
}

func (p *Page) GetContents() []byte {
	return p.contents
}
