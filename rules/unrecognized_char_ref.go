package rules

import (
	"bytes"
	"html"
	"regexp"

	"github.com/toba/go-html-validate/parser"
)

// namedCharRefPattern matches named character references like &amp; or &aacute;
// Excludes numeric references (&#123; &#x1F;).
var namedCharRefPattern = regexp.MustCompile(`&([a-zA-Z][a-zA-Z0-9]*);`)

// UnrecognizedCharRef checks that named character references are valid HTML5 entities.
type UnrecognizedCharRef struct{}

func (r *UnrecognizedCharRef) Name() string { return RuleUnrecognizedCharRef }

func (r *UnrecognizedCharRef) Description() string {
	return "character references must be valid HTML5 entities"
}

// Check implements Rule but returns nil - this rule uses CheckRaw instead.
func (r *UnrecognizedCharRef) Check(_ *parser.Document) []Result {
	return nil
}

// CheckRaw examines the raw content for unrecognized named character references.
func (r *UnrecognizedCharRef) CheckRaw(filename string, content []byte) []Result {
	var results []Result

	lines := bytes.Split(content, []byte("\n"))
	for lineNum, line := range lines {
		// Skip lines inside Go template expressions
		lineStr := string(line)

		matches := namedCharRefPattern.FindAllStringSubmatchIndex(lineStr, -1)
		for _, match := range matches {
			// match[0]:match[1] is the full match &name;
			// match[2]:match[3] is the captured name
			fullRef := lineStr[match[0]:match[1]]
			name := lineStr[match[2]:match[3]]

			// Skip if inside a Go template expression
			if isInsideTemplateExpr(lineStr, match[0]) {
				continue
			}

			// Use html.UnescapeString to check validity:
			// if the result equals the input, the entity is unrecognized
			unescaped := html.UnescapeString(fullRef)
			if unescaped == fullRef {
				results = append(results, Result{
					Rule:     r.Name(),
					Message:  "unrecognized character reference &" + name + ";",
					Filename: filename,
					Line:     lineNum + 1,
					Col:      match[0] + 1,
					Severity: Warning,
				})
			}
		}
	}

	return results
}

// isInsideTemplateExpr checks if a position is within a {{ ... }} template expression.
func isInsideTemplateExpr(line string, pos int) bool {
	// Find all {{ and }} positions and check if pos falls inside one
	depth := 0
	for i := 0; i < len(line)-1 && i < pos; i++ {
		if line[i] == '{' && line[i+1] == '{' {
			depth++
			i++ // skip second brace
		} else if line[i] == '}' && line[i+1] == '}' {
			depth--
			i++
		}
	}
	return depth > 0
}
