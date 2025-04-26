// rpm-get is a CLI tool to download and manage rpm packages
package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"

    // "github.com/spf13/cobra"
    "github.com/FlawlessCasual17/rpm-get/cmds"
    "github.com/fatih/color"
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

// BEGIN: Message types

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

// END: Message types

const SUCCESS_EXIT_CODE int = 0
const ERROR_EXIT_CODE int = 1

func main() {
    flag.Usage = cmds.Usage

    helpFlag := flag.Bool("help", false, "Print help message.")

    flag.Parse()

    if *helpFlag {
        cmds.Usage()
        os.Exit(SUCCESS_EXIT_CODE)
    }

    args := flag.Args()

    if len(args) > 0 {
        switch args[0] {
        case "help":
            cmds.Usage()
            os.Exit(SUCCESS_EXIT_CODE)
        }
    }
}

// spellcheck: ignore

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

// spellcheck: ignore

// which looks for the given command in the
// PATH and prints an error if it's not found.
func which(cmd string)  {
    _, err := exec.LookPath(cmd)
    msg := fmt.Sprintf("Command `%s` not found in PATH. Exiting...", cmd)

    if err != nil { printc(msg, FATAL, false) }
}

// createCacheDir creates the cache directory.
func createCacheDir() {
    err := os.MkdirAll(CacheDir, 0755)
    if err != nil { printc("Unable to create cache dir!", ERROR, false) }
}

// createEtcDir creates the etc directory.
func createEtcDir() {
    err := os.MkdirAll(ETC_DIR, 0755)
    if err != nil { printc("Unable to create etc dir!", ERROR, false) }
}

// getReleases retrieves the list of releases from the GitHub/GitLab API.
func getReleases() {

}

// `printc` prints messages with colored text to the console.
func printc(msg any, msgType any, newLine bool) {
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
        os.Exit(ERROR_EXIT_CODE)
    default:
        fmt.Printf("%s  [%s]: %s\n", cr, GRAY("UNKNOWN"), msg)
    }
}
