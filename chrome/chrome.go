package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

type LLMTool interface {
	tools.Tool
	SetLLM(llms.LLM)
}

// LLMGuidedCrawler implements the Tool interface
type LLMGuidedCrawler struct {
	Tool *ChromedpTool // Scraping tool (e.g., Chromedp)
	Llm  llms.Model    // LLM model (e.g., GPT)
}

func (c *LLMGuidedCrawler) SetLLM(llm llms.LLM) {
	c.Llm = llm
}

// Name returns the name of the tool
func (c *LLMGuidedCrawler) Name() string {
	return "LLM Guided Web Crawler"
}

// Description returns a brief description of what the tool does
func (c *LLMGuidedCrawler) Description() string {
	return "Crawls web pages using an LLM to identify the presence of answers to search queries."
}

// Call performs the crawling process, stopping when the answer is found
func (c *LLMGuidedCrawler) Call(ctx context.Context, query string) (string, error) {
	// Start by searching on Google with the query
	searchURL := fmt.Sprintf("https://www.google.com/search?q=%s", query)

	// Initialize scraping context
	result := ""
	for i := 0; i < 5; i++ { // Limit to 5 pages for example
		fmt.Printf("Crawling page %d...\n", i+1)

		// Scrape the search results
		pageContent, err := c.Tool.Run(searchURL)
		if err != nil {
			return "", fmt.Errorf("error scraping page: %w", err)
		}

		// Use LLM to check if the answer is present on the page
		foundAnswer, answer, err := c.checkForAnswer(ctx, pageContent, query)
		if err != nil {
			return "", fmt.Errorf("error checking LLM: %w", err)
		}

		// If the LLM detects the answer, stop early
		if foundAnswer {
			result = answer // Store the relevant page content
			fmt.Println("Answer found!")
			break
		}

		// Move to the next page (pagination URL)
		searchURL = fmt.Sprintf("https://www.google.com/search?q=%s&start=%d", query, (i+1)*10)
	}

	if result == "" {
		return "", fmt.Errorf("answer not found within page limit")
	}

	return result, nil
}

// checkForAnswer uses the LLM to decide if the page contains the answer
func (c *LLMGuidedCrawler) checkForAnswer(ctx context.Context, content string, query string) (bool, string, error) {
	// LLM prompt to evaluate if the answer to the query is present in the content
	prompt := fmt.Sprintf(`
		You are an expert at understanding content. Analyze the following content and determine if it contains the answer to the query: "%s".
		If the answer is found, return the answer, otherwise return ''.
		Content: %s
	`, query, content)

	// Generate response from LLM
	response, err := llms.GenerateFromSinglePrompt(ctx, c.Llm, prompt)
	if err != nil {
		return false, "", fmt.Errorf("error generating LLM response: %w", err)
	}

	// Check if the LLM confirms that the answer is present
	return strings.ToLower(response) != "", response, nil
}

// ChromedpTool scrapes a webpage content
type ChromedpTool struct{}

func (t *ChromedpTool) Run(input string) (string, error) {
	// Create context for scraping
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout for scraping task
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var result string
	err := chromedp.Run(ctx,
		chromedp.Navigate(input),
		chromedp.Text("body", &result),
	)
	if err != nil {
		return "", fmt.Errorf("error running Chromedp tool: %w", err)
	}

	return result, nil
}

// Export the tool so it can be accessed via the plugin system
var Tool LLMGuidedCrawler
