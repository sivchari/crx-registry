package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/sivchari/crx-registry/internal/generator"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		slog.Error("execution failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		showVersion bool
		basePath    string
		debug       bool
	)
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.StringVar(&basePath, "path", ".", "path to registry root")
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.Parse()

	initLogger(debug)

	if showVersion {
		fmt.Printf("crx-registry %s\n", version)
		return nil
	}

	args := flag.Args()
	if len(args) == 0 {
		return fmt.Errorf("no command specified. use: generate, validate")
	}

	cmd := args[0]
	slog.Debug("executing command", "command", cmd, "path", basePath)

	switch cmd {
	case "generate":
		return runGenerate(basePath)
	case "validate":
		return runValidate(basePath)
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func initLogger(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	slog.SetDefault(slog.New(handler))
}

func runGenerate(basePath string) error {
	gen := generator.New(basePath)

	slog.Info("scanning packages", "path", basePath)
	result, err := gen.Generate()
	if err != nil {
		return err
	}

	if len(result.Errors) > 0 {
		for _, e := range result.Errors {
			slog.Error("validation failed", "file", e.Filename, "error", e.Err)
		}
		return fmt.Errorf("%d package(s) failed validation", len(result.Errors))
	}

	if err := gen.WriteRegistry(result.Registry); err != nil {
		return err
	}

	slog.Info("generated registry.yaml", "packages", len(result.Packages))
	for _, pkg := range result.Packages {
		slog.Debug("included package", "name", pkg.Name, "id", pkg.ID)
	}

	return nil
}

func runValidate(basePath string) error {
	gen := generator.New(basePath)

	slog.Info("validating packages", "path", basePath)
	result, err := gen.Generate()
	if err != nil {
		return err
	}

	if len(result.Errors) > 0 {
		for _, e := range result.Errors {
			slog.Error("validation failed", "file", e.Filename, "error", e.Err)
		}
		return fmt.Errorf("%d package(s) failed validation", len(result.Errors))
	}

	slog.Info("all packages valid", "count", len(result.Packages))
	return nil
}
