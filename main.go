package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
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
	io.Copy(os.Stdout, resp.Body)
}

/*
1. GET the webpage
2. parse all the links on the page
3 build proper urls with our links
4. filter out any links w/ a diff domain
5. find all the pages (BFS)
6. print out XML
*/
