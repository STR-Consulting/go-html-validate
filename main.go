// htmlint is an HTML accessibility linter for Go templates.
//
// Usage:
//
//	htmlint [options] <files or directories>
//
// Options:
//
//	-f, --format     Output format: text, json (default: text)
//	-q, --quiet      Only show errors, not warnings
//	--no-color       Disable colored output
//	--ignore         Glob patterns to ignore (can be repeated)
//	--disable        Disable specific rules (can be repeated)
//	-h, --help       Show help
//
// Examples:
//
//	htmlint web/
//	htmlint -q web/**/*.html
//	htmlint --format=json web/ > lint-results.json
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jschaf/go-html-validate/linter"
	"github.com/jschaf/go-html-validate/reporter"
	"github.com/jschaf/go-html-validate/rules"
)

type stringSlice []string

func (s *stringSlice) String() string { return strings.Join(*s, ",") }
func (s *stringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func main() {
	os.Exit(run())
}

func run() int {
	var (
		format       string
		quiet        bool
		noColor      bool
		ignoreFlags  stringSlice
		disableFlags stringSlice
		showHelp     bool
		listRules    bool
	)

	flag.StringVar(&format, "format", "text", "Output format: text, json")
	flag.StringVar(&format, "f", "text", "Output format (shorthand)")
	flag.BoolVar(&quiet, "quiet", false, "Only show errors")
	flag.BoolVar(&quiet, "q", false, "Only show errors (shorthand)")
	flag.BoolVar(&noColor, "no-color", false, "Disable colored output")
	flag.Var(&ignoreFlags, "ignore", "Glob pattern to ignore")
	flag.Var(&disableFlags, "disable", "Rule to disable")
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showHelp, "h", false, "Show help (shorthand)")
	flag.BoolVar(&listRules, "list-rules", false, "List available rules")

	flag.Usage = usage
	flag.Parse()

	if showHelp {
		usage()
		return 0
	}

	if listRules {
		printRules()
		return 0
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "error: no files or directories specified")
		fmt.Fprintln(os.Stderr, "usage: htmlint [options] <files or directories>")
		return 1
	}

	// Build config
	cfg := linter.DefaultConfig()
	cfg.DisabledRules = disableFlags
	cfg.IgnorePatterns = ignoreFlags
	if quiet {
		cfg.ErrorsOnly()
	}

	// Create linter
	l := linter.New(cfg)

	// Set reporter
	var rep linter.Reporter
	switch format {
	case "json":
		rep = reporter.NewJSON()
	default:
		textRep := reporter.NewText()
		textRep.NoColor = noColor
		rep = textRep
	}
	l.SetReporter(rep)

	// Run linting
	errorCount, err := l.Run(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if errorCount > 0 {
		return 1
	}
	return 0
}

func usage() {
	fmt.Fprintf(os.Stderr, `htmlint - HTML accessibility linter for Go templates

Usage:
  htmlint [options] <files or directories>

Options:
  -f, --format      Output format: text, json (default: text)
  -q, --quiet       Only show errors, not warnings
  --no-color        Disable colored output
  --ignore PATTERN  Glob pattern to ignore (can be repeated)
  --disable RULE    Disable specific rule (can be repeated)
  --list-rules      List available rules
  -h, --help        Show this help

Examples:
  htmlint web/
  htmlint -q web/**/*.html
  htmlint --format=json web/ > lint-results.json
  htmlint --disable=prefer-aria web/
`)
}

func printRules() {
	registry := rules.NewRegistry()
	fmt.Println("Available rules:")
	fmt.Println()
	for _, rule := range registry.All() {
		fmt.Printf("  %-20s %s\n", rule.Name(), rule.Description())
	}
}
