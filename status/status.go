package status

type Status interface {
	Name() string
	CSSIdent() string
	HTMLChar() string
}

type status struct {
	name     string
	cssIdent string
	htmlChar string
}

func (s status) Name() string     { return s.name }
func (s status) CSSIdent() string { return s.cssIdent }
func (s status) HTMLChar() string { return s.htmlChar }

func MakeStatus(name, cssIdent, htmlChar string) Status {
	return status{
		name:     name,
		cssIdent: cssIdent,
		htmlChar: htmlChar,
	}
}

var Healthy Status = MakeStatus("Healthy", "healthy", "✓")
var Sick Status = MakeStatus("Sick", "sick", "✕")
