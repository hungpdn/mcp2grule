package dto

// Fact represents a set of key-value pairs used in rule evaluation
type Fact struct {
	M map[string]any
}

// NewFact creates a new Fact instance
func NewFact(m map[string]any) *Fact {
	if m == nil {
		m = map[string]any{}
	}
	return &Fact{M: m}
}

// Get retrieves a value by key
// Has checks if a key exists
// Set sets a value for a key
// AsMap returns the underlying map
func (f *Fact) Get(key string) any    { return f.M[key] }
func (f *Fact) Has(key string) bool   { _, ok := f.M[key]; return ok }
func (f *Fact) Set(key string, v any) { f.M[key] = v }
func (f *Fact) AsMap() map[string]any { return f.M }
