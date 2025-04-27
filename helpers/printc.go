package helpers

import (
    "fmt"
    "os"
    // third-party packages
    "github.com/fatih/color"
)

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

const SUCCESS_EXIT_CODE int = 0
const ERROR_EXIT_CODE int = 1
const USAGE_EXIT_CODE int = 2

// Printc prints messages with colored text to the console.
func Printc(msg any, msgType any, newLine bool) {
    RED := color.New(color.FgRed).SprintFunc()
    GREEN := color.New(color.FgGreen).SprintFunc()
    YELLOW := color.New(color.FgYellow).SprintFunc()
    BLUE := color.New(color.FgBlue).SprintFunc()
    ORANGE := color.RGB(255, 128, 0).SprintFunc()
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
