package style

import (
	"testing"

	"github.com/jheddings/ccglow/internal/types"
)

func TestApply_Nil(t *testing.T) {
	if result := Apply("hello", nil); result != "hello" {
		t.Errorf("expected hello, got %s", result)
	}
}

func TestApply_PrefixSuffix(t *testing.T) {
	SetColorLevel(0)
	defer SetColorLevel(1)

	result := Apply("world", &types.StyleAttrs{Prefix: "[", Suffix: "]"})
	if result != "[world]" {
		t.Errorf("expected [world], got %s", result)
	}
}

func TestApply_Bold(t *testing.T) {
	SetColorLevel(1)
	result := Apply("text", &types.StyleAttrs{Bold: true})
	expected := "\x1b[0m\x1b[1mtext\x1b[0m"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestApply_NamedColor(t *testing.T) {
	SetColorLevel(1)
	result := Apply("text", &types.StyleAttrs{Color: "red"})
	expected := "\x1b[0m\x1b[31mtext\x1b[0m"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestApply_256Color(t *testing.T) {
	SetColorLevel(1)
	result := Apply("text", &types.StyleAttrs{Color: "240"})
	expected := "\x1b[0m\x1b[38;5;240mtext\x1b[0m"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestApply_HexColor(t *testing.T) {
	SetColorLevel(1)
	result := Apply("text", &types.StyleAttrs{Color: "#ff0000"})
	expected := "\x1b[0m\x1b[38;2;255;0;0mtext\x1b[0m"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestApply_PlainMode(t *testing.T) {
	SetColorLevel(0)
	defer SetColorLevel(1)

	result := Apply("text", &types.StyleAttrs{Color: "red", Bold: true})
	if result != "text" {
		t.Errorf("expected plain text, got %q", result)
	}
}

func TestApply_PrefixInsideColor(t *testing.T) {
	SetColorLevel(1)
	result := Apply("val", &types.StyleAttrs{Color: "red", Prefix: ">> "})
	expected := "\x1b[0m\x1b[31m>> val\x1b[0m"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
