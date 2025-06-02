/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	h "github.com/FlawlessCasual17/rpm-get/helpers"
	"github.com/spf13/cobra"
)

// reinstallCmd represents the reinstall command
var reinstallCmd = &cobra.Command{
    Use:   "reinstall",
    Short: "A brief description of your command",
    Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("reinstall called")
    },
}

func init() {
    rootCmd.AddCommand(reinstallCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // reinstallCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // reinstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// reinstallPkg reinstalls the requested RPM package that is already installed.
func reinstallPkg(pkg string) {
    if !isAdmin() {
        h.Printc("rpm-get must be run as root!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    }

    cmd := which("sudo") + " " + which("dnf")
    args := []string { "reinstall", "-y", pkg }
    command := exec.Command(cmd, args...)
    out, err := command.Output()

    if err != nil {
        h.Printc(err.Error(), h.ERROR, false)
    } else {
        println(out)
    }
}
