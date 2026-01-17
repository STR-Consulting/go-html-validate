package rules

import "strings"

// HeadingTags identifies heading elements (h1-h6).
var HeadingTags = map[string]bool{
	"h1": true, "h2": true, "h3": true,
	"h4": true, "h5": true, "h6": true,
}

// HeadingRank returns the numeric rank of a heading (1-6), or 0 if not a heading.
func HeadingRank(tagName string) int {
	switch strings.ToLower(tagName) {
	case "h1":
		return 1
	case "h2":
		return 2
	case "h3":
		return 3
	case "h4":
		return 4
	case "h5":
		return 5
	case "h6":
		return 6
	default:
		return 0
	}
}

// LabelableElements are HTML elements that can be associated with a <label>.
var LabelableElements = map[string]bool{
	"button":   true,
	"input":    true, // except type="hidden"
	"meter":    true,
	"output":   true,
	"progress": true,
	"select":   true,
	"textarea": true,
}

// LandmarkElements maps HTML elements to their implicit ARIA landmark role.
var LandmarkElements = map[string]string{
	"aside":   "complementary",
	"footer":  "contentinfo",
	"form":    "form",
	"header":  "banner",
	"main":    "main",
	"nav":     "navigation",
	"section": "region",
}

// AriaLabelableElements are elements that can have aria-label/aria-labelledby.
var AriaLabelableElements = map[string]bool{
	// Landmark elements
	"main": true, "nav": true, "aside": true, "header": true,
	"footer": true, "section": true, "article": true, "form": true,
	// Interactive elements
	"a": true, "button": true, "input": true, "select": true, "textarea": true,
	// Other labelable elements
	"table": true, "dialog": true, "iframe": true, "img": true,
	"figure": true, "summary": true, "details": true,
	"meter": true, "output": true, "progress": true,
	// SVG elements
	"svg": true,
}

// ImplicitRoles maps HTML elements to their implicit ARIA roles.
var ImplicitRoles = map[string]string{
	"a":        "link", // when href present
	"article":  "article",
	"aside":    "complementary",
	"button":   "button",
	"dialog":   "dialog",
	"form":     "form",
	"h1":       "heading",
	"h2":       "heading",
	"h3":       "heading",
	"h4":       "heading",
	"h5":       "heading",
	"h6":       "heading",
	"header":   "banner",
	"footer":   "contentinfo",
	"img":      "img",
	"input":    "", // varies by type
	"li":       "listitem",
	"main":     "main",
	"nav":      "navigation",
	"ol":       "list",
	"option":   "option",
	"progress": "progressbar",
	"section":  "region",
	"select":   "combobox",
	"table":    "table",
	"tbody":    "rowgroup",
	"td":       "cell",
	"textarea": "textbox",
	"tfoot":    "rowgroup",
	"th":       "columnheader",
	"thead":    "rowgroup",
	"tr":       "row",
	"ul":       "list",
}

// InputTypeRoles maps input types to their implicit ARIA roles.
var InputTypeRoles = map[string]string{
	"button":   "button",
	"checkbox": "checkbox",
	"email":    "textbox",
	"image":    "button",
	"number":   "spinbutton",
	"radio":    "radio",
	"range":    "slider",
	"reset":    "button",
	"search":   "searchbox",
	"submit":   "button",
	"tel":      "textbox",
	"text":     "textbox",
	"url":      "textbox",
}

// FormControlElements are HTML elements that act as form controls.
var FormControlElements = map[string]bool{
	"input":    true,
	"select":   true,
	"textarea": true,
	"button":   true,
	"output":   true,
}

// ButtonLikeInputTypes are input types that behave like buttons.
var ButtonLikeInputTypes = map[string]bool{
	"button": true,
	"submit": true,
	"reset":  true,
	"image":  true,
}

// AbstractRoles are ARIA roles that must not be used directly.
// See: https://www.w3.org/TR/wai-aria/#abstract_roles
var AbstractRoles = map[string]bool{
	"command":     true,
	"composite":   true,
	"input":       true,
	"landmark":    true,
	"range":       true,
	"roletype":    true,
	"section":     true,
	"sectionhead": true,
	"select":      true,
	"structure":   true,
	"widget":      true,
	"window":      true,
}
