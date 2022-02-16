package program

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle      = lipgloss.NewStyle().Margin(1, 2)
	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

func NewProgram() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	return p.Start()
}

type model struct {
	menu     list.Model
	create   tea.Model
	command  string
	quitting bool
}

// initialModel returns the initial state of the program - a list of all possible commands
// written in user friendly language, with nothing selected.
func initialModel() model {
	items := []list.Item{
		menuItem{title: "Create", desc: "Create a new function project from a template", model: create{}},
		menuItem{title: "Build", desc: "Turn an existing function project into a runnable container"},
		menuItem{title: "Configure", desc: "View and update options for an existing function project"},
		menuItem{title: "Deploy", desc: "Run an existing function project on a cluster"},
		menuItem{title: "Undeploy", desc: "Remove an existing function from a cluster"},
		menuItem{title: "Info", desc: "See information about an existing function"},
		menuItem{title: "List", desc: "Get a list of all functions deployed on the cluster"},
		menuItem{title: "Run", desc: "Run an existing function in a local container"},
		menuItem{title: "Invoke", desc: "Invoke a running function, either locally or on a cluster"},
		menuItem{title: "Templates", desc: "Install and update reusable function templates"},
	}

	m := model{menu: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.menu.Title = "⚡ Knative Functions ⚡"
	m.quitting = false
	return m
}

func (m model) Init() tea.Cmd {
	// Initialize sub-models
	for _, i := range m.menu.Items() {
		i.Model().Init()
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
		case "q":
			m.quitting = true
			return m, tea.Quit

		case "esc":
			m.command = ""
			return m, nil

		case "enter":
			i, ok := m.menu.SelectedItem().(menuItem)
			if ok {
				m.command = string(i.title)
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.menu.SetSize(msg.Width-left-right, msg.Height-top-bottom)

	}

	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		// The user will probably not see this
		return quitTextStyle.Render("OK")
	}

	switch m.command {
	case "Create":
		return docStyle.Render(m.create.View())
	case "":
		return docStyle.Render(m.menu.View())
	default:
		return docStyle.Render(fmt.Sprintf("Sorry, %s isn't available yet.", m.command))
	}
}

type menuItem struct {
	title, desc string
	model       tea.Model
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }
func (i menuItem) Model() tea.Model    { return i.model }
