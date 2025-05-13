package db

type MatchedDocument struct {
	ID         int     `json:"id"`
	Content    string  `json:"content"`
	Similarity float64 `json:"similarity"`
}
