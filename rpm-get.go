// rpm-get is a CLI tool to download and manage rpm packages
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
)

// Important variables/constants

// VERSION is the current version of _rpm-get_.
const VERSION string = "0.0.1"

var homeDir, _ = os.UserHomeDir()
// CACHE_DIR is the directory where rpm-get will
// cache JSON files from **_GitHub_**/**_GitLab_**.
// As well as downloaded packages.
var CACHE_DIR = filepath.Join(homeDir, ".cache/rpm-get")

const ETC_DIR string = "/etc/rpm-get"

// HOST_CPU is the host CPU architecture.
const HOST_CPU string = runtime.GOARCH

// USER_AGENT is the user agent string used for HTTP requests.
var USER_AGENT = fmt.Sprintf(
    "Mozilla/5.0 (X11; Linux %s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36", HOST_CPU)

// Important variables/constants

// Message types

// INFO is the message type for informational messages.
const INFO string = "INFO"
// PROGRESS is the message type for progress messages.
const PROGRESS string = "PROGRESS"
// WARNING is the message type for warning messages.
const WARNING string = "WARNING"
// ERROR is the message type for error messages.
const ERROR string = "ERROR"
// FATAL is the message type for fatal messages.
// This is followed by `os.Exit(1)`.
const FATAL string = "FATAL"

// Message types

func main() {
    // homeDir, _ := os.UserHomeDir()

}

// spellcheck: ignore

// usage prints the help text when rpm-get is called without arguments,
// Or called with the `help` argument.
func usage() {
    fmt.Printf(`
rpm-get version %s


Usage

rpm-get {update [--repos-only] [--quiet] | upgrade [--dg-only] | show <pkg list> | install <pkg list>
        | reinstall <pkg list> | remove [--remove-repo] <pkg list>
        | search [--include-unsupported] <regex> | cache | clean
        | list [--include-unsupported] [--raw|--installed|--not-installed]
        | prettylist [<repo>] | csvlist [<repo>] | fix-installed [--old-apps]
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

show
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

prettylist
    markdown formatted list the packages available in repo. repo defaults to
    01-main. If repo is 00-builtin or 01-main the packages from 00-builtin are
    included. Use this to update README.md.

csvlist
    csv formatted list the packages available in repo. repo defaults to
    01-main. If repo is 00-builtin or 01-main the packages from 00-builtin are
    included. Use this with 3rd party wrappers.

cache
    list the contents of the rpm-get cache (/var/cache/rpm-get).

fix-installed
    fix installed packages whose definitions were changed. When --old-apps is
provided, transition packages to new format. This command is only intended
    for internal use.

help
    show this help.

version
    show rpm-get version.
`, VERSION)
}

func getVersion() string {
    return fmt.Sprintf("rpm-get version %s", VERSION)
}

// checkPrivileges ensures that the user running rpm-get is using sudo or is a root.
func checkPrivileges() {
    if os.Geteuid() != 0 || os.Getenv("SUDO_USER") != "" {
        printc("rpm-get must be run as root.", FATAL, false)
    }

    printc("rpm-get is running as root.", INFO, false)
}

// checkPath looks for the given command in the
// PATH and prints an error if it's not found.
func checkPath(cmd string)  {
    _, err := exec.LookPath(cmd)
    msg := fmt.Sprintf("Command `%s` not found in PATH. Exiting...", cmd)

    if err != nil { printc(msg, FATAL, false) }
}

// createCacheDir creates the cache directory.
func createCacheDir() {
    err := os.Mkdir(CACHE_DIR, 0755)
    if err != nil { printc("Unable to create cache dir!", ERROR, false) }
}

// createEtcDir creates the etc directory.
func createEtcDir() {
    err := os.Mkdir(ETC_DIR, 0755)
    if err != nil { printc("Unable to create etc dir!", ERROR, false) }
}

// `printc` prints messages with colored text to the console.
func printc(msg string, msgType any, newLine bool) {
    RED := color.New(color.FgRed).SprintFunc()
    GREEN := color.New(color.FgGreen).SprintFunc()
    YELLOW := color.New(color.FgYellow).SprintFunc()
    BLUE := color.New(color.FgBlue).SprintFunc()
    ORANGE := color.New(color.FgHiYellow).SprintFunc()
    GRAY := color.New(color.FgHiBlack).SprintFunc()
    // RESET := color.New(color.Reset).SprintFunc()

    cr := "\n"
    if !newLine { cr = "" }

    switch msgType {
    case INFO:
        fmt.Printf("%s  [%s]: %s\n", cr, GREEN(INFO), msg)
    case PROGRESS:
        fmt.Printf("%s  [%s]: %s\n", cr, BLUE(PROGRESS), msg)
    case WARNING:
        fmt.Printf("%s  [%s]: %s\n", cr, YELLOW(WARNING), msg)
    case ERROR:
        fmt.Printf("%s  [%s]: %s\n", cr, RED(ERROR), msg)
    case FATAL:
        fmt.Printf("%s  [%s]: %s\n", cr, ORANGE(FATAL), msg)
        os.Exit(1)
    default:
        fmt.Printf("%s  [%s]: %s\n", cr, GRAY("UNKNOWN"), msg)
    }
}
