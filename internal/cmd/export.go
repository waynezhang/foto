package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynezhang/foto/internal/export"
)

var ExportCmd = func() *cobra.Command {
	var outputPath string
	var minimize bool

	fn := func(cmd *cobra.Command, args []string) {
		export.Export(outputPath, minimize)
	}

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export sites",
		Run:   fn,
	}
	cmd.Flags().StringVarP(&outputPath, "output", "o", "dist", "Output directory")
	cmd.Flags().BoolVarP(&minimize, "minimize", "m", false, "Wether minimize output files(css, html, js supported) or not")

	return cmd
}()
