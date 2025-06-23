package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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
	"github.com/goccy/go-json"
	"github.com/goccy/go-yaml"
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

// VERSION is the current version of rpm-get.
const VERSION string = "0.0.1"

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

// ETC_DIR is the directory where rpm-get will store repositories.
const ETC_DIR string = "/etc/rpm-get"

// TODO: Add support for Zypper repos

// YUM_REPOS_DIR is the directory where rpm-get will store RPM repositories.
const YUM_REPOS_DIR string = "/etc/yum.repos.d"

// CACHE_DIR is the directory where rpm-get will
// cache JSON files from GitHub/GitLab.
// As well as downloaded packages.
const CACHE_DIR string = "/var/cache/rpm-get"

// HOST_CPU is the host CPU architecture.
const HOST_CPU string = runtime.GOARCH

// MAIN_REPO is the main repository for rpm-get.
const MAIN_REPO string = "https://github.com/FlawlessCasual17/rpm-get"

// PKGS_REPO is the repository for package manifests.
const PKGS_REPO string = "https://github.com/FlawlessCasual17/rpm-get.Packages"

type LicenseObject struct {
    // Package license based on SPDX license format, and license list: https://spdx.org/licenses/
    identifier string   `yaml:"identifier"`
    // License URL
    url string          `yaml:"url"`
}

// Package license based on SPDX license format, and license list: https://spdx.org/licenses/
type License struct {
    // Package license based on SPDX license format, and license list: https://spdx.org/licenses/
    licenseString string     `yaml:",omitempty"`
    // Package license based on SPDX license format, and license list: https://spdx.org/licenses/
    license *LicenseObject   `yaml:",omitempty"`
}

type PkgArch struct {
    // Download URL for the architecture
    url string    `yaml:"url"`
}

type UrlRepo struct {
    // Package repository URL
    url string         `yaml:"url"`
    // Repository GPG key URL
    gpgKeyUrl string   `yaml:"gpg_key_url"`
}

type CoprRepo struct {
    // Copr user name
    username string   `yaml:"username"`
    // Copr project name
    project string    `yaml:"project"`
}

// Information about an RPM/Copr repository
type Repo struct {
    urlRepo *UrlRepo     `yaml:",omitempty"`
    coprRepo *CoprRepo   `yaml:",omitempty"`
}

// Custom script to extract version information. Supports BASH, FISH, ZSH, PowerShell (pwsh), Nushell, and Python
type Script struct {
    scriptType string   `yaml:"script_type,omitempty"`
    run string          `yaml:"run,omitempty"`
}

// Auto-update architecture-specific configuration
type AutoUpdatePkgArch struct {
    // Auto-update URL for the architecture
    url string            `yaml:"url"`
}

// Auto-update architecture-specific configuration
type AutoUpdateArch struct {
    x86_64 *AutoUpdatePkgArch   `yaml:"x86_64,omitempty"`
    x86 *AutoUpdatePkgArch      `yaml:"x86,omitempty"`
    arm64 *AutoUpdatePkgArch    `yaml:"arm64,omitempty"`
}

// Information about a GitHub (Gitea, Gogs, Forgejo, and Codeberg) repository
type GithubObject struct {
    // Can either be a GitHub username or organization name
    owner string      `yaml:"owner"`
    // GitHub repository name
    repo string       `yaml:"repo"`
}

// Can also be used with Gitea, Gogs, Forgejo, and Codeberg.
type GitHub struct {
    // Information about a GitHub (Gitea, Gogs, Forgejo, and Codeberg) repository. This must be in the format 'username/repository-name'
    githubString string    `yaml:",omitempty"`
    // Information about a GitHub (Gitea, Gogs, Forgejo, and Codeberg) repository
    github *GithubObject   `yaml:",omitempty"`
}

type GitLabGroup struct {
    // GitLab group name
    group string      `yaml:"group"`
    // GitLab sub-group name
    subGroup string   `yaml:"sub_group"`
    // GitLab project name
    project string    `yaml:"project"`
}

type GitLabProfile struct {
    // GitLab profile name
    profile string   `yaml:"profile"`
    // GitLab project name
    project string   `yaml:"project"`
}

type GitLab struct {
    // Information about a GitLab repository. This must be in the format 'group/sub-group/project'
    groupString string       `yaml:",omitempty"`
    // Information about a GitLab repository. This must be in the format 'profile/project'
    profileString string     `yaml:",omitempty"`
    // Information about a GitLab repository
    group *GitLabGroup       `yaml:",omitempty"`
    // Information about a GitLab repository
    profile *GitLabProfile   `yaml:",omitempty"`
}

// Schema for package manifests
type Pkg struct {
    // List of operating systems (that use RPM) supported by this package
    supported_os []string              `yaml:"supported_os"`
    // Package version
    version string                     `yaml:"version"`
    // Package name
    name string                        `yaml:"name"`
    // Package license based on SPDX license format, and license list: https://spdx.org/licenses/
    license *License                    `yaml:"license"`
    // Package homepage
    homepage string                    `yaml:"homepage"`
    // Package description
    description string                 `yaml:"description"`
    // Additional notes about the package
    notes string                       `yaml:"notes,omitempty"`
    // Architecture-specific download information
    pkg_arches []string                `yaml:"pkg_arches"`
    arch struct {
        x86_64 *PkgArch                `yaml:"x86_64,omitempty"`
        x86 *PkgArch                   `yaml:"x86,omitempty"`
        arm64 *PkgArch                 `yaml:"arm64,omitempty"`
    }                                  `yaml:"arch"`
    // Information about an RPM/Copr repository
    repo *Repo                         `yaml:"repo,omitempty"`
    // List of package dependencies
    depends []string                   `yaml:"depends,omitempty"`
    // List of recommended packages
    recommends []string                `yaml:"recommends,omitempty"`
    // List of suggested packages
    suggests []string                  `yaml:"suggests,omitempty"`
    // List of conflicting packages
    conflicts []string                 `yaml:"conflicts,omitempty"`
    // List of packages that this package replaces
    replaces []string                  `yaml:"replaces,omitempty"`
    // Auto-update configuration
    auto_update struct {
        // Configuration for checking package version.
        check_version struct {
            url string                 `yaml:"url"`
            jsonpath string            `yaml:"jsonpath,omitempty"`
            xpath string               `yaml:"xpath,omitempty"`
            script Script              `yaml:"script,omitempty"`
            regex string               `yaml:"regex,omitempty"`
            regex_replace string       `yaml:"regex_replace,omitempty"`
            use_latest bool            `yaml:"use_latest,omitempty"`
            github *GitHub             `yaml:"github,omitempty"`
            gitlab *GitLab             `yaml:"gitlab,omitempty"`
        }                              `yaml:"check_version"`
        // Architecture-specific auto-update information
        arch *AutoUpdateArch           `yaml:"arch"`
    }                                  `yaml:"auto_update"`
}

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

// rateLimited checks if the given feedback message contains a rate limit error.
func rateLimited(feedbackMsg string) bool {
    targets := []string { "API rate limit exceeded", "API rate limit exceeded for" }
    return strings.Contains(feedbackMsg, targets[0]) || strings.Contains(feedbackMsg, targets[1])
}

// parseJsonFile parses a JSON file using a given JSONPath and returns the result.
func parseJsonFile(filePath string, jsonpathExpr string) (string, error) {
    content, readErr := os.ReadFile(filePath)
    if readErr != nil {
        h.Printc("Failed to read file!", h.ERROR, false)
        return "", fmt.Errorf("Failed to read file: %w", readErr)
    }

    value, err := parseJson(content, ".", "", jsonpathExpr)
    if err != nil {
        h.Printc("Failed to parse JSON!", h.ERROR, true)
        return "", fmt.Errorf("Failed to parse JSON: %w", err)
    }

    return value, nil
}

// parseJson parses JSON content and returns the matches of a given regex.
func parseJson(content []byte, regexStr string, regexRepl string, jsonpathExpr string) (string, error) {
    result := ""
    data := any (nil)

    // Compile regex from string
    regex, regexErr := regexp.Compile(regexStr)
    if regexErr != nil {
        h.Printc("Failed to parse regex!", h.ERROR, false)
        return result, fmt.Errorf("Failed to parse regex: %w", regexErr)
    }

    // Compile JSONPath from string
    jpath, jpathErr := json.CreatePath(jsonpathExpr)
    if jpathErr != nil {
        h.Printc("Failed to parse JSONPath!", h.ERROR, false)
        return result, fmt.Errorf("Failed to parse JSONPath: %w", jpathErr)
    }

    if err := jpath.Unmarshal(content, &data); err != nil {
        h.Printc("Failed to parse JSON with JSONPath!", h.ERROR, false)
        return result, fmt.Errorf("Failed to parse JSON with JSONPath: %w", err)
    }

    strValue, ok := data.(string)
    if !ok {
        println("JSONPath did not return a string value!")
        return result, fmt.Errorf("JSONPath did not return a string value: %T, expected string", data)
    }

    if regexRepl != "" {
        matches := regex.ReplaceAllString(strValue, regexRepl)
        if regex.MatchString(strValue) { result = matches }
    } else {
        matches := regex.ReplaceAllString(strValue, "$1")
        if regex.MatchString(strValue) { result = matches }
    }

    return result, nil
}

func parseYaml(content []byte, regexStr string, regexRepl string, yamlpathExpr string) (string, error) {
    result := ""
    data := any (nil)

    // Compile regex from string
    regex, regexErr := regexp.Compile(regexStr)
    if regexErr != nil {
        h.Printc("Failed to parse regex!", h.ERROR, false)
        return result, fmt.Errorf("Failed to parse regex: %w", regexErr)
    }

    // Compile yamlpath from string
    yamlpath, yamlpathErr := yaml.PathString(yamlpathExpr)
    if yamlpathErr != nil {
        h.Printc("Failed to parse YAMLPath!", h.ERROR, false)
        return result, fmt.Errorf("Failed to parse YAMLPath: %w", yamlpathErr)
    }

    contentReader := bytes.NewReader(content)
    if err := yamlpath.Read(contentReader, &data); err != nil {
        h.Printc("Failed to parse YAML with YAMLPath!", h.ERROR, false)
        return result, fmt.Errorf("Failed to parse YAML with YAMLPath: %w", err)
    }

    strValue, ok := data.(string)
    if !ok {
        println("YAMLPath did not return a string value!")
        return result, fmt.Errorf("YAMLPath did not return a string value: %T, expected string", data)
    }

    if regexRepl != "" {
        matches := regex.ReplaceAllString(strValue, regexRepl)
        if regex.MatchString(strValue) { result = matches }
    } else {
        matches := regex.ReplaceAllString(strValue, "$1")
        if regex.MatchString(strValue) { result = matches }
    }

    return result, nil
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
    filePath := filepath.Join(YUM_REPOS_DIR, baseName)

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
    args := []string { "copr", "enable", "-y", username + "/" + project }
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
