package pagegetter

import (
	"io/ioutil"
	"net/http"
)

type PageGetter struct{}

func CreatePageGetter() *PageGetter {
	return &PageGetter{}
}

func (pg *PageGetter) GetPage(pageUrl string) ([]byte, error) {

	response, err := http.Get(pageUrl)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	html, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return html, nil
}

func (pg *PageGetter) GetPageAsString(pageUrl string) (string, error) {
	contents, err := pg.GetPage(pageUrl)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}
