package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := NewModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		if msg.Type == tea.KeyCtrlC {
// 			return m, tea.Quit
// 		}

// func (m model) View() string {
// 	switch m.State {
// 	case Loading:
// 		return "Loading...\n"
// 	case Ready:
// 		return issuesToString(m.Issues)
// 	case Error:
// 		return fmt.Sprintf("Error: %v\n", m.Error)
// 	default:
// 		return "Unknown state\n"
// 	}
// }
