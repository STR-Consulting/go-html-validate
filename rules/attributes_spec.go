package rules

// AutocompleteTokens lists valid autocomplete attribute values.
// Per HTML spec: https://html.spec.whatwg.org/multipage/form-control-infrastructure.html#autofilling-form-controls:-the-autocomplete-attribute
var AutocompleteTokens = map[string]bool{
	// On/off
	"on":  true,
	"off": true,

	// Contact information
	"name":                 true,
	"honorific-prefix":     true,
	"given-name":           true,
	"additional-name":      true,
	"family-name":          true,
	"honorific-suffix":     true,
	"nickname":             true,
	"organization-title":   true,
	"username":             true,
	"new-password":         true,
	"current-password":     true,
	"one-time-code":        true,
	"organization":         true,
	"street-address":       true,
	"address-line1":        true,
	"address-line2":        true,
	"address-line3":        true,
	"address-level4":       true,
	"address-level3":       true,
	"address-level2":       true,
	"address-level1":       true,
	"country":              true,
	"country-name":         true,
	"postal-code":          true,
	"cc-name":              true,
	"cc-given-name":        true,
	"cc-additional-name":   true,
	"cc-family-name":       true,
	"cc-number":            true,
	"cc-exp":               true,
	"cc-exp-month":         true,
	"cc-exp-year":          true,
	"cc-csc":               true,
	"cc-type":              true,
	"transaction-currency": true,
	"transaction-amount":   true,
	"language":             true,
	"bday":                 true,
	"bday-day":             true,
	"bday-month":           true,
	"bday-year":            true,
	"sex":                  true,
	"url":                  true,
	"photo":                true,

	// Telephone
	"tel":              true,
	"tel-country-code": true,
	"tel-national":     true,
	"tel-area-code":    true,
	"tel-local":        true,
	"tel-local-prefix": true,
	"tel-local-suffix": true,
	"tel-extension":    true,

	// Email
	"email": true,
	"impp":  true,

	// Webauthn
	"webauthn": true,
}

// AutocompleteSectionPrefixes are valid prefixes for autocomplete section tokens.
var AutocompleteSectionPrefixes = []string{
	"section-",
	"shipping",
	"billing",
}

// ValidInputTypes lists valid values for <input type="">.
var ValidInputTypes = map[string]bool{
	"button":         true,
	"checkbox":       true,
	"color":          true,
	"date":           true,
	"datetime-local": true,
	"email":          true,
	"file":           true,
	"hidden":         true,
	"image":          true,
	"month":          true,
	"number":         true,
	"password":       true,
	"radio":          true,
	"range":          true,
	"reset":          true,
	"search":         true,
	"submit":         true,
	"tel":            true,
	"text":           true,
	"time":           true,
	"url":            true,
	"week":           true,
}

// ValidButtonTypes lists valid values for <button type="">.
var ValidButtonTypes = map[string]bool{
	"submit": true,
	"reset":  true,
	"button": true,
}

// ValidScriptTypes lists valid MIME types for <script type="">.
// Empty string is also valid (defaults to JavaScript).
var ValidScriptTypes = map[string]bool{
	"":                       true, // Default JavaScript
	"text/javascript":        true,
	"application/javascript": true,
	"text/ecmascript":        true,
	"application/ecmascript": true,
	"module":                 true, // ES modules
	"importmap":              true, // Import maps
	"speculationrules":       true, // Speculation rules
	"text/html":              true, // Data block
	"application/json":       true, // Data block
	"application/ld+json":    true, // JSON-LD
	"text/plain":             true, // Data block
}

// ValidLinkRels lists valid values for <link rel="">.
var ValidLinkRels = map[string]bool{
	"alternate":        true,
	"author":           true,
	"canonical":        true,
	"dns-prefetch":     true,
	"expect":           true,
	"help":             true,
	"icon":             true,
	"license":          true,
	"manifest":         true,
	"modulepreload":    true,
	"next":             true,
	"pingback":         true,
	"preconnect":       true,
	"prefetch":         true,
	"preload":          true,
	"prerender":        true,
	"prev":             true,
	"privacy-policy":   true,
	"search":           true,
	"stylesheet":       true,
	"terms-of-service": true,
}

// ValidAnchorRels lists valid values for <a rel="">.
var ValidAnchorRels = map[string]bool{
	"alternate":        true,
	"author":           true,
	"bookmark":         true,
	"external":         true,
	"help":             true,
	"license":          true,
	"me":               true,
	"next":             true,
	"nofollow":         true,
	"noopener":         true,
	"noreferrer":       true,
	"opener":           true,
	"prev":             true,
	"privacy-policy":   true,
	"search":           true,
	"tag":              true,
	"terms-of-service": true,
}

// ValidFormMethods lists valid values for <form method="">.
var ValidFormMethods = map[string]bool{
	"get":    true,
	"post":   true,
	"dialog": true,
}

// ValidFormEnctypes lists valid values for <form enctype="">.
var ValidFormEnctypes = map[string]bool{
	"application/x-www-form-urlencoded": true,
	"multipart/form-data":               true,
	"text/plain":                        true,
}

// ValidTargets lists valid values for target attribute.
var ValidTargets = map[string]bool{
	"_self":   true,
	"_blank":  true,
	"_parent": true,
	"_top":    true,
}

// ValidScopeValues lists valid values for <th scope="">.
var ValidScopeValues = map[string]bool{
	"row":      true,
	"col":      true,
	"rowgroup": true,
	"colgroup": true,
}

// ValidDirValues lists valid values for dir="" attribute.
var ValidDirValues = map[string]bool{
	"ltr":  true,
	"rtl":  true,
	"auto": true,
}

// ValidLoadingValues lists valid values for loading="" attribute.
var ValidLoadingValues = map[string]bool{
	"eager": true,
	"lazy":  true,
}

// ValidDecodingValues lists valid values for decoding="" attribute.
var ValidDecodingValues = map[string]bool{
	"sync":  true,
	"async": true,
	"auto":  true,
}

// ValidCrossOriginValues lists valid values for crossorigin="" attribute.
var ValidCrossOriginValues = map[string]bool{
	"":                true, // Same as anonymous
	"anonymous":       true,
	"use-credentials": true,
}

// ValidReferrerPolicies lists valid values for referrerpolicy="" attribute.
var ValidReferrerPolicies = map[string]bool{
	"":                                true, // Empty is valid
	"no-referrer":                     true,
	"no-referrer-when-downgrade":      true,
	"origin":                          true,
	"origin-when-cross-origin":        true,
	"same-origin":                     true,
	"strict-origin":                   true,
	"strict-origin-when-cross-origin": true,
	"unsafe-url":                      true,
}

// ValidSandboxTokens lists valid tokens for iframe sandbox="" attribute.
var ValidSandboxTokens = map[string]bool{
	"allow-downloads":                          true,
	"allow-forms":                              true,
	"allow-modals":                             true,
	"allow-orientation-lock":                   true,
	"allow-pointer-lock":                       true,
	"allow-popups":                             true,
	"allow-popups-to-escape-sandbox":           true,
	"allow-presentation":                       true,
	"allow-same-origin":                        true,
	"allow-scripts":                            true,
	"allow-storage-access-by-user-activation":  true,
	"allow-top-navigation":                     true,
	"allow-top-navigation-by-user-activation":  true,
	"allow-top-navigation-to-custom-protocols": true,
}

// InputTypeAttributes maps input types to their supported attributes.
// Attributes not listed here may still be valid globally.
var InputTypeAttributes = map[string]map[string]bool{
	"text": {
		"autocomplete": true, "dirname": true, "list": true,
		"maxlength": true, "minlength": true, "pattern": true,
		"placeholder": true, "readonly": true, "required": true, "size": true,
	},
	"search": {
		"autocomplete": true, "dirname": true, "list": true,
		"maxlength": true, "minlength": true, "pattern": true,
		"placeholder": true, "readonly": true, "required": true, "size": true,
	},
	"url": {
		"autocomplete": true, "list": true,
		"maxlength": true, "minlength": true, "pattern": true,
		"placeholder": true, "readonly": true, "required": true, "size": true,
	},
	"tel": {
		"autocomplete": true, "list": true,
		"maxlength": true, "minlength": true, "pattern": true,
		"placeholder": true, "readonly": true, "required": true, "size": true,
	},
	"email": {
		"autocomplete": true, "list": true,
		"maxlength": true, "minlength": true, "multiple": true, "pattern": true,
		"placeholder": true, "readonly": true, "required": true, "size": true,
	},
	"password": {
		"autocomplete": true,
		"maxlength":    true, "minlength": true, "pattern": true,
		"placeholder": true, "readonly": true, "required": true, "size": true,
	},
	"date": {
		"autocomplete": true, "list": true,
		"max": true, "min": true, "readonly": true, "required": true, "step": true,
	},
	"month": {
		"autocomplete": true, "list": true,
		"max": true, "min": true, "readonly": true, "required": true, "step": true,
	},
	"week": {
		"autocomplete": true, "list": true,
		"max": true, "min": true, "readonly": true, "required": true, "step": true,
	},
	"time": {
		"autocomplete": true, "list": true,
		"max": true, "min": true, "readonly": true, "required": true, "step": true,
	},
	"datetime-local": {
		"autocomplete": true, "list": true,
		"max": true, "min": true, "readonly": true, "required": true, "step": true,
	},
	"number": {
		"autocomplete": true, "list": true,
		"max": true, "min": true, "placeholder": true,
		"readonly": true, "required": true, "step": true,
	},
	"range": {
		"autocomplete": true, "list": true,
		"max": true, "min": true, "step": true,
	},
	"color": {
		"autocomplete": true, "list": true,
	},
	"checkbox": {
		"checked": true, "required": true,
	},
	"radio": {
		"checked": true, "required": true,
	},
	"file": {
		"accept": true, "capture": true, "multiple": true, "required": true,
	},
	"submit": {
		"formaction": true, "formenctype": true, "formmethod": true,
		"formnovalidate": true, "formtarget": true,
	},
	"image": {
		"alt": true, "formaction": true, "formenctype": true, "formmethod": true,
		"formnovalidate": true, "formtarget": true, "height": true, "src": true, "width": true,
	},
	"reset":  {},
	"button": {},
	"hidden": {
		"autocomplete": true,
	},
}

// HTMLCharacterReferences lists valid named character references.
// This is a subset of the most common ones; the full list is very large.
// See: https://html.spec.whatwg.org/multipage/named-characters.html
var HTMLCharacterReferences = map[string]bool{
	// Common entities
	"amp": true, "lt": true, "gt": true, "quot": true, "apos": true,
	"nbsp": true, "copy": true, "reg": true, "trade": true,
	// Math symbols
	"plusmn": true, "times": true, "divide": true, "minus": true,
	"equals": true, "ne": true, "le": true, "ge": true,
	"infin": true, "sum": true, "prod": true, "radic": true,
	// Greek letters (common)
	"alpha": true, "beta": true, "gamma": true, "delta": true,
	"epsilon": true, "zeta": true, "eta": true, "theta": true,
	"iota": true, "kappa": true, "lambda": true, "mu": true,
	"nu": true, "xi": true, "omicron": true, "pi": true,
	"rho": true, "sigma": true, "tau": true, "upsilon": true,
	"phi": true, "chi": true, "psi": true, "omega": true,
	// Currency
	"cent": true, "pound": true, "yen": true, "euro": true, "curren": true,
	// Punctuation
	"mdash": true, "ndash": true, "lsquo": true, "rsquo": true,
	"ldquo": true, "rdquo": true, "hellip": true, "bull": true,
	// Arrows
	"larr": true, "rarr": true, "uarr": true, "darr": true, "harr": true,
	// Other common
	"deg": true, "para": true, "sect": true, "dagger": true, "Dagger": true,
	"laquo": true, "raquo": true, "iexcl": true, "iquest": true,
	"frac14": true, "frac12": true, "frac34": true,
}
