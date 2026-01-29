package generator

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/sivchari/crx-registry/internal/types"
	"github.com/sivchari/crx-registry/internal/validator"
)

// Generator generates registry.yaml from pkgs/ directory.
type Generator struct {
	basePath  string
	validator *validator.Validator
}

// New creates a new Generator.
func New(basePath string) *Generator {
	return &Generator{
		basePath:  basePath,
		validator: validator.New(),
	}
}

// Result contains the generation result.
type Result struct {
	Registry *types.Registry
	Packages []*types.Package
	Errors   []PackageError
}

// PackageError represents an error for a specific package.
type PackageError struct {
	Filename string
	Err      error
}

func (e PackageError) Error() string {
	return fmt.Sprintf("%s: %v", e.Filename, e.Err)
}

// Generate scans pkgs/ directory and generates registry.yaml.
func (g *Generator) Generate() (*Result, error) {
	pkgsDir := filepath.Join(g.basePath, "pkgs")

	slog.Debug("reading packages directory", "dir", pkgsDir)
	entries, err := os.ReadDir(pkgsDir)
	if err != nil {
		return nil, fmt.Errorf("reading pkgs directory: %w", err)
	}

	result := &Result{
		Registry: &types.Registry{
			Version: types.RegistryVersion,
		},
	}

	for _, entry := range entries {
		if !g.isPackageFile(entry) {
			slog.Debug("skipping non-package file", "file", entry.Name())
			continue
		}

		filename := filepath.Join(pkgsDir, entry.Name())
		slog.Debug("loading package", "file", entry.Name())

		pkg, err := g.loadPackage(filename)
		if err != nil {
			slog.Debug("failed to load package", "file", entry.Name(), "error", err)
			result.Errors = append(result.Errors, PackageError{
				Filename: entry.Name(),
				Err:      err,
			})
			continue
		}

		if err := g.validator.ValidatePackage(pkg, filename); err != nil {
			slog.Debug("package validation failed", "file", entry.Name(), "error", err)
			result.Errors = append(result.Errors, PackageError{
				Filename: entry.Name(),
				Err:      err,
			})
			continue
		}

		slog.Debug("package loaded", "name", pkg.Name, "id", pkg.ID)
		result.Packages = append(result.Packages, pkg)
		result.Registry.Packages = append(result.Registry.Packages, pkg.Name)
	}

	slices.SortFunc(result.Packages, func(a, b *types.Package) int {
		return strings.Compare(a.Name, b.Name)
	})
	slices.Sort(result.Registry.Packages)

	slog.Debug("generation complete", "valid", len(result.Packages), "errors", len(result.Errors))
	return result, nil
}

// WriteRegistry writes the registry to registry.yaml.
func (g *Generator) WriteRegistry(reg *types.Registry) error {
	data, err := yaml.Marshal(reg)
	if err != nil {
		return fmt.Errorf("marshaling registry: %w", err)
	}

	path := filepath.Join(g.basePath, "registry.yaml")
	slog.Debug("writing registry", "path", path)

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing registry.yaml: %w", err)
	}

	return nil
}

func (g *Generator) isPackageFile(entry fs.DirEntry) bool {
	if entry.IsDir() {
		return false
	}
	name := entry.Name()
	return strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")
}

func (g *Generator) loadPackage(filename string) (*types.Package, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var pkg types.Package
	if err := yaml.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	return &pkg, nil
}
