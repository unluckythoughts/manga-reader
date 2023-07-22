package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/unluckythoughts/manga-reader/models"
)

type model struct {
	view    view
	wait    bool
	spinner spinner.Model
	favs    []models.MangaFavorite
	table   table.Model
}

func getInitialModel() model {
	m := model{
		view: favs,
		wait: true,
	}

	return m
}

func (m model) Init() tea.Cmd {
	if m.wait {
		m.spinner = spinner.New()
		m.spinner.Spinner = spinner.Jump
		return m.spinner.Tick
	}

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
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
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	switch m.view {
	case "favs":
		return favsView(m)
	}
	return ""
}

func main() {
	var err error
	m := getInitialModel()
	m.favs, err = getFavs()
	if err != nil {
		fmt.Println("could not get favs:", err)
		os.Exit(1)
	}

	m.table = populateFavsTable(m.favs)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
