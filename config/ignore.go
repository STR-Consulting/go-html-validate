package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// IgnoreFileName is the name of the ignore file.
const IgnoreFileName = ".htmlvalidateignore"

// LoadIgnorePatterns searches for and loads .htmlvalidateignore from dir upward.
// Returns nil if no ignore file is found.
func LoadIgnorePatterns(dir string) ([]string, error) {
	path, err := FindIgnoreFile(dir)
	if err != nil {
		return nil, err
	}
	if path == "" {
		return nil, nil
	}
	return LoadIgnoreFile(path)
}

// LoadIgnoreFile loads patterns from a specific ignore file.
func LoadIgnoreFile(path string) ([]string, error) {
	f, err := os.Open(path) //nolint:gosec // user-specified ignore file path
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var patterns []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return patterns, nil
}

// FindIgnoreFile searches for .htmlvalidateignore from dir upward.
// Returns empty string if no ignore file is found.
func FindIgnoreFile(dir string) (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	for {
		path := filepath.Join(absDir, IgnoreFileName)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}

		parent := filepath.Dir(absDir)
		if parent == absDir {
			// Reached root
			return "", nil
		}
		absDir = parent
	}
}

// MatchesIgnorePattern checks if a path matches any of the ignore patterns.
// Supports gitignore-style patterns:
//   - Regular globs: *.html, test_*.go
//   - Directory patterns (ending with /): node_modules/, vendor/
//   - Recursive patterns: **/*.generated.html
//   - Negation: !important.html (not implemented yet)
func MatchesIgnorePattern(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if matchPattern(path, pattern) {
			return true
		}
	}
	return false
}

func matchPattern(path, pattern string) bool {
	// Handle negation (skip for now, would need to track and apply later)
	if strings.HasPrefix(pattern, "!") {
		return false
	}

	// Handle directory patterns (ending with /)
	if strings.HasSuffix(pattern, "/") {
		dir := strings.TrimSuffix(pattern, "/")
		// Match if path contains this directory component
		if strings.Contains(path, "/"+dir+"/") ||
			strings.HasPrefix(path, dir+"/") ||
			path == dir {
			return true
		}
		return false
	}

	// Handle ** recursive patterns
	if strings.Contains(pattern, "**") {
		return matchDoublestar(path, pattern)
	}

	// Simple glob matching against basename
	base := filepath.Base(path)
	if matched, _ := filepath.Match(pattern, base); matched {
		return true
	}

	// Try matching full path
	if matched, _ := filepath.Match(pattern, path); matched {
		return true
	}

	return false
}

// matchDoublestar handles ** patterns.
func matchDoublestar(path, pattern string) bool {
	// Split pattern by **
	parts := strings.Split(pattern, "**")

	if len(parts) == 2 {
		prefix := strings.TrimSuffix(parts[0], "/")
		suffix := strings.TrimPrefix(parts[1], "/")

		// **/*.html - match any path ending with suffix
		if prefix == "" {
			if matched, _ := filepath.Match(suffix, filepath.Base(path)); matched {
				return true
			}
			// Also try matching the suffix against the full path for nested patterns
			if suffix != "" && strings.HasSuffix(path, suffix) {
				return true
			}
		}

		// prefix/**/suffix - match prefix at start, suffix at end
		if prefix != "" && suffix != "" {
			hasPrefix := strings.HasPrefix(path, prefix+"/") || strings.HasPrefix(path, prefix)
			hasSuffix := strings.HasSuffix(path, suffix) ||
				func() bool { m, _ := filepath.Match(suffix, filepath.Base(path)); return m }()
			return hasPrefix && hasSuffix
		}

		// prefix/** - match anything under prefix
		if prefix != "" && suffix == "" {
			return strings.HasPrefix(path, prefix+"/") || path == prefix
		}
	}

	return false
}
