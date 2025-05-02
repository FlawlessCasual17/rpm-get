package cmd

import (
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"

    h "github.com/FlawlessCasual17/rpm-get/helpers"
    "github.com/samber/lo"
    "github.com/schollz/progressbar/v3"
    "github.com/spf13/cobra"
)

// VERSION is the current version of _rpm-get_.
const VERSION string = "0.0.1"

var (
    // private variables


    // public variables

    // CacheDir is the directory where rpm-get will
    // cache JSON files from **_GitHub_**/**_GitLab_**.
    // As well as downloaded packages.
    CacheDir = filepath.Join(getEnv("HOME"), ".cache/rpm-get")
    Project = ""
    RelType = ""
    Creator = ""
    ProjectId = ""
    // UserAgent is the user agent string used for HTTP requests.
    UserAgent = fmt.Sprintf(
        "Mozilla/5.0 (X11; Linux %s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
        HOST_CPU)
    GhHeaderAuth = fmt.Sprintf("Bearer %s", getEnv("GITHUB_TOKEN"))
    GlHeaderAuth = getEnv("GITLAB_TOKEN")
)

const ETC_DIR string = "/etc/rpm-get"

// HOST_CPU is the host CPU architecture.
const HOST_CPU string = runtime.GOARCH

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

// getEnv returns the value of the environment variable. Empty string if not found.
func getEnv(key string) string {
    v, ok := os.LookupEnv(key)
    return lo.Ternary(ok, v, "")
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
func which(cmd string) string {
    result, err := exec.LookPath(cmd)
    msg := fmt.Sprintf("Command `%s` not found in PATH. Exiting...", cmd)

    if err != nil { h.Printc(msg, h.FATAL, false) }
    return result
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

// spellcheck: ignore
// # Gets the releases from either GitHub or GitLab
// def get_releases
//   c_file_path = "#{$CACHE_DIR}/#{$app}_cache.json"
//   url = ''
//   headers = {}
//   feedback_msg = nil
//
//   case $type
//   when 'github'
//     base_url = 'https://api.github.com/repos'
//     url = "#{base_url}/#{$creator}/#{$app}/releases/latest"
//     headers = {
//       'User-Agent': $USER_AGENT,
//       'Authorization': $gh_header_auth
//     }
//     feedback_msg = parse_json(c_file_path, ['message'])
//   when 'gitlab'
//     base_url = 'https://gitlab.com/api/v4/projects'
//     url = "#{base_url}/#{$project_id}/releases/permalink/latest"
//     headers = {
//       'User-Agent': $USER_AGENT,
//       'PRIVATE-TOKEN': $gl_header_auth
//     }
//     feedback_msg = File.read(c_file_path)
//   end
//
//   # Ensure cache directory exists
//   create_cache_dir() unless Dir.exist?($CACHE_DIR)
//   # If cache file exists, skip
//   return if File.exist?(c_file_path)
//
//   # Notify the user
//   printc "Updating #{c_file_path}", 'info'
//   printc "Downloading JSON cache of #{$app} to #{c_file_path}", 'progress'
//
//   # Try downloading, or raise an error
//   begin
//     Down.download(
//       url,
//       destination: c_file_path,
//       headers: headers,
//       # @param [Integer] content_length
//       content_length_proc: lambda do |content_length|
//         # If `content_length` is null, change the format
//         content_length.nil? &&
//           $format = $format.sub('TOTAL::total_byte', '').sub('/:total', '')
//
//         bar = TTY::ProgressBar.new($format, total: content_length, head: '>')
//         bar.resize(50)
//         $progress_bar = bar
//       end,
//       # @param [Float] progress
//       progress_proc: lambda do |progress|
//         raise TypeError if progress.nil?
//
//         $progress_bar&.advance(progress - $progress_bar.current)
//       end
//     )
//
//     printc 'Download complete!', 'progress', true
//   rescue StandardError => e
//     printc "Failed to update cache for #{$app}", 'error', true
//     printc e.detailed_message, 'error'
//   end
//
//   # If the GitHub/GitLab API rate limit is exceeded, tell the user
//   return unless rate_limited?(feedback_msg)
//
//   printc 'API rate limit exceeded!', 'warn'
//   printc "Deleting #{c_file_path}", 'info'
//   FileUtils.remove_file(c_file_path)
// end
// ⬆️ Convert this from Ruby to Go

// getReleases retrieves the list of releases from the GitHub/GitLab API.
func getReleases() {
    filePath := fmt.Sprintf("%s_cache.json", Project)
    cacheFilePath := filepath.Join(CacheDir, filePath)
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

    if _, e := os.Stat(CacheDir); e != nil && os.IsNotExist(e) { createCacheDir() }

    // NOTE: `//nolint:all` is used to suppress annoying linter warnings/errors.

    lo.TryCatch(func() error {
        resp, _ := http.DefaultClient.Do(request)
        //nolint:all
        defer resp.Body.Close()

        file, _ := os.OpenFile(cacheFilePath, os.O_CREATE|os.O_WRONLY, 0644)
        //nolint:all
        defer file.Close()

        // NOTE: Might switch to "github.com/cheggaaa/pb" if this doesn't meet my needs.
        bar := progressbar.DefaultBytes(resp.ContentLength, "Downloading...")

        h.Printc(fmt.Sprintf("Downloading %s to %s", url, cacheFilePath), h.INFO, true)
        //nolint:all
        io.Copy(io.MultiWriter(file, bar), resp.Body)

        return nil // No need to return an error here
    }, func() {
        h.Printc("Unable to create cache file!", h.ERROR, false)
        os.Exit(h.ERROR_EXIT_CODE)
    })

    if rateLimited(feedbackMsg) {
        h.Printc("API rate limit exceeded!", h.WARNING, false)
        h.Printc("Deleting cache file...", h.INFO, true)
        err := os.Remove(cacheFilePath)
        if err != nil { h.Printc(err.Error(), h.ERROR, false) }
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
