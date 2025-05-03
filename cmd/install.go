package cmd

import (
    "github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
    Use:   "install",
    Short: "Install packages",
    Long: "Install packages",
    Run: func(_ *cobra.Command, _ []string) {
        getReleases()
    },
}

func init() { rootCmd.AddCommand(installCmd) }
