package statusline

import "testing"

func TestApplyStyle_Nil(t *testing.T) {
	result := ApplyStyle("hello", nil)
	if result != "hello" {
		t.Errorf("expected hello, got %s", result)
	}
}

func TestApplyStyle_PrefixSuffix(t *testing.T) {
	SetColorLevel(0)
	defer SetColorLevel(1)

	result := ApplyStyle("world", &StyleAttrs{Prefix: "[", Suffix: "]"})
	if result != "[world]" {
		t.Errorf("expected [world], got %s", result)
	}
}

func TestApplyStyle_Bold(t *testing.T) {
	SetColorLevel(1)
	result := ApplyStyle("text", &StyleAttrs{Bold: true})
	expected := "\x1b[0m\x1b[1mtext\x1b[0m"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestApplyStyle_NamedColor(t *testing.T) {
	SetColorLevel(1)
	result := ApplyStyle("text", &StyleAttrs{Color: "red"})
	expected := "\x1b[0m\x1b[31mtext\x1b[0m"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestApplyStyle_256Color(t *testing.T) {
	SetColorLevel(1)
	result := ApplyStyle("text", &StyleAttrs{Color: "240"})
	expected := "\x1b[0m\x1b[38;5;240mtext\x1b[0m"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestApplyStyle_HexColor(t *testing.T) {
	SetColorLevel(1)
	result := ApplyStyle("text", &StyleAttrs{Color: "#ff0000"})
	expected := "\x1b[0m\x1b[38;2;255;0;0mtext\x1b[0m"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestApplyStyle_PlainMode(t *testing.T) {
	SetColorLevel(0)
	defer SetColorLevel(1)

	result := ApplyStyle("text", &StyleAttrs{Color: "red", Bold: true})
	if result != "text" {
		t.Errorf("expected plain text, got %q", result)
	}
}

func TestApplyStyle_PrefixOutsideColor(t *testing.T) {
	SetColorLevel(1)
	result := ApplyStyle("val", &StyleAttrs{Color: "red", Prefix: ">> "})
	expected := ">> \x1b[0m\x1b[31mval\x1b[0m"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
