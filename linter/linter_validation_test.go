package linter_test

import (
	"strings"
	"testing"

	"github.com/toba/go-html-validate/linter"
	"github.com/toba/go-html-validate/rules"
)

func TestLintContent_ValidID(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
		severity rules.Severity
	}{
		{
			name:     "empty id",
			html:     `<div id="">Content</div>`,
			wantRule: "valid-id",
			severity: rules.Error,
		},
		{
			name:     "id with space",
			html:     `<div id="foo bar">Content</div>`,
			wantRule: "valid-id",
			severity: rules.Error,
		},
		{
			name:     "id starts with digit",
			html:     `<div id="123abc">Content</div>`,
			wantRule: "valid-id",
			severity: rules.Warning,
		},
		{
			name: "valid id with hyphen",
			html: `<div id="my-id">Content</div>`,
		},
		{
			name: "valid id with underscore",
			html: `<div id="my_id">Content</div>`,
		},
		{
			name: "no id attribute",
			html: `<div>Content</div>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			if tt.wantRule == "" {
				for _, r := range results {
					if r.Rule == "valid-id" {
						t.Errorf("expected no valid-id results, got %v", results)
					}
				}
				return
			}
			found := false
			for _, r := range results {
				if r.Rule == tt.wantRule {
					found = true
					if r.Severity != tt.severity {
						t.Errorf("expected severity %v, got %v", tt.severity, r.Severity)
					}
					break
				}
			}
			if !found {
				t.Errorf("expected rule %q in results, got %v", tt.wantRule, results)
			}
		})
	}
}

func TestLintContent_RequireLang(t *testing.T) {
	// Note: LintContent uses ParseFragment which doesn't preserve html element structure.
	// The require-lang rule is tested via LintFile in integration tests for full documents.
	// These tests verify the rule doesn't flag fragments.
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "fragment without html element (no flag)",
			html: `<div>Content</div>`,
		},
		{
			name: "fragment with main content (no flag)",
			html: `<main>Content</main>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleRequireLang, tt.wantRule)
		})
	}
}

func TestLintContent_ElementName(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "valid element",
			html: `<div>content</div>`,
		},
		{
			name: "valid custom element",
			html: `<my-component>content</my-component>`,
		},
		{
			name:     "unknown element",
			html:     `<foobar>content</foobar>`,
			wantRule: rules.RuleElementName,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleElementName, tt.wantRule)
		})
	}
}

func TestLintContent_AttributeAllowedValues(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "valid input type",
			html: `<input type="text">`,
		},
		{
			name:     "invalid input type",
			html:     `<input type="foobar">`,
			wantRule: rules.RuleAttributeAllowedValues,
		},
		{
			name: "valid button type",
			html: `<button type="submit">Click</button>`,
		},
		{
			name:     "invalid button type",
			html:     `<button type="invalid">Click</button>`,
			wantRule: rules.RuleAttributeAllowedValues,
		},
		{
			name: "valid form method",
			html: `<form method="post"></form>`,
		},
		{
			name:     "invalid form method",
			html:     `<form method="put"></form>`,
			wantRule: rules.RuleAttributeAllowedValues,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleAttributeAllowedValues, tt.wantRule)
		})
	}
}

func TestLintContent_NoMissingReferences(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "for references existing id",
			html: `<label for="name">Name</label><input id="name">`,
		},
		{
			name:     "for references non-existent id",
			html:     `<label for="missing">Name</label><input id="name">`,
			wantRule: rules.RuleNoMissingReferences,
		},
		{
			name: "aria-labelledby references existing id",
			html: `<span id="label">Label</span><input aria-labelledby="label">`,
		},
		{
			name:     "aria-labelledby references non-existent id",
			html:     `<input aria-labelledby="missing">`,
			wantRule: rules.RuleNoMissingReferences,
		},
		{
			name: "template expression in for (skip)",
			html: `<label for="TMPL">Name</label>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleNoMissingReferences, tt.wantRule)
		})
	}
}

func TestLintContent_FormDupName(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "unique names",
			html: `<form><input name="a"><input name="b"></form>`,
		},
		{
			name:     "duplicate names",
			html:     `<form><input name="a"><input name="a"></form>`,
			wantRule: rules.RuleFormDupName,
		},
		{
			name: "radio buttons can share names",
			html: `<form><input type="radio" name="choice"><input type="radio" name="choice"></form>`,
		},
		{
			name: "checkboxes can share names",
			html: `<form><input type="checkbox" name="opts"><input type="checkbox" name="opts"></form>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleFormDupName, tt.wantRule)
		})
	}
}

func TestLintContent_MapIDName(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "map with matching id and name",
			html: `<map id="nav" name="nav"></map>`,
		},
		{
			name:     "map with mismatched id and name",
			html:     `<map id="nav1" name="nav2"></map>`,
			wantRule: rules.RuleMapIDName,
		},
		{
			name: "map with name only",
			html: `<map name="nav"></map>`,
		},
		{
			name:     "map without name",
			html:     `<map id="nav"></map>`,
			wantRule: rules.RuleMapIDName,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleMapIDName, tt.wantRule)
		})
	}
}

func TestLintContent_NoDupClass(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "no class attr",
			html: `<p>text</p>`,
		},
		{
			name: "unique classes",
			html: `<p class="foo bar">text</p>`,
		},
		{
			name: "other attrs ok",
			html: `<p attr="foo bar foo">text</p>`,
		},
		{
			name:     "duplicate class",
			html:     `<p class="foo bar foo">text</p>`,
			wantRule: rules.RuleNoDupClass,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleNoDupClass, tt.wantRule)
		})
	}
}

func TestLintContent_AllowedLinks(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "valid http link",
			html: `<a href="https://example.com">Link</a>`,
		},
		{
			name: "valid relative link",
			html: `<a href="/page">Link</a>`,
		},
		{
			name: "valid anchor link",
			html: `<a href="#section">Link</a>`,
		},
		{
			name:     "javascript protocol",
			html:     `<a href="javascript:alert(1)">Link</a>`,
			wantRule: rules.RuleAllowedLinks,
		},
		{
			name:     "data protocol",
			html:     `<a href="data:text/html,<h1>Hi</h1>">Link</a>`,
			wantRule: rules.RuleAllowedLinks,
		},
		{
			name: "template expression (skip)",
			html: `<a href="TMPL">Link</a>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleAllowedLinks, tt.wantRule)
		})
	}
}

func TestLintContent_LongTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "title over 70 chars",
			html:     `<html><head><title>This is a very long title that exceeds the recommended seventy character limit for SEO</title></head></html>`,
			wantRule: "long-title",
		},
		{
			name: "title under 70 chars",
			html: `<html><head><title>Short Title</title></head></html>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleLongTitle, tt.wantRule)
		})
	}
}

func TestLintContent_NoInlineStyle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name:     "element with inline style",
			html:     `<div style="color: red;">Red text</div>`,
			wantRule: "no-inline-style",
		},
		{
			name: "element without style",
			html: `<div class="red">Red text</div>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleNoInlineStyle, tt.wantRule)
		})
	}
}

func TestLintContent_InputAttributes_HTMX(t *testing.T) {
	tests := []struct {
		name        string
		html        string
		htmxEnabled bool
		htmxVersion string
		wantRule    string
		wantMessage string
	}{
		{
			name:        "htmx attr without config warns",
			html:        `<input type="text" hx-get="/api">`,
			htmxEnabled: false,
			wantRule:    rules.RuleInputAttributes,
			wantMessage: "htmx attribute 'hx-get' used but htmx not enabled",
		},
		{
			name:        "htmx attr with config enabled passes",
			html:        `<input type="text" hx-get="/api">`,
			htmxEnabled: true,
			htmxVersion: "2",
		},
		{
			name:        "htmx v4-only attr with v2 config warns",
			html:        `<input type="text" hx-optimistic>`,
			htmxEnabled: true,
			htmxVersion: "2",
			wantRule:    rules.RuleInputAttributes,
			wantMessage: "only available in htmx 4",
		},
		{
			name:        "htmx v4-only attr with v4 config passes",
			html:        `<input type="text" hx-optimistic>`,
			htmxEnabled: true,
			htmxVersion: "4",
		},
		{
			name:        "htmx deprecated attr with v4 config warns",
			html:        `<input type="text" hx-vars="foo:bar">`,
			htmxEnabled: true,
			htmxVersion: "4",
			wantRule:    rules.RuleInputAttributes,
			wantMessage: "deprecated in htmx 4",
		},
		{
			name:        "htmx deprecated attr with v2 config passes",
			html:        `<input type="text" hx-vars="foo:bar">`,
			htmxEnabled: true,
			htmxVersion: "2",
		},
		{
			name:        "hx-on event handler passes",
			html:        `<input type="text" hx-on:click="alert()">`,
			htmxEnabled: true,
			htmxVersion: "2",
		},
		{
			name:        "multiple htmx attrs with config enabled passes",
			html:        `<input type="text" hx-get="/api" hx-trigger="change" hx-target="#result">`,
			htmxEnabled: true,
			htmxVersion: "2",
		},
		{
			name:        "unknown htmx attr warns",
			html:        `<input type="text" hx-invalid-attr>`,
			htmxEnabled: true,
			htmxVersion: "2",
			wantRule:    rules.RuleInputAttributes,
			wantMessage: "unknown htmx attribute",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := linter.DefaultConfig()
			cfg.Frameworks.HTMX = tt.htmxEnabled
			cfg.Frameworks.HTMXVersion = tt.htmxVersion
			l := linter.New(cfg)

			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}

			if tt.wantRule == "" {
				// Should have no input-attributes violations
				for _, r := range results {
					if r.Rule == rules.RuleInputAttributes {
						t.Errorf("expected no %s rule, but got: %s", rules.RuleInputAttributes, r.Message)
					}
				}
				return
			}

			// Should have the expected violation
			found := false
			for _, r := range results {
				if r.Rule == tt.wantRule {
					found = true
					if tt.wantMessage != "" && !messageContains(r.Message, tt.wantMessage) {
						t.Errorf("expected message containing %q, got %q", tt.wantMessage, r.Message)
					}
					break
				}
			}
			if !found {
				t.Errorf("expected rule %q in results, got %v", tt.wantRule, results)
			}
		})
	}
}

// messageContains checks if message contains the expected substring.
func messageContains(message, expected string) bool {
	return strings.Contains(message, expected)
}

func TestLintContent_ValidFor(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "label for references input",
			html: `<label for="name">Name</label><input id="name">`,
		},
		{
			name: "label for references textarea",
			html: `<label for="bio">Bio</label><textarea id="bio"></textarea>`,
		},
		{
			name: "label for references select",
			html: `<label for="color">Color</label><select id="color"><option>Red</option></select>`,
		},
		{
			name: "label for references button",
			html: `<label for="btn">Action</label><button id="btn">Click</button>`,
		},
		{
			name: "label for references output",
			html: `<label for="result">Result</label><output id="result">42</output>`,
		},
		{
			name: "label for references meter",
			html: `<label for="fuel">Fuel</label><meter id="fuel" value="0.5">50%</meter>`,
		},
		{
			name: "label for references progress",
			html: `<label for="prog">Progress</label><progress id="prog" value="50" max="100">50%</progress>`,
		},
		{
			name:     "label for references div",
			html:     `<label for="foo">Name</label><div id="foo">text</div>`,
			wantRule: rules.RuleValidFor,
		},
		{
			name:     "label for references p",
			html:     `<label for="foo">Name</label><p id="foo">text</p>`,
			wantRule: rules.RuleValidFor,
		},
		{
			name:     "label for references span",
			html:     `<label for="foo">Name</label><span id="foo">text</span>`,
			wantRule: rules.RuleValidFor,
		},
		{
			name:     "label for references hidden input",
			html:     `<label for="tok">Token</label><input type="hidden" id="tok">`,
			wantRule: rules.RuleValidFor,
		},
		{
			name: "label for references missing id (skip)",
			html: `<label for="missing">Name</label>`,
		},
		{
			name: "label without for (skip)",
			html: `<label>Name <input></label>`,
		},
		{
			name: "label for with template expression (skip)",
			html: `<label for="TMPL">Name</label>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleValidFor, tt.wantRule)
		})
	}
}

func TestLintContent_UnrecognizedCharRef(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantRule string
	}{
		{
			name: "valid entity amp",
			html: `<p>Tom &amp; Jerry</p>`,
		},
		{
			name: "valid entity lt",
			html: `<p>1 &lt; 2</p>`,
		},
		{
			name: "valid entity copy",
			html: `<p>&copy; 2024</p>`,
		},
		{
			name: "valid entity aacute",
			html: `<p>&aacute;</p>`,
		},
		{
			name: "valid entity nbsp",
			html: `<p>hello&nbsp;world</p>`,
		},
		{
			name:     "invalid entity foobar",
			html:     `<p>&foobar;</p>`,
			wantRule: rules.RuleUnrecognizedCharRef,
		},
		{
			name:     "invalid entity bloop",
			html:     `<p>&bloop;</p>`,
			wantRule: rules.RuleUnrecognizedCharRef,
		},
		{
			name: "numeric decimal entity (skip)",
			html: `<p>&#8212;</p>`,
		},
		{
			name: "numeric hex entity (skip)",
			html: `<p>&#x2014;</p>`,
		},
	}

	l := linter.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := l.LintContent("test.html", []byte(tt.html))
			if err != nil {
				t.Fatalf("LintContent() error = %v", err)
			}
			checkRule(t, results, rules.RuleUnrecognizedCharRef, tt.wantRule)
		})
	}
}
