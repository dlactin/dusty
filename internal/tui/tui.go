package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dlactin/dusty/internal/git"
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
		branches: branches,
		// A map which indicates which choices are selected
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Check which key was pressed
		switch msg.String() {
		// Exit TUI
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
			if len(m.branches) == 0 {
				return m, nil
			}

			branch := m.branches[m.cursor]
			err := git.DelBranch(branch.Name, false)
			if err != nil {
				fmt.Printf("error deleting branches: %v", err)
			}

			// Remove the branch from the branches slice
			m.branches = append(m.branches[:m.cursor], m.branches[m.cursor+1:]...)

			// Adjust cursor position if needed
			if m.cursor >= len(m.branches) && len(m.branches) > 0 {
				m.cursor = len(m.branches) - 1
			}

			// If no branches left, quit
			if len(m.branches) == 0 {
				return m, tea.Quit
			}

			// Refresh TUI when a branch is deleted
			return m, tea.ClearScreen
		}
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	return m, nil
}

func (m model) View() string {
	// TUI Header
	s := "-- Matching local branches: --\n\n"
	// Iterate over our choices
	for i, choice := range m.branches {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		// Render the row, forcing column widths
		s += fmt.Sprintf("%-s %-30.30s | Author: %-13.13s | Merged: %-5t | Age: %-3d days\n", cursor, choice.Name, choice.Author, choice.Merged, choice.Age)
	}
	// TUI footer
	s += "\nPress d to delete branch, press q to quit.\n"
	// Return the UI for rendering
	return s
}
