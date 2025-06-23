package cmd

import (
    // "bytes"
    "fmt"
    "os"
    "path/filepath"
    // "regexp"
    "strings"

    h "github.com/FlawlessCasual17/rpm-get/helpers"
    // "github.com/goccy/go-json"
    "github.com/goccy/go-yaml"
    "github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
    Use:   "info",
    Short: "Display information about a package",
    Long: "Display information about a package",
    Run: func(cmd *cobra.Command, args []string) {
        pkg := args[0]
        result, err := pkgInfo(pkg)
        if err != nil {
            h.Printc("Failed to get package information!", h.ERROR, false)
        }
        println(result)
    },
}

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
}

// // parseJsonFile parses a JSON file using a given JSONPath and returns the result.
// func parseJsonFile(filePath string, jsonpathExpr string) (string, error) {
//     content, readErr := os.ReadFile(filePath)
//     if readErr != nil {
//         h.Printc("Failed to read file!", h.ERROR, false)
//         return "", fmt.Errorf("Failed to read file: %w", readErr)
//     }
//
//     value, err := parseJson(content, ".", "", jsonpathExpr)
//     if err != nil {
//         h.Printc("Failed to parse JSON!", h.ERROR, true)
//         return "", fmt.Errorf("Failed to parse JSON: %w", err)
//     }
//
//     return value, nil
// }
//
// // parseJson parses JSON content and returns the matches of a given regex.
// func parseJson(content []byte, regexStr string, regexRepl string, jsonpathExpr string) (string, error) {
//     result := ""
//     data := any (nil)
//
//     // Compile regex from string
//     regex, regexErr := regexp.Compile(regexStr)
//     if regexErr != nil {
//         h.Printc("Failed to parse regex!", h.ERROR, false)
//         return result, fmt.Errorf("Failed to parse regex: %w", regexErr)
//     }
//
//     // Compile JSONPath from string
//     jpath, jpathErr := json.CreatePath(jsonpathExpr)
//     if jpathErr != nil {
//         h.Printc("Failed to parse JSONPath!", h.ERROR, false)
//         return result, fmt.Errorf("Failed to parse JSONPath: %w", jpathErr)
//     }
//
//     if err := jpath.Unmarshal(content, &data); err != nil {
//         h.Printc("Failed to parse JSON with JSONPath!", h.ERROR, false)
//         return result, fmt.Errorf("Failed to parse JSON with JSONPath: %w", err)
//     }
//
//     strValue, ok := data.(string)
//     if !ok {
//         println("JSONPath did not return a string value!")
//         return result, fmt.Errorf("JSONPath did not return a string value: %T, expected string", data)
//     }
//
//     if regexRepl != "" {
//         matches := regex.ReplaceAllString(strValue, regexRepl)
//         if regex.MatchString(strValue) { result = matches }
//     } else {
//         matches := regex.ReplaceAllString(strValue, "$1")
//         if regex.MatchString(strValue) { result = matches }
//     }
//
//     return result, nil
// }
//
// func parseYaml(content []byte, regexStr string, regexRepl string, yamlpathExpr string) (string, error) {
//     result := ""
//     data := any (nil)
//
//     // Compile regex from string
//     regex, regexErr := regexp.Compile(regexStr)
//     if regexErr != nil {
//         h.Printc("Failed to parse regex!", h.ERROR, false)
//         return result, fmt.Errorf("Failed to parse regex: %w", regexErr)
//     }
//
//     // Compile yamlpath from string
//     yamlpath, yamlpathErr := yaml.PathString(yamlpathExpr)
//     if yamlpathErr != nil {
//         h.Printc("Failed to parse YAMLPath!", h.ERROR, false)
//         return result, fmt.Errorf("Failed to parse YAMLPath: %w", yamlpathErr)
//     }
//
//     contentReader := bytes.NewReader(content)
//     if err := yamlpath.Read(contentReader, &data); err != nil {
//         h.Printc("Failed to parse YAML with YAMLPath!", h.ERROR, false)
//         return result, fmt.Errorf("Failed to parse YAML with YAMLPath: %w", err)
//     }
//
//     strValue, ok := data.(string)
//     if !ok {
//         println("YAMLPath did not return a string value!")
//         return result, fmt.Errorf("YAMLPath did not return a string value: %T, expected string", data)
//     }
//
//     if regexRepl != "" {
//         matches := regex.ReplaceAllString(strValue, regexRepl)
//         if regex.MatchString(strValue) { result = matches }
//     } else {
//         matches := regex.ReplaceAllString(strValue, "$1")
//         if regex.MatchString(strValue) { result = matches }
//     }
//
//     return result, nil
// }

func init() {
    rootCmd.AddCommand(infoCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // infoCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// PkgInfo returns the package information for the given package.
func pkgInfo(pkg string) (string, error) {
    result := ""
    data := Pkg {}
    filePath := filepath.Join(DataDir, pkg + ".yaml")

    content, readErr := os.ReadFile(filePath)
    if readErr != nil {
        h.Printc("Failed to read file!", h.ERROR, false)
        return result, fmt.Errorf("Failed to read file: %w", readErr)
    }

    if err := yaml.Unmarshal(content, &data); err != nil {
        h.Printc("Failed to unmarshal file!", h.ERROR, false)
        return result, fmt.Errorf("Failed to unmarshal file: %w", err)
    }

    result = fmt.Sprintf(`
        Supported OS: %s
        Version: %s
        Name: %s
        License: %s
        Homepage: %s
        Description: %s
        Notes: %s
        Pkg Arches: %s
        `,
        strings.Join(data.supported_os, ", "),
        data.version,
        data.name,
        data.license.licenseString,
        data.homepage,
        data.description,
        data.notes,
        strings.Join(data.pkg_arches, ", "),
    )

    return result, nil
}
