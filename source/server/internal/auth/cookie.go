package auth

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type Cookie struct {
	cookieStore *sessions.CookieStore
}

func newCookie() *Cookie {
	return &Cookie{
		cookieStore: sessions.NewCookieStore([]byte("apowine")),
	}
}

func (c *Cookie) GetCookieStore() *sessions.CookieStore {

	return c.cookieStore
}

func (c *Cookie) EmptyCookieStore() *sessions.CookieStore {

	return sessions.NewCookieStore(securecookie.GenerateRandomKey(5))
}
