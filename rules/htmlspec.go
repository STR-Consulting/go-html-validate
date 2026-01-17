package rules

// ContentModel represents HTML5 content model categories.
type ContentModel int

const (
	ModelNone        ContentModel = 0
	ModelFlow        ContentModel = 1 << iota // Most elements in body
	ModelPhrasing                             // Text-level elements
	ModelInteractive                          // User-interactive elements
	ModelHeading                              // h1-h6
	ModelSectioning                           // article, aside, nav, section
	ModelEmbedded                             // img, video, audio, etc.
	ModelMetadata                             // head elements
	ModelTransparent                          // Inherits parent's content model
)

// ElementSpec defines an HTML element's constraints per HTML5 spec.
type ElementSpec struct {
	// ContentModel indicates what kind of content the element represents.
	ContentModel ContentModel
	// PermittedContent specifies allowed child element categories or specific tags.
	// Empty means element can contain any flow content.
	PermittedContent []string
	// ForbiddenContent specifies elements that cannot be descendants.
	ForbiddenContent []string
	// PermittedParents specifies allowed parent elements.
	// Empty means any element accepting the content model.
	PermittedParents []string
	// RequiredParent specifies a parent that must exist (not necessarily direct).
	RequiredParent string
	// RequiredChildren specifies children that must exist.
	RequiredChildren []string
	// RequiredAttributes lists attributes that must be present.
	RequiredAttributes []string
	// VoidElement is true for elements that cannot have children (br, img, etc.).
	VoidElement bool
	// Deprecated: indicates the element should not be used.
	Deprecated bool
	// DeprecatedMessage provides guidance on what to use instead.
	DeprecatedMessage string
}

// AttrType indicates the type of attribute value validation.
type AttrType int

const (
	AttrTypeString    AttrType = iota // Any string value
	AttrTypeEnum                      // Must be one of AllowedValues
	AttrTypeBoolean                   // Boolean attribute (presence-based)
	AttrTypePattern                   // Must match Pattern regex
	AttrTypeID                        // Valid ID value
	AttrTypeIDRef                     // Reference to an ID
	AttrTypeIDRefList                 // Space-separated list of ID refs
	AttrTypeURL                       // Valid URL
	AttrTypeInteger                   // Integer value
	AttrTypePositive                  // Positive integer
)

// AttrSpec defines an attribute's constraints.
type AttrSpec struct {
	// Type indicates how the value should be validated.
	Type AttrType
	// AllowedValues lists valid values for AttrTypeEnum.
	AllowedValues []string
	// Pattern is a regex for AttrTypePattern.
	Pattern string
	// Required indicates the attribute must be present.
	Required bool
	// Deprecated: indicates the attribute should not be used.
	Deprecated bool
	// DeprecatedMessage provides guidance on what to use instead.
	DeprecatedMessage string
	// ValidFor lists elements this attribute is valid on (empty = global).
	ValidFor []string
}
