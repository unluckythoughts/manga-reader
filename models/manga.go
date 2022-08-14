package models

type Manga struct {
	URL      string
	Title    string
	ImageURL string
	Synopsis string
	Chapters []Chapter
}

type Chapter struct {
	URL        string
	Title      string
	Number     string
	UploadDate string
	ImageURLs  []string
}
