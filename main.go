package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{Use: "basui", Short: "Compile HTML-first operator screens", SilenceUsage: true}
	build := &cobra.Command{
		Use:   "build [site-directory]",
		Short: "Compile layouts, pages and assets into the site's dist directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			site := "."
			if len(args) == 1 {
				site = args[0]
			}
			count, err := buildSite(site)
			if err == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "Built %d page(s) in %s/dist\n", count, site)
			}
			return err
		},
	}
	root.AddCommand(build)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
