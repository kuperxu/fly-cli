package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"fly/internal/model"
	"gopkg.in/yaml.v3"
)

const alertFileName = "alerts.yaml"

// AlertBook holds all user price alerts
type AlertBook struct {
	Alerts []*model.Alert `yaml:"alerts"`
}

// AlertStore manages reading/writing the alerts config file
type AlertStore struct {
	path string
}

// NewAlertStore creates an AlertStore using ~/.fly-cli/alerts.yaml
func NewAlertStore() (*AlertStore, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot find home directory: %w", err)
	}
	dir := filepath.Join(home, ".fly-cli")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("cannot create config dir: %w", err)
	}
	return &AlertStore{path: filepath.Join(dir, alertFileName)}, nil
}

// Load reads all alerts from disk. Returns empty AlertBook if file doesn't exist.
func (s *AlertStore) Load() (*AlertBook, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return &AlertBook{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cannot read alerts file: %w", err)
	}
	var book AlertBook
	if err := yaml.Unmarshal(data, &book); err != nil {
		return nil, fmt.Errorf("invalid alerts file: %w", err)
	}
	return &book, nil
}

// Save writes alerts to disk
func (s *AlertStore) Save(book *AlertBook) error {
	data, err := yaml.Marshal(book)
	if err != nil {
		return fmt.Errorf("cannot serialize alerts: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("cannot write alerts file: %w", err)
	}
	return nil
}

// Upsert adds or updates an alert by code.
// For updates, only non-zero fields in the input overwrite existing values.
func (s *AlertStore) Upsert(a *model.Alert) error {
	book, err := s.Load()
	if err != nil {
		return err
	}
	normalized := normalizeKey(a.Code)
	for i, existing := range book.Alerts {
		if normalizeKey(existing.Code) == normalized {
			// Merge: only overwrite fields that are explicitly set (non-zero)
			if a.Entry != 0 {
				existing.Entry = a.Entry
			}
			if a.TP1 != 0 {
				existing.TP1 = a.TP1
			}
			if a.SL != 0 {
				existing.SL = a.SL
			}
			book.Alerts[i] = existing
			return s.Save(book)
		}
	}
	book.Alerts = append(book.Alerts, a)
	return s.Save(book)
}

// Remove deletes an alert by code. Returns error if not found.
func (s *AlertStore) Remove(code string) error {
	book, err := s.Load()
	if err != nil {
		return err
	}
	normalized := normalizeKey(code)
	newAlerts := book.Alerts[:0]
	found := false
	for _, a := range book.Alerts {
		if normalizeKey(a.Code) == normalized {
			found = true
			continue
		}
		newAlerts = append(newAlerts, a)
	}
	if !found {
		return fmt.Errorf("alert not found: %s", code)
	}
	book.Alerts = newAlerts
	return s.Save(book)
}

// FindAlert returns an alert by code, or nil if not found
func (book *AlertBook) FindAlert(code string) *model.Alert {
	normalized := normalizeKey(code)
	for _, a := range book.Alerts {
		if normalizeKey(a.Code) == normalized {
			return a
		}
	}
	return nil
}
