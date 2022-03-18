package htmlparser

import (
	"bytes"
	"strings"

	. "sitemap.mmedic.com/m/v2/src/models/link"
	. "sitemap.mmedic.com/m/v2/src/utils/file_reader"

	"golang.org/x/net/html"
)

type HTMLParser struct {
	fileReader *FileReader
}

func CreateHTMLParser() *HTMLParser {
	return &HTMLParser{fileReader: CreateFileReader()}
}

func (hp *HTMLParser) GetLinks(contents string) (map[string]*Link, error) {

	htmlReader := strings.NewReader(string(contents))
	tokenizer := html.NewTokenizer(htmlReader)

	var links map[string]*Link = make(map[string]*Link)
	linkFound := false
	var link *Link

	for {
		tokenizerToken := tokenizer.Next()

		switch tokenizerToken {
		case html.ErrorToken:
			return nil, tokenizer.Err()
		case html.TextToken:
			handleTextToken(links, link, &linkFound, tokenizer)
		case html.StartTagToken:
			tagName, _ := tokenizer.TagName()
			if string(tagName) == "a" {
				link = handleStartTagToken(links, &linkFound, tokenizer)
			}
		case html.EndTagToken:
			tagName, _ := tokenizer.TagName()
			if string(tagName) == "html" {
				return links, nil
			}
		}
	}
}

func handleTextToken(links map[string]*Link, link *Link, linkFound *bool, tokenizer *html.Tokenizer) {
	if *linkFound {
		text := tokenizer.Text()
		if !bytes.Equal(text, []byte(html.CommentToken.String())) {
			link.SetText(string(text))
			links[link.GetHref()] = link
			*linkFound = false
		}
	}
}

func handleStartTagToken(links map[string]*Link, linkFound *bool, tokenizer *html.Tokenizer) *Link {
	_, attrValue, _ := tokenizer.TagAttr()
	if strings.Index(string(attrValue), "#") != 0 && strings.Index(string(attrValue), "mailto") != 0 {
		link := CreateEmptyLink()
		*linkFound = true

		link.SetHref(string(attrValue))
		return link
	}
	return nil
}
