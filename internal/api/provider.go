package api

import (
	"fly/internal/model"
)

// Provider is the interface for stock data sources
type Provider interface {
	Name() string
	GetQuotes(symbols []string) ([]*model.Quote, error)
}
