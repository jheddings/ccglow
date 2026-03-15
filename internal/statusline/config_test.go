package statusline

import "testing"

func TestParseConfig_Valid(t *testing.T) {
	input := `{
		"segments": [
			{"segment": "pwd.name", "style": {"color": "red"}},
			{"segment": "git.branch"}
		]
	}`

	nodes, err := ParseConfig([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}
	if nodes[0].Type != "pwd.name" {
		t.Errorf("expected pwd.name, got %s", nodes[0].Type)
	}
	if nodes[0].Provider != "pwd" {
		t.Errorf("expected inferred provider pwd, got %s", nodes[0].Provider)
	}
	if nodes[1].Provider != "git" {
		t.Errorf("expected inferred provider git, got %s", nodes[1].Provider)
	}
}

func TestParseConfig_WithChildren(t *testing.T) {
	input := `{
		"segments": [
			{
				"segment": "group",
				"children": [
					{"segment": "git.branch"},
					{"segment": "git.insertions"}
				]
			}
		]
	}`

	nodes, err := ParseConfig([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
	if len(nodes[0].Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(nodes[0].Children))
	}
	if nodes[0].Children[0].Provider != "git" {
		t.Errorf("expected inferred provider git on child, got %s", nodes[0].Children[0].Provider)
	}
}

func TestParseConfig_LiteralNoProvider(t *testing.T) {
	input := `{"segments": [{"segment": "literal", "props": {"text": "hi"}}]}`

	nodes, err := ParseConfig([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if nodes[0].Provider != "" {
		t.Errorf("expected no provider for literal, got %s", nodes[0].Provider)
	}
}

func TestParseConfig_InvalidJSON(t *testing.T) {
	_, err := ParseConfig([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
