package models

import "strings"

type Pattern struct {
	Match       string `json:"match,omitempty"`
	ReplaceWith string `json:"replace_with,omitempty"`
}

func (p *Pattern) IsRegex() bool {
	return strings.HasPrefix(p.Match, "/")
}

func (p *Pattern) GetRegexFlags() string {
	if !p.IsRegex() {
		return ""
	}

	// Find the last / to get flags
	lastSlashIdx := strings.LastIndex(p.Match, "/")
	if lastSlashIdx <= 0 || lastSlashIdx == len(p.Match)-1 {
		return ""
	}

	return p.Match[lastSlashIdx+1:]
}

func (p *Pattern) IsCaseSensitive() bool {
	flags := p.GetRegexFlags()
	return !strings.Contains(flags, "i")
}

type ContentSelectors struct {
	Data            string    `json:"data,omitempty"`
	ReplacePatterns []Pattern `json:"replace_patterns,omitempty"`
}

type ChapterSelectors struct {
	URL        string           `json:"url,omitempty"`
	Number     string           `json:"number,omitempty"`
	Title      string           `json:"title,omitempty"`
	UploadDate string           `json:"upload_date,omitempty"`
	DateFormat string           `json:"date_format,omitempty"`
	Content    ContentSelectors `json:"content,omitempty"`
}

type BookSelectors struct {
	URL             string           `json:"url,omitempty"`
	Title           string           `json:"title,omitempty"`
	CoverImage      string           `json:"cover_image,omitempty"`
	Synopsis        string           `json:"synopsis,omitempty"`
	ChapterListItem string           `json:"chapter_list_item,omitempty"`
	NextPage        string           `json:"next_page,omitempty"`
	Chapter         ChapterSelectors `json:"chapter,omitempty"`
}

type Selectors struct {
	BookListItem       string        `json:"book_list_item,omitempty"`
	NextPage           string        `json:"next_page,omitempty"`
	LastPage           string        `json:"last_page,omitempty"`
	NextPageURLPattern string        `json:"next_page_url_pattern,omitempty"`
	Book               BookSelectors `json:"book,omitempty"`
}

type Connector struct {
	Name        string    `json:"name,omitempty"`
	Domain      string    `json:"domain,omitempty"`
	IconURL     string    `json:"icon_url,omitempty"`
	BookListURL string    `json:"book_list_url,omitempty"`
	Type        BookType  `json:"type,omitempty"`
	Selectors   Selectors `json:"selectors,omitempty"`
}

type IConnector interface {
	GetName() string
	GetDomain() string
	GetIconURL() string
	GetSelectors() Selectors
	GetBooks() ([]Book, error)
	GetBookChapters(bookURL string) ([]Chapter, error)
	GetChapterContent(chapterURL string) (List, error)
}
