package util

import (
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

const (
	PKG_SEARCH_URL = "https://pkg.go.dev/search?m=package&%s"
)

func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	var result strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r := getText(c)
		if r != "" {
			result.WriteString(strings.TrimSpace(r))
		}
	}

	return result.String()
}

func getLastChild(n *html.Node) *html.Node {
	var lastChild *html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lastChild = c
	}
	return lastChild
}

func getPrevSibling(n *html.Node) *html.Node {
	return n.PrevSibling
}

func parseResultsHtml(root *html.Node) []string {
	var res []string

	stack := []*html.Node{root}

	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "data-gtmv" {
					res = append(res, getText(n))
				}
			}
		}

		for c := getLastChild(n); c != nil; c = getPrevSibling(c) {
			stack = append(stack, c)
		}
	}

	return res
}

func GetPkgUrl(value string) string {
	re := regexp.MustCompile(`\((.*?)\)`)
	match := re.FindStringSubmatch(value)
	if len(match) > 1 {
		return match[1]
	}
	return match[0]
}

func Search(term string) []string {
	params := url.Values{}
	params.Add("q", term)
	searchUrl := fmt.Sprintf(PKG_SEARCH_URL, params.Encode())
	fmt.Printf("URL: %s\n", searchUrl)
	resp, err := http.Get(searchUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	res := parseResultsHtml(doc)

	return res
}

/*
DEPRECATED: Do not run this function since GoPack does not manage global installations.
*/
func RunGoInstall(pkg string) error {
	fmt.Printf("Running: go install %s\n", pkg)
	cmd := exec.Command("go", "install", pkg)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// RunGoGet is the same as RunGoInstall but it uses 'go get' instead of 'go install'.
func RunGoGet(pkg string) error {
	cmd := exec.Command("go", "get", pkg)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
