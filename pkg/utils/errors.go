package utils

import (
	"fmt"

	"github.com/paphachanok/modelgene/pkg/types"
)

// ModelGeneError represents a custom error from modelgene
type ModelGeneError struct {
	Provider types.Provider
	Message  string
	Err      error
}

func (e *ModelGeneError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("modelgene [%s]: %s: %v", e.Provider, e.Message, e.Err)
	}
	return fmt.Sprintf("modelgene [%s]: %s", e.Provider, e.Message)
}

// NewError creates a ModelGeneError
func NewError(provider types.Provider, message string, err error) error {
	return &ModelGeneError{
		Provider: provider,
		Message:  message,
		Err:      err,
	}
}
