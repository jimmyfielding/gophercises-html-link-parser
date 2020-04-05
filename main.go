package main

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

var exampleHtml = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">
    A link to another page
    <span> some span  </span>
  </a>
  <a href="/page-two">A link to a second page</a>
</body>
</html>
`

func main() {
	r := strings.NewReader(exampleHtml)
	links, err := parse(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", links)
}

type Link struct {
	Href string
	Text string
}

func parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	nodes := getLinkNodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, constructLink(node))
	}

	return links, nil
}

func constructLink(n *html.Node) Link {
	var link Link
	for _, attr := range n.Parent.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
		}
	}

	link.Text = getLinkText(n)
	return link
}

func getLinkText(n *html.Node) string {
	if n.Type != html.ElementNode {
		return ""
	}

	if n.Type == html.TextNode {
		return n.Data
	}

	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getLinkText(c)
	}

	return strings.Join(strings.Fields(text), " ")
}

func getLinkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var nodes []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, getLinkNodes(c)...)
	}

	return nodes
}
