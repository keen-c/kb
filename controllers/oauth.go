package controllers

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func Oauth() {
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT"),
			os.Getenv("GOOGLE_SECRET"),
			"http://localhost:8080/auth/callback?provider=google",
		),
	)
	key := os.Getenv("SESSION_STORE_COOKIE") // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30                     // 30 days
	isProd := false                          // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store
}
