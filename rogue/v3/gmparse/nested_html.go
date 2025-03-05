package gmparse

import (
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

type ImageInfo struct {
	Src     string
	Width   string
	Height  string
	Alt     string
	Caption string
}

// Updated regex to capture the units
var styleRegex = regexp.MustCompile(`(width|height):\s*(\d+(?:\.\d+)?(?:px|%|em|rem|vh|vw))`)

func parseStyle(style string) (width, height string) {
	matches := styleRegex.FindAllStringSubmatch(style, -1)
	for _, match := range matches {
		if len(match) == 3 {
			if match[1] == "width" {
				width = match[2]
			} else if match[1] == "height" {
				height = match[2]
			}
		}
	}
	return
}

func splitHtml(htmlContent string) (ImageInfo, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return ImageInfo{}, err
	}

	var info ImageInfo

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "img" {
				for _, attr := range n.Attr {
					switch attr.Key {
					case "src":
						info.Src = attr.Val
					case "width":
						// Add "px" if it's just a number
						if matched, _ := regexp.MatchString(`^\d+$`, attr.Val); matched {
							info.Width = attr.Val + "px"
						} else {
							info.Width = attr.Val
						}
					case "height":
						// Add "px" if it's just a number
						if matched, _ := regexp.MatchString(`^\d+$`, attr.Val); matched {
							info.Height = attr.Val + "px"
						} else {
							info.Height = attr.Val
						}
					case "alt":
						info.Alt = attr.Val
					case "style":
						w, h := parseStyle(attr.Val)
						if info.Width == "" && w != "" {
							info.Width = w
						}
						if info.Height == "" && h != "" {
							info.Height = h
						}
					}
				}
			} else if n.Data == "figcaption" {
				if n.FirstChild != nil {
					info.Caption = n.FirstChild.Data
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	return info, nil
}
