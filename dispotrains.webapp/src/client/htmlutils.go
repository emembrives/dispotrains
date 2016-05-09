package client

import (
	"golang.org/x/net/html"
)

func findNode(n *html.Node, tagName string) *html.Node {
	return findNodeWithAttributes(n, tagName, make(map[string]string))
}

func findNodeWithAttributes(n *html.Node, tagName string, attributes map[string]string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tagName {
		var ok bool = true
		for key, value := range attributes {
			attr := findAttrByKey(n, key)
			if attr == nil || attr.Val != value {
				ok = false
				break
			}
		}
		if ok == true {
			return n
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		var result *html.Node = findNodeWithAttributes(c, tagName, attributes)
		if result != nil {
			return result
		}
	}
	return nil
}

func findNext(n *html.Node) *html.Node {
	for c := n.NextSibling; c != nil; c = c.NextSibling {
		if c.Type == n.Type && c.Data == n.Data {
			return c
		}
	}
	return nil
}

func findAttrByKey(n *html.Node, key string) *html.Attribute {
	for _, attribute := range n.Attr {
		if attribute.Key == key {
			return &attribute
		}
	}
	return nil
}
