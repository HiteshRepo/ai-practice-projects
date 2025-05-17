package models

type Vector struct {
	Content   string    `json:"content"`
	Embedding []float64 `json:"embedding"`
}

type Movie struct {
	Title       string `json:"title"`
	ReleaseYear string `json:"releaseYear"`
	Content     string `json:"content"`
}

func (m Movie) ToString() string {
	return "Title: " + m.Title + "\nRelease Year: " + m.ReleaseYear + "\nContent: " + m.Content
}
