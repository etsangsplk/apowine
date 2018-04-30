package github

import (
	"github.com/aporeto-inc/apowine/source/frontend-ui/client/internal/credential"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Github struct {
	githubCreds *credential.GithubCreds
}

func NewGithubHandler(githubCreds *credential.GithubCreds) *Github {

	return &Github{
		githubCreds: githubCreds,
	}
}

func (g *Github) GetGithubConfig() *oauth2.Config {

	clientID, clientSecret, _ := g.githubCreds.GetGithubCreds()

	OAuthConf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     github.Endpoint,
	}

	return OAuthConf
}
