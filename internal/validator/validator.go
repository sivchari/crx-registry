package validator

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sivchari/crx-registry/internal/types"
)

// ValidationError represents a validation error with context.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// Validator validates packages and registry files.
type Validator struct{}

// New creates a new Validator.
func New() *Validator {
	return &Validator{}
}

// ValidatePackage validates a package definition.
func (v *Validator) ValidatePackage(pkg *types.Package, filename string) error {
	var errs ValidationErrors

	// Required fields
	if pkg.Name == "" {
		errs = append(errs, ValidationError{Field: "name", Message: "required"})
	}
	if pkg.ID == "" {
		errs = append(errs, ValidationError{Field: "id", Message: "required"})
	}
	if pkg.DisplayName == "" {
		errs = append(errs, ValidationError{Field: "display_name", Message: "required"})
	}

	// Name must match filename
	expectedName := strings.TrimSuffix(filepath.Base(filename), ".yaml")
	if pkg.Name != expectedName {
		errs = append(errs, ValidationError{
			Field:   "name",
			Message: fmt.Sprintf("must match filename (expected %q, got %q)", expectedName, pkg.Name),
		})
	}

	// Name must be kebab-case
	if pkg.Name != "" && !types.KebabCasePattern.MatchString(pkg.Name) {
		errs = append(errs, ValidationError{Field: "name", Message: "must be kebab-case"})
	}

	// Extension ID must be valid
	if pkg.ID != "" && !types.ExtensionIDPattern.MatchString(pkg.ID) {
		errs = append(errs, ValidationError{
			Field:   "id",
			Message: "must be 32 lowercase letters (a-p)",
		})
	}

	// Optional URL fields must be valid if present
	if pkg.Homepage != "" && !types.URLPattern.MatchString(pkg.Homepage) {
		errs = append(errs, ValidationError{Field: "homepage", Message: "must be a valid URL"})
	}
	if pkg.Repository != "" && !types.URLPattern.MatchString(pkg.Repository) {
		errs = append(errs, ValidationError{Field: "repository", Message: "must be a valid URL"})
	}

	// Tags must be kebab-case
	for i, tag := range pkg.Tags {
		if !types.KebabCasePattern.MatchString(tag) {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("tags[%d]", i),
				Message: fmt.Sprintf("%q must be kebab-case", tag),
			})
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// ValidateRegistry validates a registry index.
func (v *Validator) ValidateRegistry(reg *types.Registry) error {
	var errs ValidationErrors

	if reg.Version < 1 {
		errs = append(errs, ValidationError{Field: "version", Message: "must be >= 1"})
	}

	seen := make(map[string]bool)
	for i, name := range reg.Packages {
		if name == "" {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("packages[%d]", i),
				Message: "empty package name",
			})
			continue
		}
		if seen[name] {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("packages[%d]", i),
				Message: fmt.Sprintf("duplicate package %q", name),
			})
		}
		seen[name] = true
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// IsValidationError checks if an error is a validation error.
func IsValidationError(err error) bool {
	var ve ValidationErrors
	var single *ValidationError
	return errors.As(err, &ve) || errors.As(err, &single)
}
