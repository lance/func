package program

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func newCreate() Subcommand {
	items := []list.Item{
		MenuItem{title: "Language", desc: "Select a language for your function"},
		MenuItem{title: "Template", desc: "Choose a template for your function"},
	}
	help := NewHelp(`
		NAME
			{{.Name}} create - Create a Function project.
		
		SYNOPSIS
			{{.Name}} create [-l|--language] [-t|--template] [-r|--repository]
									[-c|--confirm]  [-v|--verbose]  [path]
		
		DESCRIPTION
			Creates a new Function project.
		
				$ {{.Name}} create -l node -t http
		
			Creates a Function in the current directory '.' which is written in the
			language/runtime 'node' and handles HTTP events.
		
			If [path] is provided, the Function is initialized at that path, creating
			the path if necessary.
		
			To complete this command interactivly, use --confirm (-c):
				$ {{.Name}} create -c
		
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
	return Subcommand{
		menu:        menu,
		help:        help,
		displayHelp: false,
	}
}

func (c Subcommand) Init() tea.Cmd {
	var commands = []tea.Cmd{}
	for _, i := range c.menu.Items() {
		m := i.(MenuItem).model
		if isModel(m) {
			commands = append(commands, m.Init())
		}
	}

	return tea.Batch(commands...)
}

func (c Subcommand) View() string {
	if c.displayHelp == true {
		return docStyle.Render(c.help.View())
	}
	return docStyle.Render(c.menu.View())
}

func (c Subcommand) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "m":
			c.displayHelp = true
		}
	}

	var cmd tea.Cmd
	c.menu, cmd = c.menu.Update(msg)
	return c, cmd
}

// CreateCmd is a tea.Cmd which is executed asynchronously and returns its result
func CreateCmd(name, root, runtime, template string) tea.Msg {
	// TODO: create a client with the provided arguments and run client.Create()
	return nil
}
