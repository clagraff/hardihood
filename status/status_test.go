package status

import "testing"

const (
	name     = "It's Healthy!"
	cssIdent = "healthy"
	htmlChar = ":)"
)

func TestMakeStatus(t *testing.T) {
	status := MakeStatus(name, cssIdent, htmlChar)

	if status == nil {
		t.Fatal("expected status to not be nil")
	}
}

func TestStatus_Name(t *testing.T) {
	status := MakeStatus(name, cssIdent, htmlChar)

	if status.Name() != name {
		t.Errorf("expected: %s, but got: %s", name, status.Name())
	}
}

func TestStatus_CSSIdent(t *testing.T) {
	status := MakeStatus(name, cssIdent, htmlChar)

	if status.CSSIdent() != cssIdent {
		t.Errorf("expected: %s, but got: %s", cssIdent, status.CSSIdent())
	}
}

func TestStatus_HTMLChar(t *testing.T) {
	status := MakeStatus(name, cssIdent, htmlChar)

	if status.HTMLChar() != htmlChar {
		t.Errorf("expected: %s, but got: %s", htmlChar, status.HTMLChar())
	}
}
