package supabase

import (
	supa "github.com/nedpals/supabase-go"
)

func NewClient(
	projectUrl string,
	apiKey string,
) *supa.Client {
	return supa.CreateClient(projectUrl, apiKey)
}
