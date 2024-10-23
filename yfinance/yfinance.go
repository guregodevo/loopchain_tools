package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// YahooFinanceStockPriceInput represents the structure of the expected input.
type YahooFinanceStockPriceInput struct {
	StockSymbol string `json:"stock_symbol"`
	StockName   string `json:"stock_name"`
	Partition   string `json:"partition"` // Date in the format YYYY-MM-DD
}

// YahooFinanceStockPriceTool implements the Tool interface for fetching stock prices.
type YahooFinanceStockPriceTool struct{}

// Name returns the name of the tool.
func (t *YahooFinanceStockPriceTool) Name() string {
	return "yahoo_finance_stock_price"
}

// Description provides a brief description of the tool's functionality and the exact JSON input format.
func (t *YahooFinanceStockPriceTool) Description() string {
	return "Tool that fetches stock prices for a specific company and date from Yahoo Finance. " +
		"Input should be a JSON object with the following format: \n" +
		`{
			"stock_symbol": "AAPL",
			"stock_name": "Apple",
			"partition": "2023-10-16"
		}` +
		"\n" +
		"The stock symbol is the company's ticker symbol, stock_name is the name of the company, " +
		"and partition is the specific date in YYYY-MM-DD format."
}

// Call executes the tool using the provided context and input.
func (t *YahooFinanceStockPriceTool) Call(ctx context.Context, input string) (string, error) {
	// Parse the input string (JSON) into the YahooFinanceStockPriceInput struct
	var params YahooFinanceStockPriceInput
	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return "", fmt.Errorf("invalid input format: %v", err)
	}

	// Use the parsed stock symbol and date to get historical stock data from Yahoo Finance
	stockPrice, err := t.getYahooFinanceStockPrice(ctx, params.StockSymbol, params.Partition)
	if err != nil {
		return "", err
	}

	// Return the result as a formatted string
	return fmt.Sprintf("The stock price of %s (%s) on %s was %f.", params.StockSymbol, params.StockName, params.Partition, stockPrice), nil
}

// getYahooFinanceStockPrice queries Yahoo Finance API for historical stock price by ticker and date.
func (t *YahooFinanceStockPriceTool) getYahooFinanceStockPrice(ctx context.Context, stockSymbol, date string) (float64, error) {
	apiURL := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d&period1=%d&period2=%d", stockSymbol, toUnix(date), toUnix(date)+86400)

	// Create an HTTP request with the provided context to support cancellation or timeout
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to query Yahoo Finance: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to parse Yahoo Finance response: %v", err)
	}

	// Extract the stock price from the response
	chart, ok := result["chart"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid response structure")
	}

	resultArr, ok := chart["result"].([]interface{})
	if !ok || len(resultArr) == 0 {
		return 0, fmt.Errorf("no stock price data available")
	}

	firstResult := resultArr[0].(map[string]interface{})
	indicators, ok := firstResult["indicators"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid indicators structure")
	}

	quote, ok := indicators["quote"].([]interface{})
	if !ok || len(quote) == 0 {
		return 0, fmt.Errorf("no stock quote data available")
	}

	quoteData := quote[0].(map[string]interface{})
	closePrices, ok := quoteData["close"].([]interface{})
	if !ok || len(closePrices) == 0 {
		return 0, fmt.Errorf("no close price data available")
	}

	// Return the close price for the date
	closePrice := closePrices[0].(float64)
	return closePrice, nil
}

// toUnix converts a date string (YYYY-MM-DD) to Unix timestamp
func toUnix(date string) int64 {
	layout := "2006-01-02"
	t, _ := time.Parse(layout, date)
	return t.Unix()
}

var Tool YahooFinanceStockPriceTool
