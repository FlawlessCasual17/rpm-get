package main

import (
    "encoding/json"
    "fmt"
    "regexp"

    "github.com/PaesslerAG/jsonpath"
)

var content = []byte(`
[
  {
    "url": "https://api.github.com/repos/bitwarden/clients/releases/216016753",
    "assets_url": "https://api.github.com/repos/bitwarden/clients/releases/216016753/assets",
    "upload_url": "https://uploads.github.com/repos/bitwarden/clients/releases/216016753/assets{?name,label}",
    "html_url": "https://github.com/bitwarden/clients/releases/tag/desktop-mac-v2025.4.2",
    "id": 216016753,
    "author": {
      "login": "github-actions[bot]",
      "id": 41898282,
      "node_id": "MDM6Qm90NDE4OTgyODI=",
      "avatar_url": "https://avatars.githubusercontent.com/in/15368?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/github-actions%5Bbot%5D",
      "html_url": "https://github.com/apps/github-actions",
      "followers_url": "https://api.github.com/users/github-actions%5Bbot%5D/followers",
      "following_url": "https://api.github.com/users/github-actions%5Bbot%5D/following{/other_user}",
      "gists_url": "https://api.github.com/users/github-actions%5Bbot%5D/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/github-actions%5Bbot%5D/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/github-actions%5Bbot%5D/subscriptions",
      "organizations_url": "https://api.github.com/users/github-actions%5Bbot%5D/orgs",
      "repos_url": "https://api.github.com/users/github-actions%5Bbot%5D/repos",
      "events_url": "https://api.github.com/users/github-actions%5Bbot%5D/events{/privacy}",
      "received_events_url": "https://api.github.com/users/github-actions%5Bbot%5D/received_events",
      "type": "Bot",
      "user_view_type": "public",
      "site_admin": false
    },
    "node_id": "RE_kwDOAzDwU84M4Cdx",
    "tag_name": "desktop-mac-v2025.4.2",
    "target_commitish": "f01e878ec1251e91f4e1ce9b1878e8e24486e6cf",
    "name": "Desktop v2025.4.2",
    "draft": false,
    "prerelease": false,
    "created_at": "2025-05-01T15:26:36Z",
    "published_at": "2025-05-01T16:37:19Z"
  }
]
`)

func main() {
    result, err := parseJson(content, "Desktop v([\\d.]+)", "$[0].name")

    if err != nil {
        // If an error occurred, print the error message and exit or handle it
        print("Error: " + err.Error())
        // os.Exit(1) // You might want to exit if the error is critical
        return // Stop execution in main if there's an error
    }

    // If there was no error, print the successful result
    println(result)
}

// parseJson parses JSON content and returns the matches of a given regex.
func parseJson(content []byte, regexStr string, jsonpathExpr string) (string, error) {
    result := ""
    data := any (nil)

    // Compile regex from string
    regex, regexErr := regexp.Compile(regexStr)
    if regexErr != nil {
        println("Failed to parse regex!")
        return result, fmt.Errorf("Failed to parse regex: %w", regexErr)
    }

    if err := json.Unmarshal(content, &data); err != nil {
        println("Failed to unmarshal JSON!")
        return result, fmt.Errorf("Failed to unmarshal JSON: %w", err)
    }

    value, err := jsonpath.Get(jsonpathExpr, data)
    if err != nil {
        println("Failed to parse JSON!")
        return result, fmt.Errorf("Failed to parse JSON: %w", err)
    }

    strValue, ok := value.(string)
    if !ok {
        println("JSONPath did not return a string value!")
        return result, fmt.Errorf("JSONPath did not return a string value: %T, expected string", value)
    }

    matches := regex.ReplaceAllString(strValue, "$1")
    if regex.MatchString(strValue) { result = matches }

    return result, nil
}
