package cmd

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"

    h "github.com/FlawlessCasual17/rpm-get/helpers"
    "github.com/spf13/cobra"
)

// BEGIN: Important variables/constants

// VERSION is the current version of _rpm-get_.
const VERSION string = "0.0.1"

var homeDir, _ = os.UserHomeDir()
// CacheDir is the directory where rpm-get will
// cache JSON files from **_GitHub_**/**_GitLab_**.
// As well as downloaded packages.
var CacheDir = filepath.Join(homeDir, ".cache/rpm-get")

const ETC_DIR string = "/etc/rpm-get"

// HOST_CPU is the host CPU architecture.
const HOST_CPU string = runtime.GOARCH

// UserAgent is the user agent string used for HTTP requests.
var UserAgent = fmt.Sprintf(
    "Mozilla/5.0 (X11; Linux %s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36", HOST_CPU)

// END: Important variables/constants

const SUCCESS_EXIT_CODE int = 0
const ERROR_EXIT_CODE int = 1
const USAGE_EXIT_CODE int = 2

var rootCmd = &cobra.Command {
    Use: "rpm-get",
    Short: "rpm-get is a CLI tool for downloading and managing RPM packages.",
    Long: `rpmget is a CLI tool for aquiring RPM packages that are not convieniently
    available in the default repositories. These can be either 3rd party repositories
    or direct download packages from the internet.`,
    Run: func(cmd *cobra.Command, _ []string) {
        Usage(); os.Exit(h.USAGE_EXIT_CODE)
    },
}

func Execute() {
    if error := rootCmd.Execute(); error != nil {
        h.Printc(error, h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    }
}

// spellcheck: ignore

// TODO: Add the following commands:
//   - install
//   - cache
//   - reinstall
//   - remove
//   - update
//   - upgrade
//   - info
//   - list
//   - search
//   - clean

func init() {
    for _, arg := range os.Args[1:] {
        switch arg {
        case "-version":
            h.Printc("`-version` is not a valid flag. Use `-v` or `--version` instead", h.WARNING, false)
            os.Exit(h.USAGE_EXIT_CODE)
        case "-help":
            h.Printc("`-help` is not a valid flag. Use `-h` or `--help` instead", h.WARNING, false)
            os.Exit(h.USAGE_EXIT_CODE)
        case "version":
            getVersion(); os.Exit(h.SUCCESS_EXIT_CODE)
        }
    }

    // Set custom usage function before defining flags
    flag.Usage = Usage

    hFlag := flag.Bool("h", false, "Display help information")
    helpFlag := flag.Bool("help", false, "Display help information")
    questionflag := flag.Bool("?", false, "Display help information")

    vFlag := flag.Bool("v", false, "Display version information")
    versionFlag := flag.Bool("version", false, "Display version information")

    // Parse flags
    flag.Parse()

    // Check for `--help` in remaining args
    for _, arg := range flag.Args() {
        switch arg {
        case "--help":
            Usage(); os.Exit(h.SUCCESS_EXIT_CODE)
        }
    }

    // Check if -h was used
    if *hFlag || *helpFlag || *questionflag {
        Usage(); os.Exit(h.SUCCESS_EXIT_CODE)
    } else if *vFlag || *versionFlag {
        getVersion(); os.Exit(h.SUCCESS_EXIT_CODE)
    }
}

// spellcheck: ignore

// getVersion prints the current version of rpm-get.
func getVersion() { fmt.Printf("rpm-get version: %s\n", VERSION) }

// isAdmin ensures that the user running rpm-get is using sudo or is a root.
func isAdmin() bool {
    if os.Geteuid() != 0 || os.Getenv("SUDO_USER") != "" {
        h.Printc("rpm-get must be run as root.", h.WARNING, false)

        return false
    }

    h.Printc("rpm-get is running as root.", h.INFO, true)

    return true
}

// spellcheck: ignore

// which looks for the given command in the
// PATH and prints an error if it's not found.
func which(cmd string)  {
    _, err := exec.LookPath(cmd)
    msg := fmt.Sprintf("Command `%s` not found in PATH. Exiting...", cmd)

    if err != nil { h.Printc(msg, h.FATAL, false) }
}

// createCacheDir creates the cache directory.
func createCacheDir() {
    err := os.MkdirAll(CacheDir, 0755)
    if err != nil { h.Printc("Unable to create cache dir!", h.FATAL, false) }
}

// createEtcDir creates the etc directory.
func createEtcDir() {
    err := os.MkdirAll(ETC_DIR, 0755)
    if err != nil { h.Printc("Unable to create etc dir!", h.FATAL, false) }
}

// getReleases retrieves the list of releases from the GitHub/GitLab API.
// func getReleases() {
// }
