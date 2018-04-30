package credential

type GithubCreds struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

type GoogleCreds struct {
	clientID     string
	clientSecret string
	redirectURI  string
	refreshToken string
}
