package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/waynezhang/foto/internal/utils"
)

func Execute() {
	var verbose bool
	var rootCmd = &cobra.Command{
		Use:   "foto",
		Short: "Yet another publishing tool for photographers",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			if verbose {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			} else {
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(ClearCacheCmd)
	rootCmd.AddCommand(CreateCmd)
	rootCmd.AddCommand(ExportCmd)
	rootCmd.AddCommand(PreviewCmd)
	rootCmd.AddCommand(VersionCmd)

	err := rootCmd.Execute()
	utils.CheckFatalError(err, "Failed to execute command.")
}
