package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"

	"github.com/aporeto-inc/apowine/source/frontend-ui/client/internal/credential"
	githubAuth "github.com/aporeto-inc/apowine/source/frontend-ui/client/internal/github"
	googleAuth "github.com/aporeto-inc/apowine/source/frontend-ui/client/internal/google"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Auth struct {
	store        *Cookie
	googleHandle *googleAuth.Google
	githubHandle *githubAuth.Github
}

func NewAuth(googleCreds *credential.GoogleCreds, githubCreds *credential.GithubCreds) *Auth {

	return &Auth{
		store:        newCookie(),
		googleHandle: googleAuth.NewGoogleHandler(googleCreds),
		githubHandle: githubAuth.NewGithubHandler(githubCreds),
	}
}

// Login is
func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("loginDomain")
	if domain == "google" {
		OAuthConfig := a.googleHandle.GetGoogleConfig()
		url := OAuthConfig.AuthCodeURL(base64.URLEncoding.EncodeToString([]byte("apowine")), oauth2.AccessTypeOffline, oauth2.ApprovalForce)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	} else if domain == "github" {
		OAuthConfig := a.githubHandle.GetGithubConfig()
		url := OAuthConfig.AuthCodeURL(base64.URLEncoding.EncodeToString([]byte("apowine")), oauth2.AccessTypeOnline)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (a *Auth) GetCookie() *Cookie {

	return a.store
}

func (a *Auth) GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := a.store.GetCookieStore().Get(r, "githubSessions")
	state, _ := base64.URLEncoding.DecodeString(r.FormValue("state"))

	OAuthConf := a.githubHandle.GetGithubConfig()

	if string(state) != "apowine" {
		zap.L().Warn(fmt.Sprintf("invalid oauth state, expected '%s', got '%s'\n", "apowine", state))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	fmt.Println("CODE", code)
	token, err := OAuthConf.Exchange(context.TODO(), code)

	if err != nil {
		zap.L().Warn("oauthConf.Exchange() failed", zap.Error(err))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := OAuthConf.Client(context.TODO(), token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(context.TODO(), "")
	if err != nil {
		zap.L().Warn("client.Users.Get() failed", zap.Error(err))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	zap.L().Debug(fmt.Sprintf("Logged in as GitHub user: %s\n", *user.Login))

	session.Values["authenticated"] = true

	if err := session.Save(r, w); err != nil {
		zap.L().Error("Error in saving the session", zap.Error(err))
	}

	if redirectURL, ok := session.Values["redirectURL"]; ok {
		http.Redirect(w, r, redirectURL.(string), http.StatusFound)
	} else {
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func (a *Auth) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := a.store.GetCookieStore().Get(r, "googleSessions")
	state, _ := base64.URLEncoding.DecodeString(r.FormValue("state"))

	OAuthConf := a.googleHandle.GetGoogleConfig()

	if string(state) != "apowine" {
		zap.L().Warn(fmt.Sprintf("invalid oauth state, expected '%s', got '%s'\n", "apowine", state))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := OAuthConf.Exchange(context.TODO(), code)

	idToken, _ := a.googleHandle.RequestIDToken()

	fmt.Println("TOKEN", idToken.(string))

	_ = a.validateAndIssueJWTUsingMidgard(idToken.(string))

	if err != nil {
		zap.L().Warn("oauthConf.Exchange() failed", zap.Error(err))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := OAuthConf.Client(context.TODO(), token)
	client, err := oauthClient.Get("https://www.googleapis.com/oauth2/v3/userprofile")
	if err != nil {
		zap.L().Warn("client.Users.Get() failed", zap.Error(err))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer client.Body.Close()
	cdata, _ := ioutil.ReadAll(client.Body)
	fmt.Println("Email body: ", string(cdata))

	session.Values["authenticated"] = true
	session.Values["id_token"] = idToken.(string)

	if err := session.Save(r, w); err != nil {
		zap.L().Error("Error in saving the session", zap.Error(err))
	}

	if redirectURL, ok := session.Values["redirectURL"]; ok {
		http.Redirect(w, r, redirectURL.(string), http.StatusFound)
	} else {
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func (a *Auth) validateAndIssueJWTUsingMidgard(idToken string) string {
	//TODO: Should use midgard package to issue JWT using google tokens

	return ""
}
