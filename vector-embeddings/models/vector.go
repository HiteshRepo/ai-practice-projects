package models

type Vector struct {
	Content   string    `json:"content"`
	Embedding []float64 `json:"embedding"`
}
