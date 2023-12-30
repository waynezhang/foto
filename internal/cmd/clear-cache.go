package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynezhang/foto/internal/cache"
)

var ClearCacheCmd = &cobra.Command {
  Use: "clear-cache",
  Short: "Clear local cache",
  Run: func(cmd *cobra.Command, args []string) {
    cache.Shared().Clear()
  },
}
