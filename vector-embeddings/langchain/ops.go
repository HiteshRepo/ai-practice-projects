package langchain

import "github.com/tmc/langchaingo/textsplitter"

func SplitDocuments(
	splitter string,
	document string) ([]string, error) {
	recurCh := textsplitter.NewRecursiveCharacter(
		textsplitter.WithSeparators([]string{splitter}),
		textsplitter.WithChunkSize(250),
		textsplitter.WithChunkOverlap(15),
	)

	return recurCh.SplitText(document)
}
