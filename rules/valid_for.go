package rules

import (
	"strings"

	"github.com/toba/go-html-validate/parser"
	"golang.org/x/net/html"
)

// labelableElements are elements that can be associated with a <label>.
var labelableElements = map[string]bool{
	"input":    true,
	"textarea": true,
	"select":   true,
	"button":   true,
	"output":   true,
	"meter":    true,
	"progress": true,
}

// ValidFor checks that label for attributes reference labelable elements.
type ValidFor struct{}

func (r *ValidFor) Name() string { return RuleValidFor }

func (r *ValidFor) Description() string {
	return "label for attribute must reference a labelable element"
}

func (r *ValidFor) Check(doc *parser.Document) []Result {
	// Build a map of id -> node for lookup
	idMap := make(map[string]*parser.Node)
	doc.Walk(func(n *parser.Node) bool {
		if n.Type == html.ElementNode && n.HasAttr("id") {
			id := n.GetAttr("id")
			if id != "" {
				idMap[id] = n
			}
		}
		return true
	})

	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode || !n.IsElement("label") {
			return true
		}

		forAttr := n.GetAttr("for")
		if forAttr == "" {
			return true
		}

		// Skip template expressions
		if IsTemplateExpr(forAttr) {
			return true
		}

		target, exists := idMap[forAttr]
		if !exists {
			// Target not in this document; could be in another template fragment
			return true
		}

		tag := strings.ToLower(target.Data)
		if !labelableElements[tag] {
			results = append(results, Result{
				Rule:     r.Name(),
				Message:  "label for attribute references non-labelable element <" + tag + ">",
				Filename: doc.Filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
			return true
		}

		// input type="hidden" is not labelable
		if tag == "input" {
			inputType := strings.ToLower(target.GetAttr("type"))
			if inputType == "hidden" {
				results = append(results, Result{
					Rule:     r.Name(),
					Message:  "label for attribute references hidden input",
					Filename: doc.Filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
		}

		return true
	})

	return results
}
