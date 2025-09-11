package storage

import (
	"context"
	"sync"
	"time"

	"github.com/hungpdn/mcp2grule/internal/utils"
)

// memory is an in-memory implementation of IRulesetStorage interface
type memory struct {
	mu sync.RWMutex
	m  map[string]Ruleset
}

// NewMemory creates a new memory storage
func NewMemory() *memory {
	return &memory{m: map[string]Ruleset{}}
}

// key generates map key for a ruleset name
func (s *memory) key(name string) string {
	return name
}

// GetAll returns all rulesets
func (s *memory) GetAll(ctx context.Context) ([]Ruleset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Ruleset, 0, len(s.m))
	for _, v := range s.m {
		out = append(out, v)
	}

	return out, nil
}

// GetByName returns a ruleset by name
func (s *memory) GetByName(ctx context.Context, name string) (*Ruleset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	r, ok := s.m[s.key(name)]
	if !ok {
		return nil, ErrNotFound
	}

	return &r, nil
}

// Create creates a new ruleset
func (s *memory) Create(ctx context.Context, rule Ruleset) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.m[s.key(rule.Name)]; ok {
		return "", ErrAlreadyExists
	}

	rule.ID = utils.NewULID()
	rule.CreatedAt = time.Now().Unix()
	rule.UpdatedAt = time.Now().Unix()

	s.m[s.key(rule.Name)] = rule

	return rule.ID, nil
}

// Update updates an existing ruleset
func (s *memory) Update(ctx context.Context, name string, rule Ruleset) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := s.key(name)
	if _, ok := s.m[k]; !ok {
		return ErrNotFound
	}

	rule.UpdatedAt = time.Now().Unix()
	s.m[k] = rule

	return nil
}

// Delete deletes a ruleset by name
func (s *memory) Delete(ctx context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := s.key(name)
	if _, ok := s.m[k]; !ok {
		return ErrNotFound
	}
	delete(s.m, k)

	return nil
}
