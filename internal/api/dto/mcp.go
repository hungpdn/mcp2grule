package dto

import "github.com/hungpdn/mcp2grule/internal/storage"

// EvaluateIn is the input structure for Evaluate method
type EvaluateIn struct {
	Facts    Fact   `json:"facts" jsonschema:"Facts to be evaluated"`
	RuleName string `json:"rule_name" jsonschema:"Name of the ruleset to be used for evaluation"`
}

// EvaluateOut is the output structure for Evaluate method
type EvaluateOut struct {
	ModifiedFacts map[string]any `json:"modified_facts" jsonschema:"Modified facts after evaluation"`
}

// CreateIn is the input structure for Create method
type CreateIn struct {
	Name        string `json:"name" jsonschema:"Name of the ruleset"`
	Description string `json:"description" jsonschema:"Description of the ruleset"`
	Salience    int    `json:"salience" jsonschema:"Priority of the rule"`
	GRL         string `json:"grl" jsonschema:"The actual GRL content"`
}

// CreateOut is the output structure for Create method
type CreateOut struct {
	ID string `json:"id" jsonschema:"ID of the created ruleset"`
}

// UpdateIn is the input structure for Update method
type UpdateIn struct {
	Name        string `json:"name" jsonschema:"Name of the ruleset"`
	Description string `json:"description" jsonschema:"Description of the ruleset"`
	Salience    int    `json:"salience" jsonschema:"Priority of the rule"`
	GRL         string `json:"grl" jsonschema:"The actual GRL content"`
}

// UpdateOut is the output structure for Update method
type UpdateOut struct {
	Success bool `json:"success" jsonschema:"Indicates if the update was successful"`
}

// DeleteIn is the input structure for Delete method
type DeleteIn struct {
	Name string `json:"name" jsonschema:"Name of the ruleset to be deleted"`
}

// DeleteOut is the output structure for Delete method
type DeleteOut struct {
	Success bool `json:"success" jsonschema:"Indicates if the deletion was successful"`
}

// GetByNameIn is the input structure for GetByName method
type GetByNameIn struct {
	Name string `json:"name" jsonschema:"Name of the ruleset to retrieve"`
}

// GetByNameOut is the output structure for GetByName method
type GetByNameOut struct {
	Ruleset storage.Ruleset `json:"ruleset" jsonschema:"Retrieved ruleset"`
}

// GetAllOut is the output structure for GetAll method
type GetAllOut struct {
	Rulesets []storage.Ruleset `json:"rulesets" jsonschema:"List of all rulesets"`
}
