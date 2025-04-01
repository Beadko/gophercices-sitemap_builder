package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Beadko/gophercises_link/link"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

var (
	urlFlag = flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
	depth   = flag.Int("depth", 10, "the maximum number of links deep to traverse")
)

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	flag.Parse()

	if *urlFlag == "" {
		fmt.Println("Please provide a URL using -url flag")
		return
	}

	pages := bfs(*urlFlag, *depth)
	toXML := urlset{
		Urls:  make([]loc, len(pages)),
		Xmlns: xmlns,
	}
	for i, page := range pages {
		toXML.Urls[i] = loc{page}
	}
	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXML); err != nil {
		fmt.Printf("Could not encode data into XML file: %v\n", err)
	}
	fmt.Println()
}

func bfs(urlString string, depth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlString: {},
	}
	for i := 0; i <= depth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for u := range q {
			if _, ok := seen[u]; ok {
				continue
			}
			seen[u] = struct{}{}
			for _, l := range get(u) {
				if _, ok := seen[l]; !ok {
					nq[l] = struct{}{}
				}
			}
		}
	}
	r := make([]string, 0, len(seen))
	for u := range seen {
		r = append(r, u)
	}
	return r
}

func get(u string) []string {
	resp, err := http.Get(u)
	if err != nil {
		fmt.Printf("Could not get the response from %s: %v\n", *urlFlag, err)
		return []string{}
	}
	defer resp.Body.Close()
	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()

	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, err := link.Parse(r)
	if err != nil {
		fmt.Printf("Could not parse the links from the url %s body: %v\n", *urlFlag, err)
	}
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var r []string
	for _, link := range links {
		if keepFn(link) {
			r = append(r, link)
		}
	}
	return r
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
