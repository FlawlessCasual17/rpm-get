package cmd

import (
    "fmt"
    "os"

    h "github.com/FlawlessCasual17/rpm-get/helpers"
    "github.com/spf13/cobra"
)

var helpCmd = &cobra.Command {
    Use: "help",
    Short: "Show usage",
    Long: "Show usage",
    Run: func(_ *cobra.Command, _ []string) {
        Usage(); os.Exit(h.SUCCESS_EXIT_CODE)
    },
}

func init() { rootCmd.AddCommand(helpCmd) }

// Usage prints the help text when rpm-get is called without arguments,
// Or called with the `help` argument.
func Usage() {
    fmt.Print(`


Usage

rpm-get {update [--repos-only] [--quiet] | upgrade [--dg-only] | info <pkg list> | install <pkg list>
        | reinstall <pkg list> | remove [--remove-repo] <pkg list>
        | search [--include-unsupported] <regex> | cache | clean
        | list [--include-unsupported] [--raw|--installed|--not-installed]
        | help | version}

rpm-get provides a high-level commandline interface for the package management
system to easily install and update packages published in 3rd party rpm
repositories or via direct download.

update
    update is used to resynchronize the package index files from their sources.
    When --repos-only is provided, only initialize and update rpm-get's
    external repositories, without updating rpm or looking for updates of
    installed packages.
    When --quiet is provided the fetching of rpm-get repository updates is done without progress feedback.

upgrade
    upgrade is used to install the newest versions of all packages currently
    installed on the system.
    When --dg-only is provided, only the packages which have been installed by rpm-get will be upgraded.

install
    install is followed by-repo is provided, also remove the rpm repository
    of rpm/ppa packages.

clean
    clean clears out the local repository (/var/cache/rpm-get) of retrieved
    package files.

search
    search for the given regex(7) term(s) from the list of available packages
    supported by rpm-get and display matches. When --include-unsupported is
    provided, include packages with unsupported architecture or upstream
    codename and include PPAs for Debian-derived distributions.

info
    show information about the given package (or a space-separated list of
    packages) including their install source and update mechanism.

list
    list the packages available via rpm-get. When no option is provided, list
    all supported packages and tell which ones are installed (slower). When
    --include-unsupported is provided, include packages with unsupported
    architecture or upstream codename and include PPAs for Debian-derived
    distributions (faster). When --raw is provided, list all packages and do
    not tell which ones are installed (faster). When --installed is provided,
    only list the packages installed (faster). When --not-installed is provided,
    only list the packages not installed (faster).

cache
    list the contents of the rpm-get cache (/var/cache/rpm-get).

help
    show this help.

version
    show rpm-get version.
`)
}
