package cmd

import (
	"os"
	"os/exec"

	h "github.com/FlawlessCasual17/rpm-get/helpers"
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

// installPkg installs the requested RPM package.
func installPkg(pkg string) {
    if !isAdmin() {
        h.Printc("rpm-get must be run as root!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    }

    cmd := which("sudo") + " " + which("dnf")
    args := []string { "install", "-y", pkg }
    command := exec.Command(cmd, args...)
    out, err := command.Output()

    if err != nil {
        h.Printc(err.Error(), h.ERROR, false)
    } else {
        println(out)
    }
}
