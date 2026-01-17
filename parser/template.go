// Package parser provides HTML parsing with Go template syntax support.
package parser

import (
	"bytes"
	"regexp"
)

// templatePattern matches Go template syntax: {{ ... }}
// Uses non-greedy matching to handle nested braces correctly.
var templatePattern = regexp.MustCompile(`\{\{[\s\S]*?\}\}`)

// ifElseEndPattern matches {{if ...}}...{{else}}...{{end}} blocks including nested.
// Captures: full match with if-branch content kept, else-branch removed.
var ifElseEndPattern = regexp.MustCompile(`(?s)\{\{-?\s*if\s[^}]*\}\}(.*?)\{\{-?\s*else\s*-?\}\}.*?\{\{-?\s*end\s*-?\}\}`)

// ifEndPattern matches {{if ...}}...{{end}} blocks without else.
var ifEndPattern = regexp.MustCompile(`(?s)\{\{-?\s*if\s[^}]*\}\}(.*?)\{\{-?\s*end\s*-?\}\}`)

// SourceMap tracks the mapping between processed and original source positions.
// Used to report errors at their original line/column locations.
type SourceMap struct {
	// Original source content
	Original []byte
	// Processed content with templates replaced
	Processed []byte
}

// OriginalPosition converts a position in processed content to original position.
func (sm *SourceMap) OriginalPosition(line, col int) (origLine, origCol int) {
	// For now, line numbers are preserved since we only replace inline content
	// Column offsets would need more sophisticated tracking for accuracy
	return line, col
}

// Preprocessor handles Go template syntax in HTML files.
type Preprocessor struct{}

// NewPreprocessor creates a new template preprocessor.
func NewPreprocessor() *Preprocessor {
	return &Preprocessor{}
}

// Process replaces Go template syntax with placeholders to produce valid HTML.
// Returns the processed content and a source map for error location recovery.
//
// Replacement strategies:
//   - {{ .Field }} in text content → empty string (preserves structure)
//   - {{ .Field }} in attribute values → "tmpl" (keeps attribute valid)
//   - {{if}}...{{else}}...{{end}} blocks → content of if-branch kept only
//   - {{if}}...{{end}} blocks → content kept
//   - {{range}}...{{end}} → single iteration content
//   - {{template "name"}} → empty (included template not available)
func (p *Preprocessor) Process(input []byte) ([]byte, *SourceMap, error) {
	sm := &SourceMap{
		Original: input,
	}

	// First, handle {{if}}...{{else}}...{{end}} blocks - keep only if-branch
	processed := ifElseEndPattern.ReplaceAll(input, []byte("$1"))

	// Then handle {{if}}...{{end}} without else - keep content
	processed = ifEndPattern.ReplaceAll(processed, []byte("$1"))

	// Replace remaining template expressions with appropriate placeholders
	processed = templatePattern.ReplaceAllFunc(processed, func(match []byte) []byte {
		return p.replaceTemplate(match)
	})

	sm.Processed = processed
	return processed, sm, nil
}

// replaceTemplate determines the appropriate replacement for a template expression.
func (p *Preprocessor) replaceTemplate(match []byte) []byte {
	content := bytes.TrimSpace(match[2 : len(match)-2]) // Remove {{ and }}

	// Handle different template constructs
	switch {
	case bytes.HasPrefix(content, []byte("/*")):
		// Template comment: {{/* comment */}} → empty
		return nil

	case bytes.HasPrefix(content, []byte("if ")),
		bytes.HasPrefix(content, []byte("if(")),
		bytes.Equal(content, []byte("else")),
		bytes.HasPrefix(content, []byte("else if")):
		// Control flow: keep structure, remove directive
		return nil

	case bytes.HasPrefix(content, []byte("end")):
		// End of block
		return nil

	case bytes.HasPrefix(content, []byte("range ")):
		// Range loop start
		return nil

	case bytes.HasPrefix(content, []byte("template ")),
		bytes.HasPrefix(content, []byte("block ")):
		// Template inclusion
		return nil

	case bytes.HasPrefix(content, []byte("define ")):
		// Template definition
		return nil

	case bytes.HasPrefix(content, []byte("with ")):
		// With block
		return nil

	case bytes.HasPrefix(content, []byte("-")):
		// Whitespace trimming prefix
		return nil

	default:
		// Variable or function call: {{ .Field }} or {{ func .arg }}
		// Replace with placeholder that works in most contexts
		return []byte("TMPL")
	}
}

// ProcessFile reads a file and processes its template content.
func (p *Preprocessor) ProcessFile(content []byte) ([]byte, *SourceMap, error) {
	return p.Process(content)
}
