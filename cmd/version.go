/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	h "github.com/FlawlessCasual17/rpm-get/helpers"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Show version",
    Long: "Show version",
    Run: func(_ *cobra.Command, _ []string) { getVersion(); os.Exit(h.SUCCESS_EXIT_CODE) },
}

func init() { rootCmd.AddCommand(versionCmd) }

// getVersion prints the current version of rpm-get.
func getVersion() { fmt.Printf("rpm-get version: %s\n", VERSION) }
