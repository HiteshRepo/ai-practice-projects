package main

import (
	"context"
	"flag"
	"log"
	openaipkg "react/openai"
	"react/versions"

	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

type envvars struct {
	OpenApiKey string `env:"OPEN_API_KEY"`
}

func main() {
	ctx := context.Background()

	var envs envvars
	if err := env.Parse(&envs); err != nil {
		log.Fatalln("failed to parse env variables", errors.Wrap(err, "missing required env"))
	}

	versionFlag := flag.String(
		"version",
		"v3",
		"Allowed values: v1, v2, v3.")
	queryFlag := flag.String(
		"query",
		"",
		"your query related to a location and its weather")
	flag.Parse()

	version := ""
	if versionFlag != nil {
		version = *versionFlag
	}

	query := ""
	if queryFlag != nil {
		query = *queryFlag
	}

	openaiClient := openaipkg.NewOpenAiClient(envs.OpenApiKey)

	// go run main.go -version=v1/v2/v3 -query="Give me a list of activity ideas based on my current location and weather"
	switch version {
	case "v1":
		versions.V1(ctx, openaiClient, query)
	case "v2":
		versions.V2(ctx, openaiClient, query)
	case "v3":
		versions.V3(ctx, openaiClient, query)
	}
}
