package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/chriswhitcombe/rbac"
	"github.com/google/go-github/github"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type Auth struct {
	rbacMapper *rbac.RoleMapper
	store      *Cookie
	request    *http.Request
}

func NewAuth() *Auth {
	return &Auth{
		rbacMapper: rbac.NewRoleMapper(),
		store:      newCookie(),
		request:    &http.Request{},
	}
}

//Login is
func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {

	OAuthConf := &oauth2.Config{ClientID: "560f2688c98130ea6234",
		ClientSecret: "1c914abf80c3edd5b93cff3e47ab45afe96928bf",
		Endpoint:     githuboauth.Endpoint}

	url := OAuthConf.AuthCodeURL(base64.URLEncoding.EncodeToString([]byte("apowine")), oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (a *Auth) GetCookie() *Cookie {

	return a.store
}

func (a *Auth) GetRequest() *http.Request {

	return a.request
}

func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {

	fmt.Println(gcontext.Get(r, "req"))
	store, err := a.store.GetCookieStore().Get(gcontext.Get(a.request, "req").(*http.Request), "sessions")

	if err != nil {
		fmt.Println("Error retrieving cookies")
	} else if store != nil {
		store.Values["authenticated"] = false
		if err := store.Save(r, w); err != nil {
			zap.L().Error("Error in saving the session", zap.Error(err))
		}
		a.store = &Cookie{
			cookieStore: sessions.NewCookieStore(securecookie.GenerateRandomKey(5)),
		}
		a.request = &http.Request{}
	}
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func (a *Auth) GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := a.store.GetCookieStore().Get(r, "sessions")
	state, _ := base64.URLEncoding.DecodeString(r.FormValue("state"))

	OAuthConf := &oauth2.Config{ClientID: "560f2688c98130ea6234",
		ClientSecret: "1c914abf80c3edd5b93cff3e47ab45afe96928bf",
		Endpoint:     githuboauth.Endpoint}

	if string(state) != "apowine" {
		zap.L().Warn(fmt.Sprintf("invalid oauth state, expected '%s', got '%s'\n", "apowine", state))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
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
	gcontext.Set(a.request, "req", r)

	if err := session.Save(r, w); err != nil {
		zap.L().Error("Error in saving the session", zap.Error(err))
	}

	if redirectURL, ok := session.Values["redirectURL"]; ok {
		http.Redirect(w, r, redirectURL.(string), http.StatusFound)
	} else {
		if _, err := w.Write([]byte("Login Successful. Try accessing the URL for the scenario log")); err != nil {
			zap.L().Error("Failed to send Login Successful msg", zap.Error(err))
		}
	}
}

func (a *Auth) CreateAndPoliceRBACHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fmt.Print("RBAC")
		a.addMethodMapping()
		return http.HandlerFunc(nil)
	}
}

func (a *Auth) addMethodMapping() {
	a.rbacMapper.AddMethodMapping("/random", http.MethodGet, []string{"admin"})

	a.rbacMapper.AddMethodMapping("/beer", http.MethodGet, []string{"admin"})
	a.rbacMapper.AddMethodMapping("/beer", http.MethodPost, []string{"admin"})
	a.rbacMapper.AddMethodMapping("/beer", http.MethodPut, []string{"admin"})
	a.rbacMapper.AddMethodMapping("/beer/random", http.MethodGet, []string{"admin"})
	a.rbacMapper.AddMethodMapping("/beer/{id}", http.MethodGet, []string{"admin"})
	a.rbacMapper.AddMethodMapping("/beer/{id}", http.MethodDelete, []string{"admin"})

	a.rbacMapper.AddMethodMapping("/wine", http.MethodGet, []string{"admin"})
	a.rbacMapper.AddMethodMapping("/wine", http.MethodPost, []string{"admin"})
	a.rbacMapper.AddMethodMapping("/wine", http.MethodPut, []string{"admin"})
	a.rbacMapper.AddMethodMapping("/wine/random", http.MethodGet, []string{"admin"})
	a.rbacMapper.AddMethodMapping("/wine/{id}", http.MethodGet, []string{"admin"})
	a.rbacMapper.AddMethodMapping("/wine/{id}", http.MethodDelete, []string{"admin"})
}
