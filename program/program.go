package program

import (
	"github.com/charmbracelet/bubbles/key"
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
	viewport   viewport.Model // the primary view window
	menu       list.Model     // name and descriptive text for each func command
	subcommand tea.Model      // model for the currently active subcommand
	active     bool           // a flag to indicate that a subcommand is active
	quitting   bool           // a flag to indicate when the application is shutting down
	ready      bool           // signals initialization of view components is complete
}

type MenuItem struct {
	title, desc string
	model       tea.Model
}

func (i MenuItem) Title() string       { return i.title }
func (i MenuItem) Description() string { return i.desc }
func (i MenuItem) FilterValue() string { return i.title }

// subcommand is a simple placeholder model until all commands have
// been implemented
type subcommand struct {
	name string
}

func (s subcommand) Init() tea.Cmd                       { return nil }
func (s subcommand) Update(tea.Msg) (tea.Model, tea.Cmd) { return s, nil }
func (s subcommand) View() string                        { return s.name + ": Not implemented" }

func newSubcommand(name string) tea.Model {
	return subcommand{name: name}
}

// initialModel returns the initial state of the program - a list of all possible commands
// written in user friendly language, with nothing selected.
func initialModel() (m model) {
	items := []list.Item{
		MenuItem{
			title: "Create",
			desc:  "Create a new function project from a template",
			model: newCreate()},
		MenuItem{
			title: "Build",
			desc:  "Turn an existing function project into a runnable container",
			model: newSubcommand("Build")},
		MenuItem{
			title: "Configure",
			desc:  "View and update options for an existing function project",
			model: newSubcommand("Configure")},
		MenuItem{
			title: "Deploy",
			desc:  "Run an existing function project on a cluster",
			model: newSubcommand("Deploy")},
		MenuItem{
			title: "Undeploy",
			desc:  "Remove an existing function from a cluster",
			model: newSubcommand("Undeploy")},
		MenuItem{
			title: "Info", desc: "See information about an existing function",
			model: newSubcommand("Info")},
		MenuItem{
			title: "List",
			desc:  "Get a list of all functions deployed on the cluster",
			model: newSubcommand("List")},
		MenuItem{
			title: "Run",
			desc:  "Run an existing function in a local container",
			model: newSubcommand("Run")},
		MenuItem{
			title: "Invoke",
			desc:  "Invoke a running function, either locally or on a cluster",
			model: newSubcommand("Invoke")},
		MenuItem{
			title: "Templates",
			desc:  "Install and update reusable function templates",
			model: newSubcommand("Templates")},
	}

	m = model{
		menu:     list.New(items, list.NewDefaultDelegate(), 0, 0),
		active:   false,
		quitting: false,
	}
	m.menu.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "back"),
		)}
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
		cmds = append(cmds, m.Init())
	}
	return tea.Batch(cmds...)
}

// Update handles changes in the program state
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
		w, h := windowSize(msg)

		if !m.ready {
			m.viewport = viewport.Model{Width: w, Height: h}
			m.ready = true
		}

		m.menu.SetSize(w, h)
		m.viewport.Width = w
		m.viewport.Height = h
	}

	if m.active {
		// A subcommand is active, send the received
		// message to it and update its model
		var mm tea.Model
		mm, cmd = m.subcommand.Update(msg)
		m.subcommand = mm
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

func windowSize(msg tea.WindowSizeMsg) (w, h int) {
	top, right, bottom, left := docStyle.GetMargin()
	w = msg.Width - left - right
	h = msg.Height - top - bottom
	return
}

// View builds and returns a string based on the state of the program model
// If there is a current command, it delegates the view to the sub-model
func (m model) View() string {
	if m.quitting {
		return "\n👋 Bye!"
	}
	if m.active {
		// A command has been selected, render the sub-model's view
		m.viewport.SetContent(m.subcommand.View())
	} else {
		// No subcommand is selected, just show the menu
		m.viewport.SetContent(m.menu.View())
	}
	return docStyle.Render(m.viewport.View())
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
