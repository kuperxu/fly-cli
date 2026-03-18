package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fly/internal/model"
	"gopkg.in/yaml.v3"
)

const configFileName = "portfolio.yaml"

// Portfolio holds all user positions
type Portfolio struct {
	Holdings []*model.Holding `yaml:"holdings"`
}

// Store manages reading/writing the portfolio config file
type Store struct {
	path string
}

// NewStore creates a Store using ~/.fly-cli/portfolio.yaml
func NewStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot find home directory: %w", err)
	}
	dir := filepath.Join(home, ".fly-cli")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("cannot create config dir: %w", err)
	}
	return &Store{path: filepath.Join(dir, configFileName)}, nil
}

// Load reads the portfolio from disk. Returns empty portfolio if file doesn't exist.
func (s *Store) Load() (*Portfolio, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return &Portfolio{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cannot read portfolio file: %w", err)
	}
	var p Portfolio
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("invalid portfolio file: %w", err)
	}
	return &p, nil
}

// Save writes the portfolio to disk
func (s *Store) Save(p *Portfolio) error {
	data, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Errorf("cannot serialize portfolio: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("cannot write portfolio file: %w", err)
	}
	return nil
}

// Upsert adds or updates a holding by code
func (s *Store) Upsert(h *model.Holding) error {
	p, err := s.Load()
	if err != nil {
		return err
	}
	normalized := normalizeKey(h.Code)
	for i, existing := range p.Holdings {
		if normalizeKey(existing.Code) == normalized {
			p.Holdings[i] = h
			return s.Save(p)
		}
	}
	p.Holdings = append(p.Holdings, h)
	return s.Save(p)
}

// Remove deletes a holding by code. Returns error if not found.
func (s *Store) Remove(code string) error {
	p, err := s.Load()
	if err != nil {
		return err
	}
	normalized := normalizeKey(code)
	newHoldings := p.Holdings[:0]
	found := false
	for _, h := range p.Holdings {
		if normalizeKey(h.Code) == normalized {
			found = true
			continue
		}
		newHoldings = append(newHoldings, h)
	}
	if !found {
		return fmt.Errorf("holding not found: %s", code)
	}
	p.Holdings = newHoldings
	return s.Save(p)
}

// FindHolding returns a holding by code, or nil if not found
func (p *Portfolio) FindHolding(code string) *model.Holding {
	normalized := normalizeKey(code)
	for _, h := range p.Holdings {
		if normalizeKey(h.Code) == normalized {
			return h
		}
	}
	return nil
}

func normalizeKey(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}
