package status

// Status represents a possible service status and associated rendering
// information.
type Status interface {
	// Name returns a human-readable name of the statusm e.g. Healthy or Sick.
	Name() string

	// CSSIdent returns a CSS identifier associated with the status, used for
	// styling.
	CSSIdent() string

	// HTMLChar returns a string to be used as an icon to represent
	// the status, such as a ✓ or ✘.
	HTMLChar() string
}

type status struct {
	name     string
	cssIdent string
	htmlChar string
}

// Name returns a human-readable name of the status.
func (s status) Name() string { return s.name }

// CSSIdent returns the current CSS identifier.
func (s status) CSSIdent() string { return s.cssIdent }

// HTMLChar returns symbolic characters for the status.
func (s status) HTMLChar() string { return s.htmlChar }

// MakeStatus returns a Status instance, populated with the provided
// name, CSS identifier, and symbolic icon(s).
func MakeStatus(name, cssIdent, htmlChar string) Status {
	return status{
		name:     name,
		cssIdent: cssIdent,
		htmlChar: htmlChar,
	}
}

var Healthy Status = MakeStatus("Healthy", "healthy", "✓")
var Sick Status = MakeStatus("Sick", "sick", "✕")
