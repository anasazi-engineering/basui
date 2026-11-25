package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
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
	init := &cobra.Command{
		Use:   "init",
		Short: "Initialize a basui project in the current directory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := initSite("."); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Initialized basui project in the current directory")
			return nil
		},
	}
	root.AddCommand(build, init)
	return root
}
