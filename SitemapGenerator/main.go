package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/QuickOrBeDead/GoLangLearning/datastructures"

	"golang.org/x/net/html"
)

func main() {
	rootUrl, err := url.Parse("https://www.calhoun.io")
	if err != nil {
		fmt.Println("invalid url")
		return
	}

	rootNode, err := parseHtml(rootUrl.String())
	if err != nil {
		fmt.Println("error parsing html: ", err.Error())
		return
	}

	sitemap := &datastructures.SitemapNode{Url: "/", Children: make([]*datastructures.SitemapNode, 0), Parent: nil, Level: 0}

	var searchLinks func(*html.Node, *datastructures.SitemapNode)
	searchLinks = func(n *html.Node, s *datastructures.SitemapNode) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := attr.Val
					currentUrl, err := url.Parse(href)

					if err == nil {
						if currentUrl.Host == "" || (currentUrl.Host == rootUrl.Host && currentUrl.Scheme == rootUrl.Scheme) {
							if currentUrl.Path != "" {
								c := s.AddChild(currentUrl.Path, s)
								if c == nil || c.Level > 1 {
									break
								}

								time.Sleep(1 * time.Second)
								node, err := parseHtml(rootUrl.JoinPath(c.Url).String())
								if err == nil {
									searchLinks(node, c)
								} else {
									fmt.Println("error parsing html: ", err.Error())
								}
							}
						} else {
							break
						}
					}

					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			searchLinks(c, s)
		}
	}

	searchLinks(rootNode, sitemap)

	printSitemap(sitemap)
}

func printSitemap(s *datastructures.SitemapNode) {
	fmt.Printf("%s%s\n", strings.Repeat("   ", s.Level), s.Url)

	for _, c := range s.Children {
		printSitemap(c)
	}
}

func parseHtml(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("page get error. status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rootNode, err := html.Parse(bytes.NewReader(respBody))
	if err != nil {
		return nil, err
	}

	return rootNode, nil
}
