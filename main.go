package main

import (
	"flag"
	"fmt"
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

	resp, err := http.Get(*urlFlag)
	if err != nil {
		fmt.Printf("Could not get the response from %s: %v\n", *urlFlag, err)
		return
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()

	links, err := link.Parse(resp.Body)
	if err != nil {
		fmt.Printf("Could not parse the links from the url %s body: %v\n", *urlFlag, err)
	}
	var hrefs []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}
	for _, href := range hrefs {
		fmt.Println(href)
	}
}

/*
2. parse all the links on the page
3 build proper urls with our links
4. filter out any links w/ a diff domain
5. find all the pages (BFS)
6. print out XML
*/
