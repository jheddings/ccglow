package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jheddings/ccnow/internal/config"
	"github.com/jheddings/ccnow/internal/preset"
	"github.com/jheddings/ccnow/internal/provider"
	"github.com/jheddings/ccnow/internal/render"
	"github.com/jheddings/ccnow/internal/segment"
	"github.com/jheddings/ccnow/internal/session"
	"github.com/jheddings/ccnow/internal/style"
	"github.com/jheddings/ccnow/internal/types"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	var presetName, configPath, format, tee string

	root := &cobra.Command{
		Use:     "ccnow",
		Short:   "Composable statusline for Claude Code",
		Long:    "Reads session JSON from stdin, outputs styled statusline to stdout.",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			stdinBytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				stdinBytes = []byte{}
			}

			if tee != "" {
				if err := os.WriteFile(tee, stdinBytes, 0644); err != nil {
					fmt.Fprintf(os.Stderr, "ccnow: failed to write tee file: %v\n", err)
				}
			}

			output := run(presetName, configPath, format, string(stdinBytes))
			if output != "" {
				fmt.Print(output)
			}

			return nil
		},
	}

	root.Flags().StringVar(&presetName, "preset", "default", "Use a named preset (default, minimal, full)")
	root.Flags().StringVar(&configPath, "config", "", "Load JSON config file")
	root.Flags().StringVar(&format, "format", "ansi", "Output format: ansi, plain")
	root.Flags().StringVar(&tee, "tee", "", "Write raw stdin JSON to file before processing")

	root.SetVersionTemplate("{{.Version}}\n")
	root.SilenceUsage = true

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(presetName, configPath, format, stdin string) string {
	sess := session.Parse(stdin)
	if sess == nil {
		return ""
	}

	if format == "plain" {
		style.SetColorLevel(0)
	} else {
		style.SetColorLevel(1)
	}
	defer style.SetColorLevel(1)

	segments := segment.NewRegistry()
	segment.RegisterBuiltin(segments)

	providers := provider.NewRegistry()
	provider.RegisterBuiltin(providers)

	tree := resolveTree(presetName, configPath)

	providerNames := render.CollectProviderNames(tree)
	providerData := render.ResolveProviders(providerNames, providers.All(), sess)

	return render.Tree(tree, segments, sess, providerData)
}

func resolveTree(presetName, configPath string) []types.SegmentNode {
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ccnow: failed to load config: %v\n", err)
		} else {
			tree, err := config.Parse(data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ccnow: failed to parse config: %v\n", err)
			} else if len(tree) > 0 {
				return tree
			}
		}
	}

	if tree := preset.Get(presetName); tree != nil {
		return tree
	}

	return preset.Get("default")
}
