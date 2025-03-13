package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

const (
	PKG_SEARCH_URL = "https://pkg.go.dev/search?m=package&%s"
)

type Package struct {
	Path    string `json:"path"`
	Version string `json:"version"`
	Dir     string `json:"dir"`
}

func GetPkgUrl(value string) string {
	re := regexp.MustCompile(`\((.*?)\)`)
	match := re.FindStringSubmatch(value)
	if len(match) > 1 {
		return match[1]
	}
	return match[0]
}

// Search searches and parses the results from pkg.go.dev and returns the first 25 results.
func Search(term string) []string {
	params := url.Values{}
	params.Add("q", term)
	searchUrl := fmt.Sprintf(PKG_SEARCH_URL, params.Encode())
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

// RunGoGet is the same as RunGoInstall but it uses 'go get' instead of 'go install'.
func RunGoGet(pkg string) error {
	cmd := exec.Command("go", "get", pkg)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func ParseGoMod() ([]string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return nil, fmt.Errorf("Error opening go.mod file: %v\n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	inRequireBlock := false

	var modules []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// check if we're entering a require block
		if strings.HasPrefix(line, "require (") {
			inRequireBlock = true
			continue
		}

		// check if we're exiting a require block
		if line == ")" {
			inRequireBlock = false
			continue
		}

		// handle single-line require statement
		if strings.HasPrefix(line, "require ") {
			module := extractModuleFromLine(true, line)
			modules = append(modules, module)
			continue
		}

		// process modules within require block
		if inRequireBlock && line != "" && !strings.HasPrefix(line, "//") {
			module := extractModuleFromLine(false, line)
			modules = append(modules, module)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading go.mod file: %v\n", err)
	}

	return modules, nil
}

func extractModuleFromLine(singular bool, line string) string {
	var (
		module  string
		version string
	)
	fields := strings.Fields(line)

	if singular {
		module = fields[1]
		version = fields[2]
	} else {
		module = fields[0]
		version = fields[1]
	}

	return module + "@" + version
}

func GetDependencyList() ([]Package, error) {
	output, err := exec.Command("go", "list", "-m", "-json", "all").Output()
	if err != nil {
		return nil, fmt.Errorf("error executing command: %v", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(output))
	var packages []Package

	for decoder.More() {
		var pkg Package
		if err := decoder.Decode(&pkg); err != nil {
			return nil, fmt.Errorf("error parsing JSON: %v", err)
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

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
