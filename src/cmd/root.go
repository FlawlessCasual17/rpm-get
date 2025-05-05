package cmd

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "runtime"
    "strings"

    // third-party imports
    h "github.com/FlawlessCasual17/rpm-get/helpers"
    "github.com/PuerkitoBio/goquery"
    "github.com/samber/lo"
    "github.com/schollz/progressbar/v3"
    "github.com/spf13/cobra"
)

// VERSION is the current version of rpm-get.
const VERSION string = "0.0.1"

var (
    Project = ""
    RelType = ""
    Creator = ""
    ProjectId = ""
    ConfigDir = filepath.Join(os.Getenv("HOME"), ".config/rpm-get")
    ConfigFile = filepath.Join(ConfigDir, "config.json")
    // UserAgent is the user agent string used for HTTP requests.
    UserAgent = fmt.Sprintf(
        "Mozilla/5.0 (X11; Linux %s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
        HOST_CPU)
    GhHeaderAuth = fmt.Sprintf("Bearer %s", getEnv("GITHUB_TOKEN"))
    GlHeaderAuth = getEnv("GITLAB_TOKEN")
)

// ETC_DIR is the directory where rpm-get will store repositories.
const ETC_DIR string = "/etc/rpm-get"

// TODO: Add support for Zypper repos

// YUM_REPOS_DIR is the directory where rpm-get will store repositories.
const YUM_REPOS_DIR string = "/etc/yum.repos.d"

// CACHE_DIR is the directory where rpm-get will
// cache JSON files from GitHub/GitLab.
// As well as downloaded packages.
const CACHE_DIR string = "/var/cache/rpm-get"

// HOST_CPU is the host CPU architecture.
const HOST_CPU string = runtime.GOARCH

// MAIN_REPO is the main repository for rpm-get.
const MAIN_REPO string = "https://github.com/FlawlessCasual17/rpm-get"

var rootCmd = &cobra.Command {
    Use: "rpm-get",
    Short: "rpm-get is a CLI tool for downloading and managing RPM packages.",
    Long: `rpm-get is a CLI tool for aquiring RPM packages that are not convieniently
    available in the default repositories. These can be either 3rd party repositories
    or direct download packages from the internet.`,
    Run: func(_ *cobra.Command, _ []string) {
        Usage(); os.Exit(h.USAGE_EXIT_CODE)
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        h.Printc(err.Error(), h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    }
}

func init() {
    for _, arg := range os.Args[1:] {
        switch arg {
        case "-version":
            h.Printc("`-version` is not a valid flag. Use `-v` or `--version` instead", h.WARNING, false)
            os.Exit(h.USAGE_EXIT_CODE)
        case "-help":
            h.Printc("`-help` is not a valid flag. Use `-h` or `--help` instead", h.WARNING, false)
            os.Exit(h.USAGE_EXIT_CODE)
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

// createCacheDir creates the cache directory.
func createCacheDir() {
    if err := os.MkdirAll(CACHE_DIR, 0755); err != nil {
        h.Printc("Unable to create cache dir!", h.FATAL, false)
    }
}

// createEtcDir creates the etc directory.
func createEtcDir() {
    if err := os.MkdirAll(ETC_DIR, 0755); err != nil {
        h.Printc("Unable to create etc dir!", h.FATAL, false)
    }
}

// getReleases retrieves the list of releases from the GitHub/GitLab API.
func getReleases() {
    filePath := fmt.Sprintf("%s_cache.json", Project)
    cacheFilePath := filepath.Join(CACHE_DIR, filePath)
    url := ""
    feedbackMsg := ""
    request, _ := http.NewRequest("", url, nil)
    request.Header.Set("User-Agent", UserAgent)

    switch RelType {
    case "github":
        baseUrl := "https://api.github.com/repos"
        url = fmt.Sprintf("%s/%s/%s/releases/latest", baseUrl, Creator, Project)
        request.Header.Set("Authorization", GhHeaderAuth)
        feedbackMsg = parseJson(cacheFilePath, rateLimitedMsg {})
    case "gitlab":
        baseUrl := "https://gitlab.com/api/v4/projects"
        url = fmt.Sprintf("%s/%s/releases/permalink/latest", baseUrl, ProjectId)
        request.Header.Set("PRIVATE-TOKEN", GlHeaderAuth)
        v, _ := os.ReadFile(cacheFilePath);
        feedbackMsg = string(v)
    }

    if _, err := os.Stat(CACHE_DIR); err != nil && os.IsNotExist(err) {
        createCacheDir()
    }

    // NOTE: `//nolint:all` is used to suppress annoying linter warnings/errors.

    lo.TryCatch(func() error { // try
        resp, _ := http.DefaultClient.Do(request)
        //nolint:all
        defer resp.Body.Close()

        file, _ := os.OpenFile(cacheFilePath, os.O_CREATE|os.O_WRONLY, 0644)
        //nolint:all
        defer file.Close()

        // NOTE: Might switch to "github.com/cheggaaa/pb" if this doesn't meet my needs.
        bar := progressbar.DefaultBytes(resp.ContentLength, "Downloading...")
        //nolint:all
        io.Copy(io.MultiWriter(file, bar), resp.Body)

        return nil
    }, func() { // catch
        h.Printc("Unable to create cache file!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    })

    if rateLimited(feedbackMsg) {
        h.Printc("API rate limit exceeded!", h.WARNING, true)
        h.Printc("Deleting cache file...", h.INFO, true)
        if err := os.Remove(cacheFilePath); err != nil {
            h.Printc(err.Error(), h.ERROR, false)
        }
    }
}

//nolint:unused
type rateLimitedMsg struct {
    message string
}

// TODO: Add handling for more complex JSON queries

// parseJson parses the JSON file at the given path.
func parseJson(filePath string, jsonPath any) string {
    switch obj := jsonPath.(type) {
    case rateLimitedMsg:
        data, _ := os.ReadFile(filePath)

        err := json.Unmarshal(data, &obj)
        if err != nil { h.Printc(err.Error(), h.FATAL, false) }

        return obj.message
    default:
        return ""
    }
}

// rateLimited checks if the given feedback message contains a rate limit error.
func rateLimited(feedbackMsg string) bool {
    targets := []string { "API rate limit exceeded", "API rate limit exceeded for" }
    return strings.Contains(feedbackMsg, targets[0]) || strings.Contains(feedbackMsg, targets[1])
}

// scrapeWebsite parses a website and returns the matches of a given regex.
func scrapeWebsite(url string, regex string, elementRefs []string) string {
    newRegex, _ := regexp.Compile(regex)
    request, _ := http.NewRequest("", url, nil)
    request.Header.Set("User-Agent", UserAgent)
    result := ""

    lo.TryCatch(func() error { // try
        resp, _ := http.DefaultClient.Do(request)
        //nolint:all
        defer resp.Body.Close()

        doc, err := goquery.NewDocumentFromReader(resp.Body)
        if err != nil { h.Printc(err.Error(), h.WARNING, true) }

        // Parse HTML
        selection := doc.Find(elementRefs[0])
        selection.Each(func(i int, s *goquery.Selection) {
            element := s.Find(elementRefs[1]).Text()
            match := newRegex.FindString(element)
            if newRegex.MatchString(element) { result = match }
        })

        return nil
    }, func() { // catch
        h.Printc("Failed to scrape the requested website!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    })

    return result
}

// getSha256Hash returns the SHA256 hash of the given file.
func getSha256Hash(filePath string) string {
    file, _ := os.Open(filePath)
    //nolint:all
    defer file.Close()

    hash := sha256.New()
    if _, err := io.Copy(hash, file); err != nil {
        h.Printc(err.Error(), h.ERROR, false)
        return ""
    }

    return hex.EncodeToString(hash.Sum(nil))
}

// downloadPkg downloads the requested RPM package.
func downloadPkg(url string, filePath string) {
    cacheFilePath := filepath.Join(CACHE_DIR, filePath)
    request, _ := http.NewRequest("", url, nil)
    request.Header.Set("User-Agent", UserAgent)

    lo.TryCatch(func() error { // try
        resp, _ := http.DefaultClient.Do(request)
        //nolint:all
        defer resp.Body.Close()

        file, _ := os.OpenFile(cacheFilePath, os.O_CREATE|os.O_WRONLY, 0644)
        //nolint:all
        defer file.Close()

        bar := progressbar.DefaultBytes(resp.ContentLength, "Downloading...")
        //nolint:all
        io.Copy(io.MultiWriter(file, bar), resp.Body)

        return nil
    }, func() { // catch
        h.Printc("Failed to download the requested RPM!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    })
}

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
    } else { println(out) }
}

// upgradePkg upgrades the given RPM packages.
func upgradePkg(pkgs []string) {
    if !isAdmin() {
        h.Printc("rpm-get must be run as root!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    }

    cmd := which("sudo") + " " + which("dnf")
    args := []string { "install", "-y", strings.Join(pkgs, " ") }
    command := exec.Command(cmd, args...)
    out, err := command.Output()

    if err != nil {
        h.Printc(err.Error(), h.ERROR, false)
    } else { println(out) }
}

// checkUpdates checks for updates to the packages list.
func checkUpdates() {
    url := MAIN_REPO + "/raw/refs/heads/master/packages/packages-list.json"
    tmpFilePath := filepath.Join(ConfigDir, "packages-list.json.tmp")
    filePath := filepath.Join(ConfigDir, "packages-list.json")
    request, _ := http.NewRequest("", url, nil)
    request.Header.Set("User-Agent", UserAgent)

    // Download packages-list.json
    lo.TryCatch(func() error { // try
        resp, _ := http.DefaultClient.Do(request)
        //nolint:all
        defer resp.Body.Close()

        file, _ := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY, 0644)
        //nolint:all
        defer file.Close()

        bar := progressbar.DefaultBytes(resp.ContentLength, "Updating packages list...")
        //nolint:all
        io.Copy(io.MultiWriter(file, bar), resp.Body)

        return nil
    }, func() { // catch
        h.Printc("Unable to update packages list!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    })

    // Compare the hashes of the downloaded file and the existing file
    tmpListHash := getSha256Hash(tmpFilePath)
    listHash := getSha256Hash(filePath)

    if tmpListHash != listHash {
        // Attempt to move the downloaded file to the existing file
        if err := os.Rename(tmpFilePath, filePath); err != nil {
            h.Printc(err.Error(), h.ERROR, true)
        } else {
            h.Printc("Packages list was sucessfully updated!", h.INFO, true)
        }
    } else {
        h.Printc("Packages list is already up to date!", h.INFO, true)
    }
}

// removePkg removes the requested RPM package.
func removePkg(pkg string) {
    if !isAdmin() {
        h.Printc("rpm-get must be run as root!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    }

    cmd := which("sudo") + " " + which("dnf")
    args := []string { "remove", "-y", pkg }
    command := exec.Command(cmd, args...)
    out, err := command.Output()

    if err != nil {
        h.Printc(err.Error(), h.ERROR, false)
    } else { println(out) }
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
    } else { println(out) }
}

// addRepo adds the given RPM repo to the YUM repos directory.
func addRepo(repoUrl string) {
    if !isAdmin() {
        h.Printc("rpm-get must be run as root!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    }

    repoName := strings.Split(repoUrl, "/")[-0]
    tmpFilePath := filepath.Join(YUM_REPOS_DIR, repoName + ".tmp")
    filePath := filepath.Join(YUM_REPOS_DIR, repoName)
    request, _ := http.NewRequest("", repoUrl, nil)
    request.Header.Set("User-Agent", UserAgent)

    // Download packages-list.json
    lo.TryCatch(func() error { // try
        resp, _ := http.DefaultClient.Do(request)
        //nolint:all
        defer resp.Body.Close()

        file, _ := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY, 0644)
        //nolint:all
        defer file.Close()

        bar := progressbar.DefaultBytes(resp.ContentLength, "Downloading RPM repo...")
        //nolint:all
        io.Copy(io.MultiWriter(file, bar), resp.Body)

        return nil
    }, func() { // catch
        h.Printc("Unable to update packages list!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    })

    if err := os.Rename(tmpFilePath, filePath); err != nil {
        h.Printc(err.Error(), h.ERROR, false)
    } else {
        msg := fmt.Sprintf("Successfully added the repo for %s", Project)
        h.Printc(msg, h.INFO, true)
    }
}

func addCoprRepo(repoName string) {
    if !isAdmin() {
        h.Printc("rpm-get must be run as root!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    }

    cmd := which("sudo") + " " + which("dnf")
    args := []string { "copr", "enable", "-y", repoName }
    command := exec.Command(cmd, args...)
    out, err := command.Output()

    if err != nil {
        h.Printc(err.Error(), h.ERROR, false)
    } else {
        println(out)
        msg := fmt.Sprintf("Successfully added the repo for %s\n", Project)
        h.Printc(msg, h.INFO, true)
    }
}
