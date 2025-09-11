package storage

import (
	"context"
	"errors"
)

// Common database errors
var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrDatabase      = errors.New("database error")
)

// Ruleset represents the structure of a business rule in the database.
type Ruleset struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Salience    int    `json:"salience"`   // Priority of the rule
	GRL         string `json:"grl"`        // The actual GRL content
	CreatedAt   int64  `json:"created_at"` // Unix timestamp
	UpdatedAt   int64  `json:"updated_at"` // Unix timestamp
}

// IRulesetStorage defines the interface for database operations on rules.
// This allows for different database backends to be implemented.
type IRulesetStorage interface {
	// GetAll retrieves all rulesets.
	GetAll(ctx context.Context) ([]Ruleset, error)
	// GetByName retrieves a ruleset by its name.
	GetByName(ctx context.Context, name string) (*Ruleset, error)
	// Create adds a new ruleset to the database.
	Create(ctx context.Context, rule Ruleset) (string, error)
	// Update modifies an existing ruleset identified by name.
	Update(ctx context.Context, name string, rule Ruleset) error
	// Delete removes a ruleset identified by name.
	Delete(ctx context.Context, name string) error
}
