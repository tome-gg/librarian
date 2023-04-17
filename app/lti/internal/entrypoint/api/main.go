package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tome-gg/librarian/app/lti/internal/component/github"
)

var srv *http.Server

// Stop stops the web server.
func Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %v", err)
	}
	return nil
}

// Start starts the web server.
func Start(port int) error {
	r := mux.NewRouter()

	r.HandleFunc("/lti/launch", ltiLaunchHandler).Methods("POST")

	// Add additional routes for your application here

	http.Handle("/", r)
	
	addr := fmt.Sprintf(":%d", port)
	srv = &http.Server{Addr: addr, Handler: r}

	fmt.Printf("Starting server on %s...\n", addr)

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("server error: %v", err)
	}

	return nil
}

func ltiLaunchHandler(w http.ResponseWriter, r *http.Request) {
	// Determine the LTI version based on request parameters
	ltiVersion := r.FormValue("lti_version")

	switch ltiVersion {
	case "LTI-1p0", "LTI-1p1":
		// Validate LTI 1.1 launch request
		if err := validateLTI11Request(r); err != nil {
			http.Error(w, "Invalid LTI 1.1 launch request", http.StatusBadRequest)
			return
		}
	case "LTI-1p3":
		// Validate LTI 1.3 launch request
		if err := validateLTI13Request(r); err != nil {
			http.Error(w, "Invalid LTI 1.3 launch request", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Invalid LTI version", http.StatusBadRequest)
		return
	}

	// Set up user session and SSO
	session, _ := sessionStore.Get(r, "user-session")
	session.Values["authenticated"] = true
	session.Save(r, w)

	// Check if the user has a GitHub access token stored in the session
	githubAccessToken, ok := session.Values["github_access_token"].(string)

	// Redirect to GitHub authorization page if the access token is not available
	if !ok || githubAccessToken == "" {
		// Set up a state parameter to prevent CSRF attacks
		state := generateState()
		session.Values["github_auth_state"] = state
		session.Save(r, w)

		// Build the GitHub authorization URL
		githubAuthURL := github.BuildGitHubAuthURL(state)

		// Redirect the user to the GitHub authorization page
		http.Redirect(w, r, githubAuthURL, http.StatusFound)
		return
	}

	// Call the GitHub integration function
	data, err := handleReadGitHubRepo(r); 

	if err != nil {
		http.Error(w, "Failed to integrate with GitHub", http.StatusInternalServerError)
		return
	}

	fmt.Printf("%+v", data)

	// Render resource within LMS
	// You can use your Next.js frontend here
	fmt.Fprint(w, "Resource rendered within LMS")
}

const stateLength = 32

func generateState() string {
	randomBytes := make([]byte, stateLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// You may want to handle this error differently depending on your use case
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(randomBytes)
}