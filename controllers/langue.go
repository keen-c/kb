package controllers

import (
	"fmt"
	"kb/kb/models"
	"log"
	"net/http"
)

type LangueController struct {
	Lm *models.LanguageManager
}

func (lc *LangueController) HandleGetThemeByLanguage(w http.ResponseWriter, r *http.Request) {
	l, err := lc.Lm.QueryUserLangue(r.Context())
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(l)
}
