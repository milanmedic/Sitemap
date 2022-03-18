package linkformatter

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	. "sitemap.mmedic.com/m/v2/src/models/link"
)

func AddDomainToLinks(links map[string]*Link, pageDomain string) map[string]*Link {
	for _, l := range links {
		//var l *Link = links[k]

		linkUrl, err := url.Parse(l.GetHref())
		if err != nil {
			log.Fatal(err)
		}

		var domain string = linkUrl.Hostname()
		if len(domain) == 0 {

			pageUrl, err := url.Parse(pageDomain)
			if err != nil {
				panic(err)
			}

			parts := strings.Split(pageUrl.Hostname(), ".")
			domain = strings.Join(parts, ".")
			l.SetHref(fmt.Sprintf("https://%s/%s", domain, l.GetHref()))
		}
	}
	return links
}
