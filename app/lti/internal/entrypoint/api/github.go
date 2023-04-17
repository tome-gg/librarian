package api

import (
	"context"
	"net/http"

	"github.com/google/go-github/v51/github"
	github2 "github.com/tome-gg/librarian/app/lti/internal/component/github"
	"golang.org/x/oauth2"
)

func githubOAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the state parameter matches the one stored in the session
	session, _ := sessionStore.Get(r, "session-name")
	state := r.URL.Query().Get("state")
	storedState, _ := session.Values["github_auth_state"].(string)

	if state != storedState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for an access token
	code := r.URL.Query().Get("code")
	accessToken, err := github2.ExchangeGitHubCodeForAccessToken(code)
	if err != nil {
		http.Error(w, "Failed to obtain GitHub access token", http.StatusInternalServerError)
		return
	}

	// Store the access token in the user's session
	session.Values["github_access_token"] = accessToken
	session.Save(r, w)

	// Redirect the user back to the LTI launch endpoint
	// This will call the handleGitHubIntegration function with the access token now available
	http.Redirect(w, r, "/lti/launch", http.StatusFound)
}

func getGithubRepo() ( owner string, repo string) {
	return "", ""
}


func handleReadGitHubRepo(r *http.Request) (github2.ParsedYAML, error) {
	// Get the user's GitHub access token (either from the session or via OAuth)
	// For example:
	// accessToken := session.Values["github_access_token"].(string)

	// Replace this with the actual GitHub access token retrieval
	accessToken := "your-access-token"

	// Create a new GitHub client using the user's access token
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	const filepath = "training/dsu.yaml"

	owner, repo := getGithubRepo()

	growthData, err := github2.GetAndParseYAMLFile(client, owner, repo, filepath)

	if err != nil {
		return nil, err
	}

	return growthData, nil
}