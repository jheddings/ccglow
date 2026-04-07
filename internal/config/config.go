package config

import (
	"encoding/json"

	"github.com/jheddings/ccglow/internal/types"
)

type configFile struct {
	Segments    []json.RawMessage `json:"segments"`
	Width       int               `json:"width,omitempty"`
	WidthOffset int               `json:"width_offset,omitempty"`
}

// Options carries top-level render configuration parsed from a config file.
type Options struct {
	Width       int
	WidthOffset int
}

// ParseOptions extracts top-level render options from a config file. Returns
// a zero-value Options on parse failure.
func ParseOptions(data []byte) Options {
	var cfg configFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Options{}
	}
	return Options{Width: cfg.Width, WidthOffset: cfg.WidthOffset}
}

// Parse parses a JSON config file into a segment tree.
func Parse(data []byte) ([]types.SegmentNode, error) {
	var cfg configFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	var nodes []types.SegmentNode
	for _, raw := range cfg.Segments {
		var node types.SegmentNode
		if err := json.Unmarshal(raw, &node); err != nil {
			continue
		}
		if node.Expr == "" && node.Value == nil && node.Command == "" && len(node.Children) == 0 && !node.Flex {
			continue
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}
