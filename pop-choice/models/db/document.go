package db

type Document struct {
	Content   string `json:"content"`
	Embedding string `json:"embedding"`
}
