package google

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aporeto-inc/apowine/source/frontend-ui/client/internal/credential"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Google struct {
	googleCreds *credential.GoogleCreds
	idToken     string
}

func NewGoogleHandler(googleCreds *credential.GoogleCreds) *Google {

	return &Google{
		googleCreds: googleCreds,
	}
}

func (g *Google) GetGoogleConfig() *oauth2.Config {

	clientID, clientSecret, redirectURI, _ := g.googleCreds.GetGoogleCreds()

	OAuthConf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
			"https://www.googleapis.com/auth/bigquery",
			"https://www.googleapis.com/auth/userinfo.email",
			"openid",
		},
		RedirectURL: redirectURI,
		Endpoint:    google.Endpoint}

	return OAuthConf
}

func (g *Google) RequestIDToken() (interface{}, error) {
	clientID, clientSecret, _, refreshToken := g.googleCreds.GetGoogleCreds()

	apiUrl := "https://www.googleapis.com"
	resource := "/oauth2/v4/token"
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	clients := &http.Client{}
	rTok, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	rTok.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resptok, _ := clients.Do(rTok)

	body, err := ioutil.ReadAll(resptok.Body)
	if err != nil {
		return "", err
	}
	var idToken map[string]interface{}
	json.Unmarshal(body, &idToken)
	token := idToken["id_token"]

	return token, nil
}
