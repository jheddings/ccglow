package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jheddings/ccnow/internal/statusline"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	var preset, config, format, tee string

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

			output := statusline.Run(statusline.Options{
				Preset: preset,
				Config: config,
				Format: format,
			}, string(stdinBytes))

			if output != "" {
				fmt.Print(output)
			}

			return nil
		},
	}

	root.Flags().StringVar(&preset, "preset", "default", "Use a named preset (default, minimal, full)")
	root.Flags().StringVar(&config, "config", "", "Load JSON config file")
	root.Flags().StringVar(&format, "format", "ansi", "Output format: ansi, plain")
	root.Flags().StringVar(&tee, "tee", "", "Write raw stdin JSON to file before processing")

	root.SetVersionTemplate("{{.Version}}\n")
	root.SilenceUsage = true

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
