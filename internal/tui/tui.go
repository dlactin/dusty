package tui

import (
	"fmt"
	"os"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dlactin/dusty/internal/git"
)

var (
	appNameStyle = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1)
)

type model struct {
	branches []git.GitBranch
	cursor   int
	selected map[int]struct{}
}

func New(branches []git.GitBranch) {
	p := tea.NewProgram(initialModel(branches))
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}

func initialModel(branches []git.GitBranch) model {
	return model{
		// Our list of branches
		branches: branches,

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.branches)-1 {
				m.cursor++
			}

		// Delete the selected branch
		case "d", "del":
			branch := m.branches[m.cursor]
			git.DelBranch(branch.Name, false)
			m.branches = slices.Delete(m.branches, m.cursor, m.cursor+1)

			return m, nil
			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.

		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Local branches:\n\n"

	// Iterate over our choices
	for i, choice := range m.branches {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s -- Author: %s -- Merged: %t -- Age %d days\n", cursor, choice.Name, choice.Author, choice.Merged, choice.Age)
	}

	// The footer
	s += "\nPress d to delete branch, press q to quit.\n"

	// Send the UI for rendering
	return s
}
