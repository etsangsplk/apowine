package credential

func NewGithubCreds(clientID, clientSecret, redirectURI string) *GithubCreds {
	return &GithubCreds{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

func NewGoogleCreds(clientID, clientSecret, redirectURI, refreshToken string) *GoogleCreds {

	return &GoogleCreds{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		refreshToken: refreshToken,
	}
}

func (git *GithubCreds) GetGithubCreds() (string, string, string) {

	return git.clientID, git.clientSecret, git.redirectURI
}

func (goo *GoogleCreds) GetGoogleCreds() (string, string, string, string) {

	return goo.clientID, goo.clientSecret, goo.redirectURI, goo.refreshToken
}
