package api

import (
	"fmt"

	"fly/internal/model"
)

// Client fetches stock quotes with primary + fallback provider
type Client struct {
	primary  Provider
	fallback Provider
}

// NewClient creates a Client using Eastmoney as primary and Tencent as fallback
func NewClient() *Client {
	return &Client{
		primary:  newEastmoneyProvider(),
		fallback: newTencentProvider(),
	}
}

// GetQuotes fetches quotes for the given symbols, falling back to secondary provider on error
func (c *Client) GetQuotes(symbols []string) ([]*model.Quote, error) {
	quotes, err := c.primary.GetQuotes(symbols)
	if err == nil && len(quotes) == len(symbols) {
		return quotes, nil
	}
	// Try fallback
	return c.fallback.GetQuotes(symbols)
}

// GetQuote fetches a single stock quote
func (c *Client) GetQuote(symbol string) (*model.Quote, error) {
	quotes, err := c.GetQuotes([]string{symbol})
	if err != nil {
		return nil, err
	}
	if len(quotes) == 0 {
		return nil, nil
	}
	return quotes[0], nil
}

// GetIndexQuotes fetches quotes for a list of market indices/bonds using raw secids
func (c *Client) GetIndexQuotes(secids []string) ([]*model.Quote, error) {
	em, ok := c.primary.(*eastmoneyProvider)
	if !ok {
		return nil, fmt.Errorf("primary provider is not eastmoney")
	}
	return em.fetchIndices(secids)
}
