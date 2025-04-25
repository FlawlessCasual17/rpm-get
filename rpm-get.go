package main

import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/fatih/color"
)

// Important variables/constants

// `VERSION` is the current version of _rpm-get_.
const VERSION string = "0.0.1"

var homeDir, _ = os.UserHomeDir()
// `CACHE_DIR` is the directory where rpm-get will
// cache JSON files from **_GitHub_**/**_GitLab_**.
var CACHE_DIR string = filepath.Join(homeDir, ".cache/rpm-get")

// Important variables/constants

// Message types

const INFO string = "INFO"
const PROGRESS string = "PROGRESS"
const WARNING string = "WARNING"
const ERROR string = "ERROR"
const FATAL string = "FATAL"

// Message types

func main() {
    // homeDir, _ := os.UserHomeDir()

}

func getVersion() string { return VERSION }

// `printc` prints messages with colored text to the console.
func printc(msg string, msgType string, newLine bool) {
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
        default:
            fmt.Printf("%s  [%s]: %s\n", cr, GRAY("UNKNOWN"), msg)
    }
}
