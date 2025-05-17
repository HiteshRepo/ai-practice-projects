package supabase

import (
	"pop-choice/constants"
	"pop-choice/models"
	"pop-choice/models/db"

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

func ReadDocumentByContent(
	tableName string,
	dbClient *supa.Client,
	content string,
) ([]db.Document, error) {
	var results []db.Document

	err := dbClient.
		DB.
		From(tableName).
		Select(constants.PopChoiceTblContentColumnName).
		Eq(constants.PopChoiceTblContentColumnName, content).
		Execute(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func ReadDocuments(
	tableName string,
	dbClient *supa.Client,
) ([]db.Document, error) {
	var results []db.Document

	err := dbClient.
		DB.
		From(tableName).
		Select(constants.PopChoiceTblContentColumnName).
		Execute(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// https://supabase.com/docs/guides/ai/vector-columns#querying-a-vector--embedding
func InvokeMatchFunction(
	dbClient *supa.Client,
	functionName string,
	embedding []float64,
	numMatches int,
) ([]db.MatchedDocument, error) {
	var results []db.MatchedDocument

	err := dbClient.DB.Rpc(functionName, map[string]any{
		"query_embedding": embedding,
		"match_threshold": 0.50,
		"match_count":     numMatches,
	}).Execute(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
