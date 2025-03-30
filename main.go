package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Beadko/gophercises_link/link"
)

var (
	urlFlag = flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
)

func main() {
	flag.Parse()

	if *urlFlag == "" {
		fmt.Println("Please provide a URL using -url flag")
		return
	}

	pages := get(*&urlFlag)
	for _, page := range pages {
		fmt.Println(page)
	}
}

func get(u *string) []string {
	resp, err := http.Get(*u)
	if err != nil {
		fmt.Printf("Could not get the response from %s: %v\n", *urlFlag, err)
		return nil
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()

	return hrefs(resp.Body, base)
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

/*
3 build proper urls with our links
4. filter out any links w/ a diff domain
5. find all the pages (BFS)
6. print out XML
*/
