package program

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type create struct {
	menu        list.Model
	help        tea.Model
	displayHelp bool
	ready       bool
}

func newCreate() tea.Model {
	items := []list.Item{
		MenuItem{title: "Language", desc: "Select a language for your function"},
		MenuItem{title: "Template", desc: "Choose a template for your function"},
	}
	help := NewHelp(`
# NAME
	func create - Create a Function project.

# SYNOPSIS
	func create [-l|--language] [-t|--template] [-r|--repository]
							[-c|--confirm]  [-v|--verbose]  [path]

# DESCRIPTION
	Creates a new Function project.

		$ func create -l node -t http

	Creates a Function in the current directory '.' which is written in the
	language/runtime 'node' and handles HTTP events.

	If [path] is provided, the Function is initialized at that path, creating
	the path if necessary.

	To complete this command interactivly, use --confirm (-c):
		$ func create -c

	Available Language Runtimes and Templates:
{{ .Options | indent 2 " " | indent 1 "\t" }}

	To install more language runtimes and their templates see '{{.Name}} repository'.

EXAMPLES
	o Create a Node.js Function (the default language runtime) in the current
		directory (the default path) which handles http events (the default
		template).
		$ {{.Name}} create

	o Create a Node.js Function in the directory 'myfunc'.
		$ {{.Name}} create myfunc

	o Create a Go Function which handles CloudEvents in ./myfunc.
		$ {{.Name}} create -l go -t cloudevents myfunc
		`)
	menu := list.New(items, list.NewDefaultDelegate(), 0, 0)
	menu.Title = "✦ Create a new function project ✦"
	return create{
		menu:        menu,
		help:        help,
		ready:       false,
		displayHelp: false,
	}
}

func (c create) Init() tea.Cmd {
	var commands = []tea.Cmd{}
	commands = append(commands, c.help.Init())
	return tea.Batch(commands...)
}

func (c create) View() string {
	if c.displayHelp == true {
		return c.help.View()
	}
	return c.menu.View()
}

func (c create) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands = []tea.Cmd{}
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "H":
			c.displayHelp = true
		}

	case tea.WindowSizeMsg:
		if !c.ready {
			c.ready = true
		}
		c.menu.SetSize(windowSize(msg))
	}
	var cmd tea.Cmd
	// Update the menu
	c.menu, cmd = c.menu.Update(msg)
	commands = append(commands, cmd)

	// Update the help
	c.help, cmd = c.help.Update(msg)
	commands = append(commands, cmd)
	return c, cmd
}

// CreateCmd is a tea.Cmd which is executed asynchronously and returns its result
func CreateCmd(name, root, runtime, template string) tea.Msg {
	// TODO: create a client with the provided arguments and run client.Create()
	return nil
}
