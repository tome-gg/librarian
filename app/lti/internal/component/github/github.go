package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/v51/github"
	"gopkg.in/yaml.v2"
)

// ParsedYAML ...
type ParsedYAML map[string]interface{}

var (
	githubAuthBaseURL = "https://github.com/login/oauth/authorize"
	clientID          = os.Getenv("GITHUB_CLIENT_ID")
	scope             = "repo" // Adjust the scope based on your needs
)

// BuildGitHubAuthURL based on state
func BuildGitHubAuthURL(state string) string {
	authURL, _ := url.Parse(githubAuthBaseURL)

	query := authURL.Query()
	query.Set("client_id", clientID)
	query.Set("scope", scope)
	query.Set("state", state)

	authURL.RawQuery = query.Encode()

	return authURL.String()
}



var (
	githubTokenURL  = "https://github.com/login/oauth/access_token"
	clientSecret    = os.Getenv("GITHUB_CLIENT_SECRET")
)

type githubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// ExchangeGitHubCodeForAccessToken ...
func ExchangeGitHubCodeForAccessToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)

	req, err := http.NewRequest("POST", githubTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse githubAccessTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", err
	}

	if tokenResponse.AccessToken == "" {
		return "", errors.New("failed to obtain access token")
	}

	return tokenResponse.AccessToken, nil
}

// GetAndParseYAMLFile ...
func GetAndParseYAMLFile(client *github.Client, owner, repo, filePath string) (ParsedYAML, error) {
	ctx := context.Background()

	// Get the file content from the specified repository and file path
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get file content: %v", err)
	}

	// Ensure the content is available and not a directory
	if fileContent == nil || fileContent.Content == nil || fileContent.GetEncoding() != "base64" {
		return nil, fmt.Errorf("invalid file content")
	}

	// Decode the base64-encoded file content
	decodedContent, err := fileContent.GetContent()
	if err != nil {
		return nil, fmt.Errorf("failed to decode file content: %v", err)
	}

	// Parse the YAML content
	var parsedContent ParsedYAML
	if err := yaml.Unmarshal([]byte(decodedContent), &parsedContent); err != nil {
		return nil, fmt.Errorf("failed to parse YAML content: %v", err)
	}

	return parsedContent, nil
}
