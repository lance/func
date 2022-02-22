package program

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

type help struct {
	viewport viewport.Model
	content  string
	ready    bool
}

func NewHelp(content string) tea.Model {
	return help{content: content}
}

func (h help) Init() tea.Cmd {
	return nil
}

func (h help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
		// 	return m, tea.Quit
		// }

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(h.headerView())
		footerHeight := lipgloss.Height(h.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !h.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			h.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			h.viewport.YPosition = headerHeight
			h.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			h.viewport.SetContent(h.content)
			h.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			h.viewport.YPosition = headerHeight + 1
		} else {
			h.viewport.Width = msg.Width
			h.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	// Handle keyboard and mouse events in the viewport
	h.viewport, cmd = h.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return h, tea.Batch(cmds...)
}

func (h help) headerView() string {
	title := titleStyle.Render("Mr. Pager")
	line := strings.Repeat("─", max(0, h.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (h help) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", h.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, h.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (h help) View() string {
	if !h.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", h.headerView(), h.viewport.View(), h.footerView())
}
