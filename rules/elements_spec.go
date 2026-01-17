package rules

// VoidElements are HTML elements that cannot have content (self-closing).
// Per HTML5 spec: https://html.spec.whatwg.org/multipage/syntax.html#void-elements
var VoidElements = map[string]bool{
	"area":   true,
	"base":   true,
	"br":     true,
	"col":    true,
	"embed":  true,
	"hr":     true,
	"img":    true,
	"input":  true,
	"link":   true,
	"meta":   true,
	"param":  true, // deprecated in HTML5 but still void
	"source": true,
	"track":  true,
	"wbr":    true,
}

// ValidElements is the set of standard HTML5 element names.
// Custom elements (containing hyphen) are handled separately.
var ValidElements = map[string]bool{
	// Document metadata
	"html": true, "head": true, "title": true, "base": true, "link": true,
	"meta": true, "style": true,
	// Sectioning root
	"body": true,
	// Content sectioning
	"article": true, "section": true, "nav": true, "aside": true,
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"hgroup": true, "header": true, "footer": true, "address": true,
	// Text content
	"p": true, "hr": true, "pre": true, "blockquote": true,
	"ol": true, "ul": true, "menu": true, "li": true,
	"dl": true, "dt": true, "dd": true,
	"figure": true, "figcaption": true, "main": true, "search": true, "div": true,
	// Inline text
	"a": true, "em": true, "strong": true, "small": true, "s": true,
	"cite": true, "q": true, "dfn": true, "abbr": true, "ruby": true,
	"rt": true, "rp": true, "data": true, "time": true, "code": true,
	"var": true, "samp": true, "kbd": true, "sub": true, "sup": true,
	"i": true, "b": true, "u": true, "mark": true, "bdi": true, "bdo": true,
	"span": true, "br": true, "wbr": true,
	// Edits
	"ins": true, "del": true,
	// Embedded content
	"picture": true, "source": true, "img": true, "iframe": true,
	"embed": true, "object": true, "param": true, "video": true, "audio": true,
	"track": true, "map": true, "area": true,
	// SVG and MathML
	"svg": true, "math": true,
	// Tabular data
	"table": true, "caption": true, "colgroup": true, "col": true,
	"tbody": true, "thead": true, "tfoot": true, "tr": true, "td": true, "th": true,
	// Forms
	"form": true, "label": true, "input": true, "button": true, "select": true,
	"datalist": true, "optgroup": true, "option": true, "textarea": true,
	"output": true, "progress": true, "meter": true, "fieldset": true, "legend": true,
	// Interactive elements
	"details": true, "summary": true, "dialog": true,
	// Scripting
	"script": true, "noscript": true, "template": true, "slot": true, "canvas": true,
}

// ElementSpecs provides detailed constraints for elements that need special validation.
// Not all elements need entries here - only those with specific requirements.
var ElementSpecs = map[string]ElementSpec{
	// Document structure
	"html": {
		ContentModel:     ModelNone,
		RequiredChildren: []string{"head", "body"},
	},
	"head": {
		ContentModel:     ModelMetadata,
		RequiredChildren: []string{"title"}, // Required unless iframe srcdoc
		PermittedParents: []string{"html"},
	},
	"body": {
		ContentModel:     ModelFlow,
		PermittedParents: []string{"html"},
	},
	"title": {
		ContentModel:     ModelNone, // Text only
		PermittedParents: []string{"head"},
	},

	// Lists require specific children
	"ul": {
		ContentModel:     ModelFlow,
		PermittedContent: []string{"li", "script", "template"},
	},
	"ol": {
		ContentModel:     ModelFlow,
		PermittedContent: []string{"li", "script", "template"},
	},
	"li": {
		ContentModel:   ModelFlow,
		RequiredParent: "ul", // or ol, menu - checked in rule
	},
	"dl": {
		ContentModel:     ModelFlow,
		PermittedContent: []string{"dt", "dd", "div", "script", "template"},
	},
	"dt": {
		ContentModel:   ModelFlow,
		RequiredParent: "dl",
	},
	"dd": {
		ContentModel:   ModelFlow,
		RequiredParent: "dl",
	},

	// Tables
	"table": {
		ContentModel:     ModelFlow,
		PermittedContent: []string{"caption", "colgroup", "thead", "tbody", "tfoot", "tr", "script", "template"},
	},
	"caption": {
		ContentModel:     ModelFlow,
		PermittedParents: []string{"table"},
	},
	"colgroup": {
		ContentModel:     ModelNone,
		PermittedContent: []string{"col", "template"},
		PermittedParents: []string{"table"},
	},
	"col": {
		VoidElement:      true,
		PermittedParents: []string{"colgroup"},
	},
	"thead": {
		ContentModel:     ModelNone,
		PermittedContent: []string{"tr", "script", "template"},
		PermittedParents: []string{"table"},
	},
	"tbody": {
		ContentModel:     ModelNone,
		PermittedContent: []string{"tr", "script", "template"},
		PermittedParents: []string{"table"},
	},
	"tfoot": {
		ContentModel:     ModelNone,
		PermittedContent: []string{"tr", "script", "template"},
		PermittedParents: []string{"table"},
	},
	"tr": {
		ContentModel:     ModelNone,
		PermittedContent: []string{"td", "th", "script", "template"},
		RequiredParent:   "table", // or thead/tbody/tfoot
	},
	"td": {
		ContentModel:     ModelFlow,
		PermittedParents: []string{"tr"},
	},
	"th": {
		ContentModel:     ModelFlow,
		PermittedParents: []string{"tr"},
	},

	// Forms
	"form": {
		ContentModel:     ModelFlow,
		ForbiddenContent: []string{"form"}, // Forms cannot be nested
	},
	"fieldset": {
		ContentModel: ModelFlow,
		// Note: legend must be first child if present, but flow content follows
	},
	"legend": {
		ContentModel:     ModelPhrasing,
		PermittedParents: []string{"fieldset"},
	},
	"label": {
		ContentModel:     ModelPhrasing,
		ForbiddenContent: []string{"label"}, // Labels cannot be nested
	},
	"select": {
		ContentModel:     ModelNone,
		PermittedContent: []string{"option", "optgroup", "hr", "script", "template"},
	},
	"optgroup": {
		ContentModel:     ModelNone,
		PermittedContent: []string{"option", "script", "template"},
		PermittedParents: []string{"select"},
	},
	"option": {
		ContentModel:   ModelNone, // Text only
		RequiredParent: "select",  // or optgroup/datalist
	},
	"datalist": {
		ContentModel:     ModelPhrasing,
		PermittedContent: []string{"option", "script", "template"},
	},

	// Media and embeds
	"picture": {
		ContentModel:     ModelNone,
		PermittedContent: []string{"source", "img", "script", "template"},
	},
	"video": {
		ContentModel:     ModelTransparent,
		PermittedContent: []string{"source", "track"},
	},
	"audio": {
		ContentModel:     ModelTransparent,
		PermittedContent: []string{"source", "track"},
	},
	"source": {
		VoidElement:      true,
		PermittedParents: []string{"audio", "video", "picture"},
	},
	"track": {
		VoidElement:      true,
		PermittedParents: []string{"audio", "video"},
	},
	"map": {
		ContentModel:       ModelTransparent,
		RequiredAttributes: []string{"name"},
	},
	"area": {
		VoidElement:      true,
		PermittedParents: []string{"map"},
	},

	// Interactive
	"details": {
		ContentModel: ModelFlow,
		// Note: summary must be first child if present, but flow content follows
	},
	"summary": {
		ContentModel:     ModelPhrasing,
		PermittedParents: []string{"details"},
	},

	// Ruby annotations
	"ruby": {
		ContentModel:     ModelPhrasing,
		PermittedContent: []string{"rt", "rp"},
	},
	"rt": {
		ContentModel:     ModelPhrasing,
		PermittedParents: []string{"ruby"},
	},
	"rp": {
		ContentModel:     ModelNone, // Text only
		PermittedParents: []string{"ruby"},
	},

	// Void elements
	"br":    {VoidElement: true},
	"hr":    {VoidElement: true},
	"img":   {VoidElement: true, RequiredAttributes: []string{"src", "alt"}},
	"input": {VoidElement: true},
	"meta":  {VoidElement: true},
	"link":  {VoidElement: true},
	"base":  {VoidElement: true},
	"embed": {VoidElement: true},
	"wbr":   {VoidElement: true},
	"param": {VoidElement: true, Deprecated: true, DeprecatedMessage: "use object data attribute"},
}

// RequiredAncestors maps elements to their required ancestor elements.
// The element must have at least one of these as an ancestor.
var RequiredAncestors = map[string][]string{
	"li":         {"ul", "ol", "menu"},
	"dt":         {"dl"},
	"dd":         {"dl"},
	"td":         {"table"},
	"th":         {"table"},
	"tr":         {"table"},
	"thead":      {"table"},
	"tbody":      {"table"},
	"tfoot":      {"table"},
	"caption":    {"table"},
	"colgroup":   {"table"},
	"col":        {"colgroup"},
	"legend":     {"fieldset"},
	"optgroup":   {"select"},
	"option":     {"select", "optgroup", "datalist"},
	"source":     {"audio", "video", "picture"},
	"track":      {"audio", "video"},
	"area":       {"map"},
	"figcaption": {"figure"},
	"summary":    {"details"},
	"rt":         {"ruby"},
	"rp":         {"ruby"},
}

// UniqueElements are elements that should appear at most once in their context.
var UniqueElements = map[string]string{
	"title":   "head",     // One title per head
	"base":    "head",     // One base per head
	"main":    "body",     // One main per document (already have rule)
	"caption": "table",    // One caption per table
	"thead":   "table",    // One thead per table
	"tfoot":   "table",    // One tfoot per table
	"summary": "details",  // One summary per details (first child)
	"legend":  "fieldset", // One legend per fieldset (first child)
}

// ForbiddenDescendants maps elements to descendants they cannot contain.
var ForbiddenDescendants = map[string][]string{
	"a":        {"a"},                                          // Links cannot contain links
	"button":   {"a", "button", "input", "select", "textarea"}, // Interactive content
	"label":    {"label"},                                      // Labels cannot be nested
	"form":     {"form"},                                       // Forms cannot be nested
	"progress": {"progress"},
	"meter":    {"meter"},
	"dfn":      {"dfn"},
	"abbr":     {"abbr"},
	"header":   {"header", "footer"},
	"footer":   {"header", "footer"},
	"address":  {"address", "header", "footer"},
}
