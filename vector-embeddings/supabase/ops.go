package supabase

import (
	"vector-embeddings/models"
	"vector-embeddings/models/db"

	supa "github.com/nedpals/supabase-go"
)

func InsertDocument(
	tableName string,
	dbClient *supa.Client,
	doc models.Vector,
) ([]db.Document, error) {
	var results []db.Document

	err := dbClient.DB.From(tableName).Insert(doc).Execute(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
