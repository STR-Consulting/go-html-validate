package rules

import (
	"regexp"
	"strings"

	"github.com/STR-Consulting/go-html-validate/parser"
	"golang.org/x/net/html"
)

// HTMXAttributes validates htmx attribute values.
type HTMXAttributes struct {
	htmxEnabled bool
	htmxVersion string
}

// Configure implements HTMXConfigurable.
func (r *HTMXAttributes) Configure(htmxEnabled bool, htmxVersion string) {
	r.htmxEnabled = htmxEnabled
	r.htmxVersion = htmxVersion
}

// Name returns the rule identifier.
func (r *HTMXAttributes) Name() string { return RuleHTMXAttributes }

// Description returns what this rule checks.
func (r *HTMXAttributes) Description() string {
	return "htmx attribute values must be valid"
}

// Valid hx-swap base values.
var validSwapValues = map[string]bool{
	"innerhtml":   true,
	"outerhtml":   true,
	"beforebegin": true,
	"afterbegin":  true,
	"beforeend":   true,
	"afterend":    true,
	"delete":      true,
	"none":        true,
	// htmx 4 additions
	"textcontent": true,
	"upsert":      true,
}

// Valid hx-swap modifiers.
var validSwapModifiers = map[string]bool{
	"swap":         true,
	"settle":       true,
	"scroll":       true,
	"show":         true,
	"focus-scroll": true,
	"transition":   true,
	"ignoreTitle":  true,
}

// Valid hx-trigger event modifiers.
var validTriggerModifiers = map[string]bool{
	"once":      true,
	"changed":   true,
	"delay":     true,
	"throttle":  true,
	"from":      true,
	"target":    true,
	"consume":   true,
	"queue":     true,
	"root":      true,
	"threshold": true,
}

// Regex patterns for validation.
var (
	// Time format: number followed by s or ms (e.g., "1s", "500ms")
	timePattern = regexp.MustCompile(`^\d+(?:ms|s)$`)
	// Queue modes
	queueModes = map[string]bool{"first": true, "last": true, "all": true, "none": true}
)

// Known DOM events (valid in hx-on:*).
var knownDOMEvents = map[string]bool{
	// Mouse events
	"click": true, "dblclick": true, "mousedown": true, "mouseup": true,
	"mousemove": true, "mouseenter": true, "mouseleave": true, "mouseover": true, "mouseout": true,
	// Keyboard events
	"keydown": true, "keyup": true, "keypress": true,
	// Form events
	"submit": true, "change": true, "input": true, "focus": true, "blur": true, "reset": true,
	"invalid": true, "select": true,
	// Document/Window events
	"load": true, "unload": true, "resize": true, "scroll": true, "error": true,
	"beforeunload": true, "hashchange": true, "popstate": true,
	// Drag events
	"drag": true, "dragstart": true, "dragend": true, "dragover": true,
	"dragenter": true, "dragleave": true, "drop": true,
	// Touch events
	"touchstart": true, "touchend": true, "touchmove": true, "touchcancel": true,
	// Pointer events
	"pointerdown": true, "pointerup": true, "pointermove": true, "pointerenter": true,
	"pointerleave": true, "pointerover": true, "pointerout": true, "pointercancel": true,
	// Animation/transition events
	"animationstart": true, "animationend": true, "animationiteration": true,
	"transitionstart": true, "transitionend": true, "transitionrun": true, "transitioncancel": true,
	// Clipboard events
	"copy": true, "cut": true, "paste": true,
	// Media events
	"play": true, "pause": true, "ended": true, "volumechange": true, "seeking": true,
	"seeked": true, "timeupdate": true, "loadeddata": true, "loadedmetadata": true,
	// Other
	"contextmenu": true, "wheel": true, "compositionstart": true, "compositionend": true,
}

// Known htmx v2 events (htmx:eventName format).
var knownHTMXv2Events = map[string]bool{
	"htmx:abort":                     true,
	"htmx:afterOnLoad":               true,
	"htmx:afterProcessNode":          true,
	"htmx:afterRequest":              true,
	"htmx:afterSettle":               true,
	"htmx:afterSwap":                 true,
	"htmx:beforeCleanupElement":      true,
	"htmx:beforeOnLoad":              true,
	"htmx:beforeProcessNode":         true,
	"htmx:beforeRequest":             true,
	"htmx:beforeSend":                true,
	"htmx:beforeSwap":                true,
	"htmx:beforeTransition":          true,
	"htmx:configRequest":             true,
	"htmx:confirm":                   true,
	"htmx:historyCacheError":         true,
	"htmx:historyCacheHit":           true,
	"htmx:historyCacheMiss":          true,
	"htmx:historyCacheMissLoad":      true,
	"htmx:historyCacheMissLoadError": true,
	"htmx:historyRestore":            true,
	"htmx:beforeHistorySave":         true,
	"htmx:beforeHistoryUpdate":       true,
	"htmx:load":                      true,
	"htmx:noSSESourceError":          true,
	"htmx:oobAfterSwap":              true,
	"htmx:oobBeforeSwap":             true,
	"htmx:oobErrorNoTarget":          true,
	"htmx:onLoadError":               true,
	"htmx:prompt":                    true,
	"htmx:pushedIntoHistory":         true,
	"htmx:replacedInHistory":         true,
	"htmx:responseError":             true,
	"htmx:sendAbort":                 true,
	"htmx:sendError":                 true,
	"htmx:sseError":                  true,
	"htmx:swapError":                 true,
	"htmx:targetError":               true,
	"htmx:timeout":                   true,
	"htmx:trigger":                   true,
	"htmx:validateUrl":               true,
	"htmx:validation:validate":       true,
	"htmx:validation:failed":         true,
	"htmx:validation:halted":         true,
	"htmx:xhr:abort":                 true,
	"htmx:xhr:loadstart":             true,
	"htmx:xhr:loadend":               true,
	"htmx:xhr:progress":              true,
}

// Known htmx v4 event phases and actions (htmx:phase:action format).
var knownHTMXv4Phases = map[string]bool{
	"before":  true,
	"after":   true,
	"error":   true,
	"finally": true,
}

var knownHTMXv4Actions = map[string]bool{
	"request":        true,
	"swap":           true,
	"settle":         true,
	"send":           true,
	"process":        true,
	"cleanup":        true,
	"onLoad":         true,
	"transition":     true,
	"viewTransition": true,
	"history":        true,
	"historyUpdate":  true,
	"historySave":    true,
	"sse":            true,
	"oob":            true,
}

// Check examines the document for invalid htmx attribute values.
func (r *HTMXAttributes) Check(doc *parser.Document) []Result {
	if !r.htmxEnabled {
		return nil
	}

	var results []Result

	doc.Walk(func(n *parser.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}

		for _, attr := range n.Attr {
			attrName := strings.ToLower(attr.Key)

			if !strings.HasPrefix(attrName, "hx-") {
				continue
			}

			var validationResults []Result

			switch {
			case attrName == "hx-swap":
				validationResults = r.validateSwap(doc.Filename, n, attr.Val)
			case attrName == "hx-trigger":
				validationResults = r.validateTrigger(doc.Filename, n, attr.Val)
			case attrName == "hx-target":
				validationResults = r.validateTarget(doc.Filename, n, attr.Val)
			case strings.HasPrefix(attrName, "hx-on:") || strings.HasPrefix(attrName, "hx-on-"):
				validationResults = r.validateHxOn(doc.Filename, n, attr.Key)
			}

			results = append(results, validationResults...)
		}

		return true
	})

	return results
}

// validateSwap checks hx-swap attribute values.
func (r *HTMXAttributes) validateSwap(filename string, n *parser.Node, value string) []Result {
	if value == "" {
		return nil // Empty is valid (uses default)
	}

	// Skip template expressions (both raw and preprocessed)
	if strings.Contains(value, "{{") || strings.Contains(value, "TMPL") {
		return nil
	}

	var results []Result
	value = strings.TrimSpace(value)

	// Split on whitespace to separate base value from modifiers
	parts := strings.Fields(value)
	if len(parts) == 0 {
		return nil
	}

	// First part is the swap strategy
	baseValue := strings.ToLower(parts[0])

	// Check for htmx 4 only values when using v2
	if r.htmxVersion != "4" && (baseValue == "textcontent" || baseValue == "upsert") {
		results = append(results, Result{
			Rule:     RuleHTMXAttributes,
			Message:  "hx-swap value '" + baseValue + "' is only available in htmx 4",
			Filename: filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Warning,
		})
		return results
	}

	if !validSwapValues[baseValue] {
		results = append(results, Result{
			Rule:     RuleHTMXAttributes,
			Message:  "invalid hx-swap value '" + parts[0] + "'",
			Filename: filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Error,
		})
		return results
	}

	// Validate modifiers
	for i := 1; i < len(parts); i++ {
		modifier := parts[i]
		colonIdx := strings.Index(modifier, ":")
		if colonIdx == -1 {
			results = append(results, Result{
				Rule:     RuleHTMXAttributes,
				Message:  "invalid hx-swap modifier '" + modifier + "' (missing colon)",
				Filename: filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
			continue
		}

		modName := strings.ToLower(modifier[:colonIdx])
		modValue := modifier[colonIdx+1:]

		if !validSwapModifiers[modName] {
			results = append(results, Result{
				Rule:     RuleHTMXAttributes,
				Message:  "unknown hx-swap modifier '" + modName + "'",
				Filename: filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
			continue
		}

		// Validate modifier values
		switch modName {
		case "swap", "settle":
			if !timePattern.MatchString(modValue) {
				results = append(results, Result{
					Rule:     RuleHTMXAttributes,
					Message:  "hx-swap " + modName + " modifier requires a time value (e.g., '1s', '500ms')",
					Filename: filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
		case "scroll", "show":
			validPositions := map[string]bool{"top": true, "bottom": true}
			if !validPositions[strings.ToLower(modValue)] && !strings.HasPrefix(modValue, "#") {
				results = append(results, Result{
					Rule:     RuleHTMXAttributes,
					Message:  "hx-swap " + modName + " modifier value should be 'top', 'bottom', or a selector",
					Filename: filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Warning,
				})
			}
		case "focus-scroll":
			if modValue != "true" && modValue != "false" {
				results = append(results, Result{
					Rule:     RuleHTMXAttributes,
					Message:  "hx-swap focus-scroll modifier should be 'true' or 'false'",
					Filename: filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
		case "transition":
			if modValue != "true" && modValue != "false" {
				results = append(results, Result{
					Rule:     RuleHTMXAttributes,
					Message:  "hx-swap transition modifier should be 'true' or 'false'",
					Filename: filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
		}
	}

	return results
}

// validateTrigger checks hx-trigger attribute values.
func (r *HTMXAttributes) validateTrigger(filename string, n *parser.Node, value string) []Result {
	if value == "" {
		return nil
	}

	// Skip template expressions (both raw and preprocessed)
	if strings.Contains(value, "{{") || strings.Contains(value, "TMPL") {
		return nil
	}

	var results []Result
	value = strings.TrimSpace(value)

	// Handle multiple triggers separated by commas
	triggers := strings.Split(value, ",")

	for _, trigger := range triggers {
		trigger = strings.TrimSpace(trigger)
		if trigger == "" {
			continue
		}

		triggerResults := r.validateSingleTrigger(filename, n, trigger)
		results = append(results, triggerResults...)
	}

	return results
}

// validateSingleTrigger validates a single trigger specification.
func (r *HTMXAttributes) validateSingleTrigger(filename string, n *parser.Node, trigger string) []Result {
	var results []Result

	// Split into parts by whitespace
	parts := strings.Fields(trigger)
	if len(parts) == 0 {
		return nil
	}

	// First part is the event name
	eventName := strings.ToLower(parts[0])

	// Handle "every Xs" special syntax
	if eventName == "every" {
		if len(parts) < 2 {
			results = append(results, Result{
				Rule:     RuleHTMXAttributes,
				Message:  "hx-trigger 'every' requires a time value (e.g., 'every 1s')",
				Filename: filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
			return results
		}
		if !timePattern.MatchString(parts[1]) {
			results = append(results, Result{
				Rule:     RuleHTMXAttributes,
				Message:  "hx-trigger 'every' requires a valid time value (e.g., '1s', '500ms')",
				Filename: filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Error,
			})
		}
		return results
	}

	// Handle "intersect" with optional modifiers
	if eventName == "intersect" {
		// intersect is valid, check modifiers
		for i := 1; i < len(parts); i++ {
			mod := parts[i]
			if strings.HasPrefix(mod, "root:") || strings.HasPrefix(mod, "threshold:") {
				continue
			}
			if mod == "once" {
				continue
			}
			// Unknown modifier for intersect
			results = append(results, Result{
				Rule:     RuleHTMXAttributes,
				Message:  "unknown intersect modifier '" + mod + "'",
				Filename: filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
		}
		return results
	}

	// Validate remaining modifiers
	for i := 1; i < len(parts); i++ {
		modifier := parts[i]

		// Handle bracket syntax for filters like [ctrlKey]
		if strings.HasPrefix(modifier, "[") && strings.HasSuffix(modifier, "]") {
			continue // Filter expressions are valid
		}

		// Check for colon-style modifiers
		colonIdx := strings.Index(modifier, ":")
		if colonIdx == -1 {
			// Standalone modifier (e.g., "once", "changed", "consume")
			modName := strings.ToLower(modifier)
			if !validTriggerModifiers[modName] {
				results = append(results, Result{
					Rule:     RuleHTMXAttributes,
					Message:  "unknown hx-trigger modifier '" + modifier + "'",
					Filename: filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Warning,
				})
			}
			continue
		}

		modName := strings.ToLower(modifier[:colonIdx])
		modValue := modifier[colonIdx+1:]

		if !validTriggerModifiers[modName] {
			results = append(results, Result{
				Rule:     RuleHTMXAttributes,
				Message:  "unknown hx-trigger modifier '" + modName + "'",
				Filename: filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			})
			continue
		}

		// Validate modifier values
		switch modName {
		case "delay", "throttle":
			if !timePattern.MatchString(modValue) {
				results = append(results, Result{
					Rule:     RuleHTMXAttributes,
					Message:  "hx-trigger " + modName + " requires a time value (e.g., '1s', '500ms')",
					Filename: filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
		case "queue":
			if !queueModes[strings.ToLower(modValue)] {
				results = append(results, Result{
					Rule:     RuleHTMXAttributes,
					Message:  "hx-trigger queue mode should be 'first', 'last', 'all', or 'none'",
					Filename: filename,
					Line:     n.Line,
					Col:      n.Col,
					Severity: Error,
				})
			}
		}
	}

	return results
}

// validateTarget checks hx-target attribute values.
func (r *HTMXAttributes) validateTarget(filename string, n *parser.Node, value string) []Result {
	if value == "" {
		return nil
	}

	// Skip template expressions (both raw and preprocessed)
	if strings.Contains(value, "{{") || strings.Contains(value, "TMPL") {
		return nil
	}

	var results []Result
	value = strings.TrimSpace(value)

	// Valid special values
	specialValues := map[string]bool{
		"this":     true,
		"next":     true,
		"previous": true,
		"body":     true,
	}

	// Handle single-word values
	if !strings.Contains(value, " ") {
		lower := strings.ToLower(value)
		if specialValues[lower] {
			return nil
		}
		// Assume it's a CSS selector - basic validation
		if strings.HasPrefix(value, "#") || strings.HasPrefix(value, ".") || isValidElement(lower) {
			return nil
		}
		// Could be a valid selector, warn but don't error
		return nil
	}

	// Handle "closest X", "find X", "next X", "previous X"
	parts := strings.SplitN(value, " ", 2)
	keyword := strings.ToLower(parts[0])

	validKeywords := map[string]bool{
		"closest":  true,
		"find":     true,
		"next":     true,
		"previous": true,
	}

	if !validKeywords[keyword] && !specialValues[keyword] {
		results = append(results, Result{
			Rule:     RuleHTMXAttributes,
			Message:  "invalid hx-target keyword '" + parts[0] + "'; expected 'this', 'closest', 'find', 'next', 'previous', or a CSS selector",
			Filename: filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Warning,
		})
	}

	return results
}

// isValidElement checks if the name is a known HTML element.
func isValidElement(name string) bool {
	// Common HTML elements
	elements := map[string]bool{
		"a": true, "abbr": true, "address": true, "area": true, "article": true, "aside": true, "audio": true,
		"b": true, "base": true, "bdi": true, "bdo": true, "blockquote": true, "body": true, "br": true, "button": true,
		"canvas": true, "caption": true, "cite": true, "code": true, "col": true, "colgroup": true,
		"data": true, "datalist": true, "dd": true, "del": true, "details": true, "dfn": true, "dialog": true, "div": true, "dl": true, "dt": true,
		"em": true, "embed": true,
		"fieldset": true, "figcaption": true, "figure": true, "footer": true, "form": true,
		"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true, "head": true, "header": true, "hgroup": true, "hr": true, "html": true,
		"i": true, "iframe": true, "img": true, "input": true, "ins": true,
		"kbd": true, "label": true, "legend": true, "li": true, "link": true,
		"main": true, "map": true, "mark": true, "menu": true, "meta": true, "meter": true,
		"nav": true, "noscript": true,
		"object": true, "ol": true, "optgroup": true, "option": true, "output": true,
		"p": true, "picture": true, "pre": true, "progress": true,
		"q": true, "rp": true, "rt": true, "ruby": true,
		"s": true, "samp": true, "script": true, "search": true, "section": true, "select": true, "slot": true, "small": true, "source": true, "span": true, "strong": true, "style": true, "sub": true, "summary": true, "sup": true,
		"table": true, "tbody": true, "td": true, "template": true, "textarea": true, "tfoot": true, "th": true, "thead": true, "time": true, "title": true, "tr": true, "track": true,
		"u": true, "ul": true, "var": true, "video": true, "wbr": true,
	}
	return elements[name]
}

// validateHxOn checks hx-on:* attribute event names.
// Validates that the event name is a known DOM event or htmx event.
func (r *HTMXAttributes) validateHxOn(filename string, n *parser.Node, attrKey string) []Result {
	// Extract event name from attribute key
	// hx-on:click -> click
	// hx-on:htmx:afterRequest -> htmx:afterRequest
	// hx-on-click -> click (deprecated but still supported)
	// hx-on::htmx:afterRequest -> htmx:afterRequest (v4 shorthand)
	var eventName string
	if strings.HasPrefix(attrKey, "hx-on:") {
		eventName = strings.TrimPrefix(attrKey, "hx-on:")
	} else if strings.HasPrefix(attrKey, "hx-on-") {
		eventName = strings.TrimPrefix(attrKey, "hx-on-")
	}

	if eventName == "" {
		return []Result{{
			Rule:     RuleHTMXAttributes,
			Message:  "hx-on:* requires an event name",
			Filename: filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Error,
		}}
	}

	// Handle htmx 4 shorthand (hx-on::event for htmx: prefixed events)
	// e.g., hx-on::after-request is shorthand for hx-on:htmx:after:request
	if strings.HasPrefix(eventName, ":") {
		eventName = "htmx" + eventName
	}

	// Check for DOM events (lowercase comparison)
	lowerEvent := strings.ToLower(eventName)
	if knownDOMEvents[lowerEvent] {
		return nil // Valid DOM event
	}

	// Check for htmx events
	if strings.HasPrefix(eventName, "htmx:") || strings.HasPrefix(lowerEvent, "htmx:") {
		return r.validateHTMXEvent(filename, n, eventName)
	}

	// Unknown event - could be a custom event, warn
	return []Result{{
		Rule:     RuleHTMXAttributes,
		Message:  "unknown event '" + eventName + "' in hx-on:*; if this is a custom event, ignore this warning",
		Filename: filename,
		Line:     n.Line,
		Col:      n.Col,
		Severity: Warning,
	}}
}

// validateHTMXEvent validates an htmx event name against known events.
func (r *HTMXAttributes) validateHTMXEvent(filename string, n *parser.Node, eventName string) []Result {
	// htmx v4 uses htmx:phase:action format
	// htmx v2 uses htmx:eventName format (camelCase)

	if r.htmxVersion == "4" {
		return r.validateHTMXv4Event(filename, n, eventName)
	}

	// htmx v2 validation
	// Note: HTML attributes are lowercased by the parser, so we need case-insensitive matching.
	// We accept the event if it matches any known event case-insensitively.
	for validEvent := range knownHTMXv2Events {
		if strings.EqualFold(eventName, validEvent) {
			return nil // Valid htmx v2 event (case-insensitive match)
		}
	}

	return []Result{{
		Rule:     RuleHTMXAttributes,
		Message:  "unknown htmx event '" + eventName + "'",
		Filename: filename,
		Line:     n.Line,
		Col:      n.Col,
		Severity: Warning,
	}}
}

// validateHTMXv4Event validates an htmx 4 event name (htmx:phase:action format).
func (r *HTMXAttributes) validateHTMXv4Event(filename string, n *parser.Node, eventName string) []Result {
	// htmx 4 uses htmx:phase:action[:sub-action] format
	// e.g., htmx:after:request, htmx:before:swap, htmx:error

	// Remove htmx: prefix
	remainder := strings.TrimPrefix(eventName, "htmx:")
	parts := strings.SplitN(remainder, ":", 2)

	if len(parts) == 0 || parts[0] == "" {
		return []Result{{
			Rule:     RuleHTMXAttributes,
			Message:  "invalid htmx event format '" + eventName + "'",
			Filename: filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Error,
		}}
	}

	// First part should be a phase
	phase := strings.ToLower(parts[0])
	if !knownHTMXv4Phases[phase] {
		// Could be a standalone event like htmx:load, htmx:abort
		standaloneEvents := map[string]bool{
			"load": true, "abort": true, "trigger": true, "confirm": true, "prompt": true,
		}
		if standaloneEvents[phase] {
			return nil // Valid standalone event
		}

		return []Result{{
			Rule:     RuleHTMXAttributes,
			Message:  "unknown htmx 4 event phase '" + parts[0] + "' in '" + eventName + "'",
			Filename: filename,
			Line:     n.Line,
			Col:      n.Col,
			Severity: Warning,
		}}
	}

	// If there's a second part, it should be an action
	if len(parts) > 1 && parts[1] != "" {
		// Extract just the action (before any sub-action)
		actionParts := strings.SplitN(parts[1], ":", 2)
		action := strings.ToLower(actionParts[0])

		if !knownHTMXv4Actions[action] {
			return []Result{{
				Rule:     RuleHTMXAttributes,
				Message:  "unknown htmx 4 event action '" + actionParts[0] + "' in '" + eventName + "'",
				Filename: filename,
				Line:     n.Line,
				Col:      n.Col,
				Severity: Warning,
			}}
		}
	}

	return nil
}
