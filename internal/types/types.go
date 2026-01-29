package types

import "regexp"

// Package represents a Chrome extension package in the registry.
type Package struct {
	Name        string   `yaml:"name"`
	ID          string   `yaml:"id"`
	DisplayName string   `yaml:"display_name"`
	Description string   `yaml:"description,omitempty"`
	Homepage    string   `yaml:"homepage,omitempty"`
	Repository  string   `yaml:"repository,omitempty"`
	Tags        []string `yaml:"tags,omitempty"`
}

// Registry represents the registry index.
type Registry struct {
	Version  int      `yaml:"version"`
	Packages []string `yaml:"packages"`
}

// RegistryVersion is the current version of the registry format.
const RegistryVersion = 1

// Validation patterns.
var (
	// ExtensionIDPattern matches valid Chrome Extension IDs (32 lowercase letters a-p).
	ExtensionIDPattern = regexp.MustCompile(`^[a-p]{32}$`)

	// KebabCasePattern matches kebab-case strings.
	KebabCasePattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

	// URLPattern matches valid HTTP/HTTPS URLs.
	URLPattern = regexp.MustCompile(`^https?://[^\s]+$`)
)
