package main

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/unluckythoughts/manga-reader/models"
)

type view string

const (
	favs  view = "favs"
	manga view = "manga"
)

func populateFavsTable(favs []models.MangaFavorite) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Title", Width: 50},
		{Title: "Source", Width: 20},
		{Title: "Chapters", Width: 10},
	}

	rows := []table.Row{}
	for i, fav := range favs {
		rows = append(rows, table.Row{
			strconv.Itoa(i + 1), fav.Manga.Title, fav.Manga.Source.Name, strconv.Itoa(len(fav.Manga.Chapters)),
		})
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
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

	return t
}

func favsView(m model) string {
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	s := baseStyle.Render(m.table.View()) + "\n"

	return s
}
