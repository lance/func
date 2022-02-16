package program

import (
	tea "github.com/charmbracelet/bubbletea"
)

type create struct {
	name, root, runtime, template string
}

func (m create) Init() tea.Cmd {
	return nil
}

func (m create) View() string {
	return "This is the screen to prompt for create options"
}

func (m create) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// CreateCmd is a tea.Cmd which is executed asynchronously and returns its result
func CreateCmd(name, root, runtime, template string) tea.Msg {
	// TODO: create a client with the provided arguments and run client.Create()
	return nil
}
