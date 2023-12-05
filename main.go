package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/net/html"
)

type model struct {
	textInput textinput.Model
	choices   []string
	cursor    int
	startIdx  int
	err       error
}

const (
	pkgSearchUrl = "https://pkg.go.dev/search?m=package&%s"
)

var (
    results []string
	selection string
	step      int = 0
)

func initModel() model {
	ti := textinput.New()
	ti.Placeholder = "Package name"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		choices:   results,
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) FirstStepUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
            pkgName := m.textInput.Value()
            search(pkgName)
            m.choices = results
            step++
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) SecondStepUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
    if len(m.choices) == 0 {
        return m, tea.Quit
    }

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			selection = getPkgUrl(m.choices[m.cursor])
			return m, tea.Quit
		}
	}

	const pageSize = 5
	m.startIdx = max(0, min(len(m.choices)-pageSize, m.cursor-(pageSize/2)))

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if step == 0 {
		return m.FirstStepUpdate(msg)
	}
	return m.SecondStepUpdate(msg)
}

func (m model) View() string {
	// The header
	var s string
	if step == 0 {
		s = fmt.Sprintf("Enter package name:\n\n%s\n\n%s\n", m.textInput.View(), "(esc to quit)")
	} else {
		s = "Choose package to install:\n\n"

		// Iterate over our choices
		for i, choice := range m.choices[m.startIdx:min(m.startIdx+5, len(m.choices))] {

			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if m.cursor == m.startIdx+i {
				cursor = ">" // cursor!
			}

			// Render the row
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		// Number of results
		s += fmt.Sprintf("\nNumber of results: %d\n", len(m.choices))

		// The footer
		s += "\nPress q to quit.\n"
	}

	// Send the UI for rendering
	return s
}

func getResults(n *html.Node) {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "data-gtmv" {
				// fmt.Printf("Result: %s\n", getText(n))
				results = append(results, getText(n))
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getResults(c)
	}
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

func getPkgUrl(value string) string {
	re := regexp.MustCompile(`\((.*?)\)`)
	match := re.FindStringSubmatch(value)
	if len(match) > 1 {
		return match[1]
	}
	return match[0]
}

func search(term string) {
	params := url.Values{}
	params.Add("q", term)
	searchUrl := fmt.Sprintf(pkgSearchUrl, params.Encode())
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
	getResults(doc)
}

func runGoGet(pkg string) error {
    fmt.Printf("Running: go get %s\n", pkg)
    cmd := exec.Command("go", "get", pkg)
    err := cmd.Run()
    if err != nil {
        return err
    }
    return nil
}

func main() {
	m := initModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

    fmt.Printf("Selected: %s\n", selection)
    err := runGoGet(selection)
    if err != nil {
        panic(err)
    }
}
