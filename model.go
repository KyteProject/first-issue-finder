package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v53/github"
)

// Issue represents a GitHub issue.
type Issue struct {
	Repo        string
	Stars       int
	Title       string
	Labels      []string
	IssueNumber int
}

// State represents the current state of the application.
type State int

const (
	Loading State = iota
	Idle
	Ready
	Error
)

// Model is the application's root data structure.
type Model struct {
	state  State
	issues []*github.Issue
	total  int

	table          table.Model
	paginator      paginator.Model
	loadingSpinner spinner.Model
	help           help.Model
	err            error
}

// NewModel returns a new model for the application.
func NewModel() Model {
	var state State = Loading
	var issues []*github.Issue

	loadingSpinner := spinner.New()
	loadingSpinner.Style = activeLabelStyle
	loadingSpinner.Spinner = spinner.Dot

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(issues))

	columns := []table.Column{
		{Title: "#", Width: 5},
		{Title: "Title", Width: 50},
		{Title: "Reactions", Width: 5},
		{Title: "Labels", Width: 50},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(25),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := Model{
		state:          state,
		table:          t,
		issues:         issues,
		help:           help.New(),
		loadingSpinner: loadingSpinner,
		paginator:      p,
	}

	return m
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return checkGitHubIssues
}

// Update updates the model in response to messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case errMsg:
		m.state = Error
		m.err = msg
		return m, nil

	case fetchIssuesMsg:
		var rows []table.Row

		m.issues = msg.Issues
		m.total = *msg.Total

		fmt.Printf("\nissue: %v \n", msg.Issues[0].GetRepository())

		for _, issue := range msg.Issues {

			row := table.Row{
				fmt.Sprint(issue.GetNumber()),
				issue.GetTitle(),
				fmt.Sprint(issue.Reactions.GetTotalCount()),
				safeDereferenceLabels(issue.Labels),
			}

			rows = append(rows, row)
		}

		m.table.SetRows(rows)
		m.state = Ready
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}

		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			id := m.table.SelectedRow()[1]

			return m, tea.Batch(
				tea.Printf("Let's go to %s!", id),
			)

		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View renders the application's UI.
func (m Model) View() string {
	if m.state == Loading {
		return m.loadingSpinner.View() + "Fetching issues..." + "\n"
	}

	if m.state == Error {
		return m.err.Error()
	}

	if m.state == Ready {
		return baseStyle.Render(m.table.View()) + fmt.Sprintf("\nTotal: %v \n", m.total) + " Press q to quit.\n"
	}

	return baseStyle.Render(m.table.View()) + "\nPress q to quit.\n"
}
