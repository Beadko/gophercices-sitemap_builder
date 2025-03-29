package main

import (
	"flag"
	"fmt"
)

var (
	urlFlag = flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
)

func main() {
	flag.Parse()
	fmt.Println(*urlFlag)
}

/*
1. GET the webpage
2. parse all the links on the page
3 build proper urls with our links
4. filter out any links w/ a diff domain
5. find all the pages (BFS)
6. print out XML
*/
