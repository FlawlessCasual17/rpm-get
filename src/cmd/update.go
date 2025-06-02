package cmd

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"

    h "github.com/FlawlessCasual17/rpm-get/helpers"
    "github.com/goccy/go-json"
    "github.com/samber/lo"
    "github.com/schollz/progressbar/v3"
    "github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
    Use:   "update",
    Short: "A brief description of your command",
    Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("update called")
    },
}

func init() {
    rootCmd.AddCommand(updateCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // updateCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// getUpdates checks for updates to the packages list.
func getUpdates() {
    success := false
    url := PKGS_REPO + "/raw/refs/heads/master/packages-list.json"
    tmpFilePath := filepath.Join(ConfigDir, "packages-list.json.tmp")
    filePath := filepath.Join(ConfigDir, "packages-list.json")
    // Download packages-list.json
    lo.TryCatch(func() error { // try
        request, _ := http.NewRequest("GET", url, nil)
        request.Header.Set("User-Agent", UserAgent)
        resp, respErr := http.DefaultClient.Do(request)
        if respErr != nil {
            h.Printc("Request failed!", h.ERROR, false)
            return fmt.Errorf("Request failed: %w", respErr)
        }
        //nolint:errcheck
        defer resp.Body.Close()

        file, _ := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY, 0644)
        //nolint:errcheck
        defer file.Close()

        bar := progressbar.DefaultBytes(resp.ContentLength, "Updating packages list...")
        //nolint:errcheck
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
            success = true
        }
    } else { h.Printc("Packages list is already up to date!", h.INFO, true) }

    if success {
        data, _ := os.ReadFile(filePath)

        pkgs := []string {}
        if err := json.Unmarshal(data, &pkgs); err != nil {
            h.Printc("Failed to unmarshal packages list!", h.ERROR, false)
            os.Exit(h.ERROR_EXIT_CODE)
        }

        if err := getPkgManifests(pkgs); err != nil {
            h.Printc("Failed to download package manifests!", h.ERROR, false)
            os.Exit(h.ERROR_EXIT_CODE)
        }
    }
}

// getPkgManifests retrieves the manifests for the given packages.
func getPkgManifests(pkgs []string) error {
    downloadError := error (nil)

    h.Printc("Downloading package manifests...", h.INFO, false)
    for _, pkg := range pkgs {
        if downloadError != nil { break }

        url := PKGS_REPO + fmt.Sprintf("/raw/refs/heads/master/manifests/%s.json", pkg)
        baseName := strings.Split(url, "/")[-0]
        filePath := filepath.Join(DataDir, baseName)

        lo.TryCatch(func() error { // try
            request, _ := http.NewRequest("GET", url, nil)
            resp, respErr := http.DefaultClient.Do(request)
            if respErr != nil {
                h.Printc("Request failed!", h.ERROR, false)
                downloadError = fmt.Errorf("Request failed: %w", respErr)
                return downloadError
            }
            //nolint:errcheck
            defer resp.Body.Close()

            file, _ := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
            //nolint:errcheck
            defer file.Close()

            bar := progressbar.DefaultBytes(resp.ContentLength, pkg)

            if _, err := io.Copy(io.MultiWriter(file, bar), resp.Body); err != nil {
                msg := fmt.Sprint("Failed to download package manifest for " + pkg)
                h.Printc(msg, h.ERROR, false)
                downloadError = fmt.Errorf("%s: %w", msg, err)
                return downloadError
            }

            return nil
        }, func() { // catch
            h.Printc("Unknown error occurred!", h.ERROR, false)
            downloadError = fmt.Errorf("Unknown error occurred")
        })
    }

    return downloadError
}
