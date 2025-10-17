package controllers

import (
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

	url := config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("token_access_type", "offline"),
		oauth2.SetAuthURLParam("redirect_uri", "http://localhost:3000/oauth/dropbox/callback"),
		oauth2.S256ChallengeOption(verifier),
	)

	http.Redirect(w, r, url, http.StatusFound)
}
