package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynezhang/foto/internal/log"
	"github.com/waynezhang/foto/internal/utils"
)

func Execute() {
  var verbose bool
  var rootCmd = &cobra.Command{
    Use: "foto",
    Short: "Yet another publishing tool for photographers",
    Run: func(cmd *cobra.Command, args []string) {
      _ = cmd.Help()
    },
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
      log.SetVerbose(verbose)
    },
  }

  rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

  rootCmd.AddCommand(CreateCmd)
  rootCmd.AddCommand(PreviewCmd)
  rootCmd.AddCommand(ExportCmd)
  rootCmd.AddCommand(VersionCmd)

  err := rootCmd.Execute()
  utils.CheckFatalError(err, "Failed to execute command.")
}

