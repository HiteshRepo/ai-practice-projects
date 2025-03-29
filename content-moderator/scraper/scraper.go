package scraper

import (
	"bytes"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
)

func ScrapeURL(url string) (string, error) {
	if url == "" {
		return "", errors.New("empty URL provided")
	}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"),
	)

	var textBuffer bytes.Buffer

	c.OnHTML("script, style, meta, link, noscript", func(e *colly.HTMLElement) {
		// Do nothing, effectively skipping these elements
	})

	c.OnHTML("p, h1, h2, h3, h4, h5, h6, span, div, li, td, th, a", func(e *colly.HTMLElement) {
		// Get the text and clean it
		text := strings.TrimSpace(e.Text)
		if text != "" {
			textBuffer.WriteString(text)
			textBuffer.WriteString("\n")
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit(url)
	if err != nil {
		return "", errors.Wrap(err, "failed to visit URL")
	}

	c.Wait()

	return cleanupText(textBuffer.String()), nil
}

func cleanupText(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.Join(strings.Fields(text), " ")

	lines := strings.Split(text, "\n")
	var cleanLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanLines = append(cleanLines, trimmed)
		}
	}

	return strings.Join(cleanLines, "\n")
}
