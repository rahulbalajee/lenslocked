package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"golang.org/x/oauth2"
)

type OAuth struct {
	ProviderConfigs map[string]*oauth2.Config
}

func (oa OAuth) Connect(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	provider = strings.ToLower(provider)

	config, ok := oa.ProviderConfigs[provider]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}

	verifier := oauth2.GenerateVerifier()

	state := csrf.Token(r)
	setCookie(w, "oauth_state", state)
	setCookie(w, "oauth_verifier", verifier)

	url := config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("token_access_type", "offline"),
		oauth2.SetAuthURLParam("redirect_uri", redirectURI(r, provider)),
		oauth2.S256ChallengeOption(verifier),
	)

	http.Redirect(w, r, url, http.StatusFound)
}

func (oa OAuth) Callback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	provider = strings.ToLower(provider)

	config, ok := oa.ProviderConfigs[provider]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}

	state := r.FormValue("state")
	cookieState, err := readCookie(r, "oauth_state")
	if err != nil || cookieState != state {
		if err != nil {
			fmt.Println(err)
		}
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	deleteCookie(w, "oauth_state")

	verifier, err := readCookie(r, "oauth_verifier")
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	deleteCookie(w, "oauth_verifier")

	code := r.FormValue("code")
	token, err := config.Exchange(
		r.Context(),
		code,
		oauth2.SetAuthURLParam("redirect_uri", redirectURI(r, provider)),
		oauth2.VerifierOption(verifier),
	)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	client := config.Client(r.Context(), token)
	resp, err := client.Post(
		"https://api.dropboxapi.com/2/files/list_folder",
		"application/json",
		strings.NewReader(`{"path": ""}`),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var pretty bytes.Buffer
	err = json.Indent(&pretty, body, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	pretty.WriteTo(w)
}

func redirectURI(r *http.Request, provider string) string {
	if r.Host == "localhost:3000" {
		return fmt.Sprintf("http://localhost:3000/oauth/%s/callback", provider)
	}

	return fmt.Sprintf("http://lenslocked.bc-noc-dev.com/oauth/%s/callback", provider)
}
