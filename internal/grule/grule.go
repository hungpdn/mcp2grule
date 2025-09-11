package grule

import (
	"context"

	"github.com/hungpdn/grule-plus/engine"
	"github.com/hungpdn/mcp2grule/internal/api/dto"
	"github.com/hungpdn/mcp2grule/internal/config"
	"github.com/hungpdn/mcp2grule/internal/pkg/logger"
	"github.com/hungpdn/mcp2grule/internal/storage"
)

// IGrule is the interface for Grule service
type IGrule interface {
	Evaluate(ctx context.Context, in dto.EvaluateIn) (*dto.EvaluateOut, error)
	Create(ctx context.Context, in dto.CreateIn) (*dto.CreateOut, error)
	Update(ctx context.Context, name string, in dto.UpdateIn) (*dto.UpdateOut, error)
	Delete(ctx context.Context, name string) (*dto.DeleteOut, error)
	GetAll(ctx context.Context) (*dto.GetAllOut, error)
	GetByName(ctx context.Context, name string) (*dto.GetByNameOut, error)
}

// grule is the implementation of IGrule
type grule struct {
	cfg    config.Grule
	store  storage.IRulesetStorage
	engine engine.IGruleEngine
}

// New creates a new Grule service
func New(cfg config.Grule, store storage.IRulesetStorage) IGrule {

	engine := engine.NewSingleEngine(engine.Config{
		Type:            cfg.GetType(),
		Size:            cfg.Size,
		CleanupInterval: cfg.CleanupInterval,
		TTL:             cfg.TTL,
		FactName:        "Fact",
	})

	return &grule{
		cfg:    cfg,
		store:  store,
		engine: engine,
	}
}

// Evaluate evaluates the ruleset with the given facts
func (g *grule) Evaluate(ctx context.Context, in dto.EvaluateIn) (*dto.EvaluateOut, error) {

	rule, err := g.store.GetByName(ctx, in.RuleName)
	if err != nil {
		return nil, err
	}

	err = g.engine.Execute(ctx, rule.Name, &in.Facts)
	if err != nil {
		return nil, err
	}

	return &dto.EvaluateOut{ModifiedFacts: in.Facts.AsMap()}, nil
}

// Create creates a new ruleset
func (g *grule) Create(ctx context.Context, in dto.CreateIn) (*dto.CreateOut, error) {

	rule := storage.Ruleset{
		Name:        in.Name,
		Description: in.Description,
		Salience:    in.Salience,
		GRL:         in.GRL,
	}

	id, err := g.store.Create(ctx, rule)
	if err != nil {
		return nil, err
	}

	err = g.engine.AddRule(rule.Name, rule.GRL, 0)
	if err != nil {
		logger.Errorf("Failed to add rule %v: %v", rule.Name, err)
	}

	return &dto.CreateOut{ID: id}, nil
}

// Update updates an existing ruleset
func (g *grule) Update(ctx context.Context, name string, in dto.UpdateIn) (*dto.UpdateOut, error) {

	rule, err := g.store.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	rule.Description = in.Description
	rule.Salience = in.Salience
	rule.GRL = in.GRL

	err = g.store.Update(ctx, name, *rule)
	if err != nil {
		return nil, err
	}

	err = g.engine.BuildRule(rule.Name, rule.GRL, 0)
	if err != nil {
		logger.Errorf("Failed to build rule %v: %v", rule.Name, err)
	}

	return &dto.UpdateOut{Success: true}, nil
}

// Delete deletes a ruleset by name
func (g *grule) Delete(ctx context.Context, name string) (*dto.DeleteOut, error) {

	err := g.store.Delete(ctx, name)
	if err != nil {
		return nil, err
	}

	return &dto.DeleteOut{Success: true}, nil
}

// GetAll retrieves all rulesets
func (g *grule) GetAll(ctx context.Context) (*dto.GetAllOut, error) {

	rules, err := g.store.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.GetAllOut{Rulesets: rules}, nil
}

// GetByName retrieves a ruleset by name
func (g *grule) GetByName(ctx context.Context, name string) (*dto.GetByNameOut, error) {

	rule, err := g.store.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return &dto.GetByNameOut{Ruleset: *rule}, nil
}
