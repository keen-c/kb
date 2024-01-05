package controllers

import (
	"os"
	"github.com/markbates/goth"
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
	
}
