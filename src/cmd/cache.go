package cmd

import (
	"fmt"
	"os"

	h "github.com/FlawlessCasual17/rpm-get/helpers"
	"github.com/spf13/cobra"
)

// cacheCmd represents the cache command
var cacheCmd = &cobra.Command{
    Use:   "cache",
    Short: "A brief description of your command",
    Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("cache called")
    },
}

func init() {
    rootCmd.AddCommand(cacheCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // cacheCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // cacheCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// createCacheDir creates the cache directory.
func createCacheDir() {
    if err := os.MkdirAll(CACHE_DIR, 0755); err != nil {
        h.Printc("Unable to create cache dir!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    }
}

// createEtcDir creates the etc directory.
func createEtcDir() {
    if err := os.MkdirAll(ETC_DIR, 0755); err != nil {
        h.Printc("Unable to create etc dir!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    }
}
