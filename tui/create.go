package tui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CreateModel struct {
    textInput   textinput.Model
    err         error
    ProjectName string
    done        bool
    quit        bool
}

func NewCreateModel() *CreateModel {
    ti := textinput.New()
    ti.Placeholder = "Your project name"
    ti.Focus()
    ti.CharLimit = 150
    ti.Width = 50

    return &CreateModel{
        textInput: ti,
        err:       nil,
    }
}

func (m *CreateModel) Init() tea.Cmd {
    return textinput.Blink
}

func (m *CreateModel) Run() (tea.Model, error) {
    p := tea.NewProgram(m)
    model, err := p.Run()
    if err != nil {
        return nil, err
    }

    return model, nil
}

func (m *CreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEnter:
            if m.textInput.Value() == "" {
                return m, nil
            }
            m.ProjectName = m.textInput.Value()
            err := createProject(m.ProjectName)
            if err != nil {
                m.err = err
                return m, tea.Quit
            }
            m.done = true
            return m, tea.Quit

        case tea.KeyCtrlC, tea.KeyEsc:
            m.quit = true
            return m, tea.Quit
        }
    }

    m.textInput, cmd = m.textInput.Update(msg)
    return m, cmd
}

func (m CreateModel) View() string {
    if m.err != nil {
        return fmt.Sprintf("\nError: %v\n", m.err)
    }
    if m.done {
        return fmt.Sprintf("\nProject %s created successfully!\n", m.ProjectName)
    }
    if m.quit{
        return fmt.Sprintln("\nProject exited without entering name")
    }

    return fmt.Sprintf(
        "%s\n\n%s\n\n%s",
        titleStyle.Render("What's the name of your new Go project?"),
        inputStyle.Render(m.textInput.View()),
        "Press Enter to create project, or Esc to cancel",
    )
}

func createProject(projectName string) error {
    if err := os.MkdirAll(projectName, 0755); err != nil {
        return fmt.Errorf("failed to create project directory: %v", err)
    }

    if err := os.Chdir(projectName); err != nil {
        return fmt.Errorf("failed to change to project directory: %v", err)
    }

    dirs := []string{
        "cmd",
        "internal",
        "pkg",
        "api",
        "configs",
        "test",
    }

    for _, dir := range dirs {
        if err := os.MkdirAll(dir, 0755); err != nil {
            return fmt.Errorf("failed to create directory %s: %v", dir, err)
        }
    }

    cmd := exec.Command("go", "mod", "init", projectName)
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to initialize go module: %v", err)
    }

    return nil
}