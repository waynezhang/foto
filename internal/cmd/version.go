package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynezhang/foto/internal/log"
)

var Version = ""
var Revision = ""

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("foto v%s+%s", Version, Revision)
	},
}
