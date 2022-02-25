package program

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"

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
	viewport   viewport.Model
	menu       list.Model // name and descriptive text for each func command
	active     bool       // a hack that can be removed when all menu items have a model
	quitting   bool       // a flag to indicate when the application is shutting down
	ready      bool
	subcommand Subcommand // model for the currently active subcommand
}

type MenuItem struct {
	title, desc string
	model       Subcommand
}

func (i MenuItem) Title() string       { return i.title }
func (i MenuItem) Description() string { return i.desc }
func (i MenuItem) FilterValue() string { return i.title }

// Subcommand is implemented by all commands, e.g. create, build, deploy
type Subcommand struct {
	menu        list.Model
	help        tea.Model
	displayHelp bool
}

// initialModel returns the initial state of the program - a list of all possible commands
// written in user friendly language, with nothing selected.
func initialModel() (m model) {
	items := []list.Item{
		MenuItem{title: "Create", desc: "Create a new function project from a template", model: newCreate()},
		MenuItem{title: "Build", desc: "Turn an existing function project into a runnable container", model: Subcommand{}},
		MenuItem{title: "Configure", desc: "View and update options for an existing function project", model: Subcommand{}},
		MenuItem{title: "Deploy", desc: "Run an existing function project on a cluster", model: Subcommand{}},
		MenuItem{title: "Undeploy", desc: "Remove an existing function from a cluster", model: Subcommand{}},
		MenuItem{title: "Info", desc: "See information about an existing function", model: Subcommand{}},
		MenuItem{title: "List", desc: "Get a list of all functions deployed on the cluster", model: Subcommand{}},
		MenuItem{title: "Run", desc: "Run an existing function in a local container", model: Subcommand{}},
		MenuItem{title: "Invoke", desc: "Invoke a running function, either locally or on a cluster", model: Subcommand{}},
		MenuItem{title: "Templates", desc: "Install and update reusable function templates", model: Subcommand{}},
	}

	m = model{
		menu:     list.New(items, list.NewDefaultDelegate(), 0, 0),
		active:   false,
		quitting: false,
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
		m := i.(MenuItem).model
		if isModel(m) {
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
		} else if msg.String() == "b" {
			// "b" deactivates the current command and returns
			// us to the menu
			m.active = false
		} else if msg.Type == tea.KeyEnter {
			// Enter key activates the currently selected command
			m.active = true
			i, _ := m.menu.SelectedItem().(MenuItem)
			m.subcommand = i.model
		}

	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		w := msg.Width - left - right
		h := msg.Height - top - bottom

		if !m.ready {
			m.viewport = viewport.Model{Width: w, Height: h}
			m.ready = true
		}

		m.menu.SetSize(w, h)
		m.viewport.Width = w
		m.viewport.Height = h
		m.viewport.SetContent(m.menu.View())
	}

	if m.active && isModel(m.subcommand) {
		// A subcommand is active, send the received
		// message to it and update its model
		var mm tea.Model
		mm, cmd = m.subcommand.Update(msg)
		m.subcommand = mm.(Subcommand)
		cmds = append(cmds, cmd)
	} else {
		// No subcommand is active, just update the menu
		m.menu, cmd = m.menu.Update(msg)
		cmds = append(cmds, cmd)
	}
	// Finally, update the viewport
	m.viewport, cmd = m.viewport.Update(msg)
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
		if isModel(m.subcommand) {
			// A command has been selected, render the sub-model's view
			m.viewport.SetContent(m.subcommand.View())
		} else {
			// There is a selected command, but it doesn't have the program TUI implemented yet
			m.viewport.SetContent(fmt.Sprintf("Sorry, %s isn't available yet.", m.selectedCommand()))
		}
	} else {
		// No subcommand is selected, just show the menu
		m.viewport.SetContent(m.menu.View())
	}
	return docStyle.Render(m.viewport.View())
}

// isModel is a convenience function to determine if an interface
// implements the bubbletea Model interface. Currently only used
// to determine if a subcommand has been implemented yet.
func isModel(i interface{}) bool {
	_, ok := i.(tea.Model)
	fmt.Printf("\n\n\nModel %v, %+v\n\n\n", ok, i)
	return ok
}

// selectedCommand returns the command for the currently selected menu
// or empty string if nothing is selected
func (m model) selectedCommand() string {
	i, ok := m.menu.SelectedItem().(MenuItem)
	if ok {
		return i.Title()
	}
	return ""
}
