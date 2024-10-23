package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// YahooFinanceNewsInput represents the structure of the expected input.
type YahooFinanceNewsInput struct {
	Query string `json:"query"`
}

// YahooFinanceNewsTool implements the Tool interface.
type YahooFinanceNewsTool struct{}

// Name returns the name of the tool.
func (t *YahooFinanceNewsTool) Name() string {
	return "yahoo_finance_news"
}

// Description provides a brief description of the tool's functionality and the exact JSON input format.
func (t *YahooFinanceNewsTool) Description() string {
	return "Tool that searches financial news on Yahoo Finance. " +
		"Useful for when you need to find financial news about a public company. " +
		"Input should be a JSON object in the following format: \n" +
		`{
			"query": "AAPL"
		}` +
		"\n" +
		"Where 'query' is the company's ticker symbol, such as 'AAPL' for Apple or 'MSFT' for Microsoft."
}

// Call executes the tool using the provided context and input.
func (t *YahooFinanceNewsTool) Call(ctx context.Context, input string) (string, error) {
	// Parse the input string (JSON) into the YahooFinanceNewsInput struct
	ticker := input
	fmt.Printf("Received input: %s\n", input)

	var params YahooFinanceNewsInput
	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return "", fmt.Errorf("invalid input format: %v", err)
	}
	ticker = params.Query

	// Use the parsed query to get Yahoo Finance news
	news, err := t.getYahooFinanceNews(ctx, ticker)
	if err != nil {
		return "", err
	}

	// If no news is found
	if len(news) == 0 {
		return fmt.Sprintf("No news found for company ticker %s", ticker), nil
	}

	// Format the results
	return t.formatResults(news), nil
}

// getYahooFinanceNews queries Yahoo Finance API for company news by ticker.
func (t *YahooFinanceNewsTool) getYahooFinanceNews(ctx context.Context, ticker string) ([]map[string]interface{}, error) {
	apiURL := fmt.Sprintf("https://query1.finance.yahoo.com/v10/finance/quoteSummary/%s?modules=news", ticker)

	// Create an HTTP request with the provided context to support cancellation or timeout
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query Yahoo Finance: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse Yahoo Finance response: %v", err)
	}

	// Extract news articles
	newsItems := extractNews(result)

	if len(newsItems) == 0 {
		return nil, fmt.Errorf("no news found for ticker %s", ticker)
	}

	return newsItems, nil
}

// extractNews parses the API response and extracts news articles.
func extractNews(response map[string]interface{}) []map[string]interface{} {
	quoteSummary, ok := response["quoteSummary"].(map[string]interface{})
	if !ok {
		return nil
	}

	result, ok := quoteSummary["result"].([]interface{})
	if !ok || len(result) == 0 {
		return nil
	}

	firstResult := result[0].(map[string]interface{})
	news, ok := firstResult["news"].([]interface{})
	if !ok {
		return nil
	}

	var newsItems []map[string]interface{}
	for _, item := range news {
		newsItem, ok := item.(map[string]interface{})
		if ok {
			newsItems = append(newsItems, newsItem)
		}
	}
	return newsItems
}

// formatResults formats the news into a human-readable format.
func (t *YahooFinanceNewsTool) formatResults(news []map[string]interface{}) string {
	var results []string
	for _, article := range news {
		title, _ := article["title"].(string)
		description, _ := article["description"].(string)
		link, _ := article["link"].(string)
		results = append(results, fmt.Sprintf("%s\n%s\nLink: %s", title, description, link))
	}
	return strings.Join(results, "\n\n")
}

var Tool YahooFinanceNewsTool
