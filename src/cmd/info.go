package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	h "github.com/FlawlessCasual17/rpm-get/helpers"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
    Use:   "info",
    Short: "A brief description of your command",
    Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("info called")
    },
}

func init() {
    rootCmd.AddCommand(infoCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // infoCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// PkgInfo returns the package information for the given package.
func pkgInfo(pkg string) (string, error) {
    result := ""
    data := Pkg {}
    filePath := filepath.Join(DataDir, pkg + ".yaml")

    content, readErr := os.ReadFile(filePath)
    if readErr != nil {
        h.Printc("Failed to read file!", h.ERROR, false)
        return result, fmt.Errorf("Failed to read file: %w", readErr)
    }

    if err := yaml.Unmarshal(content, &data); err != nil {
        h.Printc("Failed to unmarshal file!", h.ERROR, false)
        return result, fmt.Errorf("Failed to unmarshal file: %w", err)
    }

    result = fmt.Sprintf(`
        Supported OS: %s
        Version: %s
        Name: %s
        License: %s
        Homepage: %s
        Description: %s
        Notes: %s
        Pkg Arches: %s
        `,
        strings.Join(data.supported_os, ", "),
        data.version,
        data.name,
        data.license.licenseString,
        data.homepage,
        data.description,
        data.notes,
        strings.Join(data.pkg_arches, ", "))

    return result, nil
}
