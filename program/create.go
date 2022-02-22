package program

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func newCreate() subcommand {
	items := []list.Item{
		menuItem{title: "Language", desc: "Select a language for your function"},
		menuItem{title: "Template", desc: "Choose a template for your function"},
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
	return subcommand{
		menu:        menu,
		help:        help,
		displayHelp: false,
	}
}

func (c subcommand) Init() tea.Cmd {
	var commands = []tea.Cmd{}
	for _, i := range c.menu.Items() {
		m := i.(menuItem).model
		if m != nil {
			commands = append(commands, m.Init())
		}
	}

	return tea.Batch(commands...)
}

func (c subcommand) View() string {
	fmt.Printf("Viewing %v %p\n", c.displayHelp, &c)
	if c.displayHelp == true {
		fmt.Println("HELP!")
		return docStyle.Render(c.help.View())
	}
	return docStyle.Render(c.menu.View())
}

func (c subcommand) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		fmt.Println("KeyMsg")
		switch keypress := msg.String(); keypress {
		case "m":
			c.displayHelp = true
		}

	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		c.menu.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	}

	var cmd tea.Cmd
	c.menu, cmd = c.menu.Update(msg)
	fmt.Printf("%+v %p\n", c.displayHelp, &c)
	return c, cmd
}

// CreateCmd is a tea.Cmd which is executed asynchronously and returns its result
func CreateCmd(name, root, runtime, template string) tea.Msg {
	// TODO: create a client with the provided arguments and run client.Create()
	return nil
}
