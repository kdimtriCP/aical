package openai

import (
	"context"
	"errors"
	"sync"
)

type FunctionDescription struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type function func(ctx context.Context, arguments string) string

type Registry struct {
	mu        sync.RWMutex
	functions map[string]function
	descs     map[string]FunctionDescription
}

func NewRegistry() *Registry {
	return &Registry{
		functions: make(map[string]function),
		descs:     make(map[string]FunctionDescription),
	}
}

func (r *Registry) Register(name string, desc FunctionDescription, function function) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.functions[name] = function
	r.descs[name] = desc
	return
}

func (r *Registry) Execute(ctx context.Context, name string, arguments string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	function, exists := r.functions[name]
	if !exists {
		return errors.New("function not found").Error()
	}
	return function(ctx, arguments)
}

func (r *Registry) Descriptions() []FunctionDescription {
	r.mu.RLock()
	defer r.mu.RUnlock()

	descs := make([]FunctionDescription, 0, len(r.descs))
	for _, desc := range r.descs {
		descs = append(descs, desc)
	}

	return descs
}
