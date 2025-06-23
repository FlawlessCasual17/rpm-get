package cmd

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"

    // third-party imports
    h "github.com/FlawlessCasual17/rpm-get/helpers"
    "github.com/samber/lo"
    "github.com/schollz/progressbar/v3"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command {
    Use: "rpm-get",
    Short: "rpm-get is a CLI tool for downloading and managing RPM packages.",
    Long: `rpm-get is a CLI tool for aquiring RPM packages that are not convieniently
available in the default repositories.
These can be either 3rd party repositories or direct download packages from the internet.`,
    Run: func(cmd *cobra.Command, _ []string) {
        if wantsVersion {
            getVersion(); os.Exit(h.SUCCESS_EXIT_CODE)
        }

        _ = cmd.Help(); os.Exit(h.USAGE_EXIT_CODE)
    },
}

var (
    wantsVersion bool
    isHTML = false
    App = ""
    Project = ""
    RelType = ""
    Creator = ""
    ProjectID = ""
    RepoName = ""
    ConfigDir = filepath.Join(os.Getenv("HOME"), ".config/rpm-get")
    ConfigFile = filepath.Join(ConfigDir, "config.json")
    DataDir = filepath.Join(os.Getenv("HOME"), ".local/share/rpm-get")
    // UserAgent is the user agent string used for HTTP requests.
    UserAgent = fmt.Sprintf(
        "Mozilla/5.0 (X11; Linux %s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
        HOST_CPU)
    GhHeaderAuth = fmt.Sprint("Bearer " + getEnv("GITHUB_TOKEN"))
    GlHeaderAuth = getEnv("GITLAB_TOKEN")
)

const (
    // VERSION is the current version of rpm-get.
    VERSION string = "0.0.1"

    // ETC_DIR is the directory where rpm-get will store repositories.
    ETC_DIR string = "/etc/rpm-get"

    // TODO: Add support for Zypper repos

    // YUM_REPOS_DIR is the directory where rpm-get will store RPM repositories.
    YUM_REPOS_DIR string = "/etc/yum.repos.d"

    // CACHE_DIR is the directory where rpm-get will
    // cache JSON files from GitHub/GitLab.
    // As well as downloaded packages.
    CACHE_DIR string = "/var/cache/rpm-get"

    // HOST_CPU is the host CPU architecture.
    HOST_CPU string = runtime.GOARCH

    // MAIN_REPO is the main repository for rpm-get.
    MAIN_REPO string = "https://github.com/FlawlessCasual17/rpm-get"

    // PKGS_REPO is the repository for package manifests.
    PKGS_REPO string = "https://github.com/FlawlessCasual17/rpm-get.Packages"
)

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        h.Printc(err.Error(), h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    }
}

func init() {
    // // Set custom usage function before defining flags
    // rootCmd.Flags().Usage = func () { _ = rootCmd.Usage() }
    // rootCmd.Flags().BoolVar(&help, "help", "h", "", "Display help information")
    // rootCmd.Flags().Bool("?", false, "Display help information")

    // vFlag := rootCmd.Flags().Bool("v", false, "Display verbose information")
    rootCmd.Flags().BoolVar(&wantsVersion, "version", false, "Display version information")

    // // Parse flags
    // rootCmd.Flags().Parse()

    // // Check for `--help` in remaining args
    // for _, arg := range rootCmd.Flags().Args() {
    //     switch arg {
    //     case "--help":
    //         _ = rootCmd.Usage(); os.Exit(h.SUCCESS_EXIT_CODE)
    //     }
    // }

    // // Check if -h was used
    // if *hFlag || *helpFlag || *questionflag {
    //     _ = rootCmd.Usage(); os.Exit(h.SUCCESS_EXIT_CODE)
    // } else if *versionFlag {
    //     getVersion(); os.Exit(h.SUCCESS_EXIT_CODE)
    // }
}

// spellcheck: ignore

// getEnv returns the value of the environment variable. Empty string if not found.
func getEnv(key string) string {
    result, ok := os.LookupEnv(key)
    return lo.Ternary(ok, result, "")
}

// spellcheck: ignore

// isAdmin ensures that the user running rpm-get is using sudo or is a root.
func isAdmin() bool {
    return os.Geteuid() == 0 || os.Getenv("SUDO_USER") == ""
}

// spellcheck: ignore

// which looks for the given command in the
// PATH and prints an error if it's not found.
func which(cmd string) string {
    result, err := exec.LookPath(cmd)
    if err != nil { return "" }
    return result
}

// getSha256Hash returns the SHA256 hash of the given file.
func getSha256Hash(filePath string) string {
    file, fileErr := os.Open(filePath)
    if fileErr != nil {
        h.Printc("Failed to open file!", h.ERROR, false)
        return ""
    }
    //nolint:errcheck
    defer file.Close()

    hash := sha256.New()
    if _, err := io.Copy(hash, file); err != nil {
        h.Printc(err.Error(), h.ERROR, false)
        return ""
    }

    return hex.EncodeToString(hash.Sum(nil))
}

// downloadPkg downloads the requested RPM package.
func downloadPkg(url string, filePath string) error {
    downloadError := error (nil)
    cacheFilePath := filepath.Join(CACHE_DIR, filePath)

    lo.TryCatch(func() error { // try
        request, _ := http.NewRequest("GET", url, nil)
        request.Header.Set("User-Agent", UserAgent)
        resp, err := http.DefaultClient.Do(request)
        if err != nil {
            h.Printc("Request failed!", h.ERROR, false)
            downloadError = fmt.Errorf("Request failed: %w", err)
            return downloadError
        }
        //nolint:errcheck
        defer resp.Body.Close()

        file, _ := os.OpenFile(cacheFilePath, os.O_CREATE|os.O_WRONLY, 0644)
        //nolint:errcheck
        defer file.Close()

        bar := progressbar.DefaultBytes(resp.ContentLength, "Downloading...")

        if _, err := io.Copy(io.MultiWriter(file, bar), resp.Body); err != nil {
            h.Printc("Failed to download the requested RPM package!", h.ERROR, false)
            downloadError = fmt.Errorf("Failed to download the requested RPM package: %w", err)
            return downloadError
        }

        return nil
    }, func() { // catch
        h.Printc("Failed to download the requested RPM package!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    })

    return downloadError
}

// addRepo adds the given RPM repo to the YUM repos directory.
func addRepo(repoUrl string) {
    if !isAdmin() {
        h.Printc("rpm-get must be run as root!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    }

    baseName := strings.Split(repoUrl, "/")[-0]
    tmpFilePath := filepath.Join(YUM_REPOS_DIR, baseName + ".tmp")
    filePath := strings.ReplaceAll(tmpFilePath, ".tmp", "")

    // Download packages-list.json
    lo.TryCatch(func() error { // try
        request, _ := http.NewRequest("GET", repoUrl, nil)
        request.Header.Set("User-Agent", UserAgent)
        resp, respErr := http.DefaultClient.Do(request)
        if respErr != nil {
            h.Printc("Request failed!", h.ERROR, false)
            return respErr
        }
        //nolint:errcheck
        defer resp.Body.Close()

        file, _ := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY, 0644)
        //nolint:errcheck
        defer file.Close()

        bar := progressbar.DefaultBytes(resp.ContentLength, "Downloading RPM repo...")
        //nolint:errcheck
        io.Copy(io.MultiWriter(file, bar), resp.Body)

        return nil
    }, func() { // catch
        h.Printc("Unable to update packages list!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    })

    if err := os.Rename(tmpFilePath, filePath); err != nil {
        h.Printc(err.Error(), h.ERROR, false)
    } else {
        RepoName = baseName
        msg := fmt.Sprint("Successfully added the repo for " + App)
        h.Printc(msg, h.INFO, true)
    }
}

// addCoprRepo adds the given Fedora COPR repo to the YUM repos directory.
func addCoprRepo(username string, project string) error {
    if !isAdmin() {
        h.Printc("rpm-get must be run as root!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    }

    cmd := which("sudo") + " " + which("dnf")
    coprRepo := username + "/" + project
    args := []string { "copr", "enable", "-y", coprRepo }
    command := exec.Command(cmd, args...)
    out, err := command.Output()

    if err != nil {
        h.Printc(err.Error(), h.ERROR, false)
        return fmt.Errorf("Command failed: %w", err)
    }

    RepoName = fmt.Sprintf("_copr:copr.fedorainfracloud.org:%s:%s", username, project)
    println(out)
    msg := fmt.Sprint("Successfully added the repo for " + App)
    h.Printc(msg, h.INFO, true)

    return nil
}

// removeRepo removes the repo of an application from the YUM repos directory.
func removeRepo() (bool, error) {
    if !isAdmin() {
        h.Printc("rpm-get must be run as root!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    }

    filePath := filepath.Join(YUM_REPOS_DIR, RepoName)

    if err := os.Remove(filePath); err != nil {
        msg := fmt.Sprint("Failed to remove the repo for " + App)
        h.Printc(msg, h.ERROR, false)
        return false, fmt.Errorf("%s: %w", msg, err)
    } else {
        msg := fmt.Sprint("Successfully removed the repo for " + App)
        h.Printc(msg, h.INFO, true)
        return true, nil
    }
}
