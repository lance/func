package program

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define some styles for our views to use
var (
	docStyle      = lipgloss.NewStyle().Margin(1, 2)
	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// NewProgram creates and starts a new TUI program, providing the user with a
// fully interactive function developer experience. The initial screen is a
// list of all commands, with a small description. Hitting <ENTER> on any of
// the commands takes the user to the command screen.
func NewProgram() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	return p.Start()
}

// The top-level model for all func commands, consisting of the command menu,
// and a pointer to the currently active command
type model struct {
	menu       list.Model // name and descriptive text for each func command
	active     bool       // a hack that can be removed when all menu items have a model
	quitting   bool       // a flag to indicate when the application is shutting down
	subcommand subcommand // model for the currently active subcommand
}

// subcommand is implemented by all commands, e.g. create, build, deploy
type subcommand struct {
	menu        list.Model
	help        tea.Model
	displayHelp bool
}

// initialModel returns the initial state of the program - a list of all possible commands
// written in user friendly language, with nothing selected.
func initialModel() (m model) {
	items := []list.Item{
		menuItem{title: "Create", desc: "Create a new function project from a template", model: newCreate()},
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

	m = model{
		menu:       list.New(items, list.NewDefaultDelegate(), 0, 0),
		active:     false,
		quitting:   false,
		subcommand: subcommand{},
	}
	m.menu.Title = "⚡ Knative Functions ⚡"
	return
}

// Init initializes the state of the program, providing an opportunity for
// sub-models to take care of any startup/initialization tasks
func (m model) Init() tea.Cmd {
	// Initialize sub-models
	var cmds = []tea.Cmd{}
	for _, i := range m.menu.Items() {
		m := i.(menuItem).Model()
		if m != nil {
			cmds = append(cmds, m.Init())
		}
	}
	return tea.Batch(cmds...)
}

// Update allows for changes in the program state based on user behavior
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.String() == "q" {
			// Typing CTRL-C or "q" exits the program
			m.quitting = true
			m.active = false
			return m, tea.Quit
		} else if msg.Type == tea.KeyEnter {
			// Enter key activates the current selection
			m.active = true
			i, _ := m.menu.SelectedItem().(menuItem)
			m.subcommand = i.model
		} else if msg.Type == tea.KeyEsc {
			// Escape key deactivates and short circuits by returning
			m.active = false
			return m, nil
		}

	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.menu.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	}

	// Only update the menu if we are not currently in a subcommand
	if m.active {
		// Get the currently selected model and update it as well
		sm := m.selectedModel()
		if sm != nil {
			var mm tea.Model
			mm, cmd = m.selectedModel().Update(msg)
			i, ok := m.menu.SelectedItem().(menuItem)
			if ok {
				i.model = mm
			}
			m.menu.SetItem(m.menu.Index(), mm)
			cmds = append(cmds, cmd)
		}
	} else {
		m.menu, cmd = m.menu.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View builds and returns a string based on the state of the program model
// If there is a current command, it delegates the view to the sub-model
func (m model) View() string {
	if m.quitting {
		return "\n👋 Bye!"
	}
	if m.active {
		if m.selectedModel() != nil {
			// A command has been selected, render the sub-model's view
			return docStyle.Render(m.selectedModel().View())
		}
		// There is a selected command, but it doesn't have the program TUI implemented yet
		return docStyle.Render(fmt.Sprintf("Sorry, %s isn't available yet.", m.selectedCommand()))
	}
	return docStyle.Render(m.menu.View())
}

// selectedModel returns the model for the currently selected menu item
// or nil if nothing is selected
func (m model) selectedModel() tea.Model {
	i, ok := m.menu.SelectedItem().(menuItem)
	if ok {
		return i.Model()
	}
	return nil
}

// selectedCommand returns the command for the currently selected menu
// or empty string if nothing is selected
func (m model) selectedCommand() string {
	i, ok := m.menu.SelectedItem().(menuItem)
	if ok {
		return i.Title()
	}
	return ""
}

type menuItem struct {
	title, desc string
	model       subcommand
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }
func (i menuItem) Model() subcommand   { return i.model }
