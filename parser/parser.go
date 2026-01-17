package parser

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Document represents a parsed HTML document with source location tracking.
type Document struct {
	Root     *Node
	Filename string
	// IsTemplateFragment indicates file starts with {{define - a Go template partial
	IsTemplateFragment bool
	// sourceMap for converting positions back to original
	sourceMap *SourceMap
}

// Node wraps html.Node with source location and traversal helpers.
type Node struct {
	*html.Node
	Line     int
	Col      int
	Parent   *Node
	Children []*Node
}

// HasAttr checks if the node has an attribute with the given name.
func (n *Node) HasAttr(name string) bool {
	if n.Node == nil {
		return false
	}
	for _, attr := range n.Attr {
		if strings.EqualFold(attr.Key, name) {
			return true
		}
	}
	return false
}

// GetAttr returns the value of an attribute, or empty string if not found.
func (n *Node) GetAttr(name string) string {
	if n.Node == nil {
		return ""
	}
	for _, attr := range n.Attr {
		if strings.EqualFold(attr.Key, name) {
			return attr.Val
		}
	}
	return ""
}

// TextContent returns the combined text content of the node and descendants.
func (n *Node) TextContent() string {
	if n.Node == nil {
		return ""
	}
	var buf strings.Builder
	n.collectText(&buf)
	return strings.TrimSpace(buf.String())
}

func (n *Node) collectText(buf *strings.Builder) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for _, child := range n.Children {
		child.collectText(buf)
	}
}

// IsElement returns true if this is an element node with the given tag name.
func (n *Node) IsElement(tag string) bool {
	return n.Type == html.ElementNode && strings.EqualFold(n.Data, tag)
}

// WalkFunc is called for each node during tree traversal.
// Return false to stop traversal.
type WalkFunc func(*Node) bool

// Walk traverses the document tree, calling fn for each node.
func (d *Document) Walk(fn WalkFunc) {
	if d.Root != nil {
		d.Root.walk(fn)
	}
}

func (n *Node) walk(fn WalkFunc) bool {
	if !fn(n) {
		return false
	}
	for _, child := range n.Children {
		if !child.walk(fn) {
			return false
		}
	}
	return true
}

// Parse parses HTML content and returns a Document with line tracking.
func Parse(filename string, content []byte) (*Document, error) {
	// Preprocess to handle Go template syntax
	prep := NewPreprocessor()
	processed, sourceMap, err := prep.Process(content)
	if err != nil {
		return nil, err
	}

	// Parse the processed HTML
	root, err := html.Parse(bytes.NewReader(processed))
	if err != nil {
		return nil, err
	}

	doc := &Document{
		Filename:  filename,
		sourceMap: sourceMap,
	}

	// Build our node tree
	doc.Root = buildNodeTree(root, nil)

	return doc, nil
}

// ParseFragment parses an HTML fragment (like a template partial).
func ParseFragment(filename string, content []byte) (*Document, error) {
	// Detect Go template fragments (files starting with {{define)
	isTemplateFragment := bytes.HasPrefix(bytes.TrimSpace(content), []byte("{{define"))

	// Preprocess to handle Go template syntax
	prep := NewPreprocessor()
	processed, sourceMap, err := prep.Process(content)
	if err != nil {
		return nil, err
	}

	// Create a context element for fragment parsing
	context := &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	}

	// Parse as fragment
	nodes, err := html.ParseFragment(bytes.NewReader(processed), context)
	if err != nil {
		return nil, err
	}

	doc := &Document{
		Filename:           filename,
		IsTemplateFragment: isTemplateFragment,
		sourceMap:          sourceMap,
	}

	// Build a synthetic root containing all fragments
	syntheticRoot := &Node{
		Node: &html.Node{Type: html.DocumentNode},
		Line: 1,
		Col:  1,
	}

	for _, n := range nodes {
		child := buildNodeTree(n, syntheticRoot)
		syntheticRoot.Children = append(syntheticRoot.Children, child)
	}

	doc.Root = syntheticRoot
	return doc, nil
}

// buildNodeTree converts html.Node tree to our Node tree.
// Note: golang.org/x/net/html doesn't provide source positions,
// so all nodes have line=1, col=1 as placeholders.
func buildNodeTree(n *html.Node, parent *Node) *Node {
	line, col := 1, 1
	if parent != nil {
		line = parent.Line
		col = parent.Col
	}

	node := &Node{
		Node:   n,
		Parent: parent,
		Line:   line,
		Col:    col,
	}

	// Process children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		child := buildNodeTree(c, node)
		node.Children = append(node.Children, child)
	}

	return node
}

// ParseReader parses HTML from an io.Reader.
func ParseReader(filename string, r io.Reader) (*Document, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return Parse(filename, content)
}
