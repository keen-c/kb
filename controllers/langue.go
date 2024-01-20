package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"kb/kb/models"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type LangueController struct {
	Lm *models.LanguageManager
}
type Game struct {
	Theme string          `json:"theme"`
	Game  models.Question `json:"game"`
}

func (lc *LangueController) HandleGetGame(w http.ResponseWriter, r *http.Request) {
	th_id, err := lc.Lm.GetCurrentTheme(r.Context())
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	th_name, err := lc.Lm.QueryNameTheme(r.Context(), th_id)
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	q, err := lc.Lm.GetUserCurrentQueux(r.Context())
	if errors.Is(err, sql.ErrNoRows) {
		five, err := lc.Lm.GetFiveWordFromTheCurrentTheme(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(five)
		if err != nil {
			log.Printf("%s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = lc.Lm.SeUserCurrentQueux(r.Context(), b)
		if err != nil {
			log.Printf("set user queux %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		games := Game{
			Theme: th_name,
			Game:  (*five)[0],
		}
		err = json.NewEncoder(w).Encode(games)
		if err != nil {
			log.Printf("error on encode : %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	games := Game{
		Theme: th_name,
		Game:  (*q)[0],
	}
	err = json.NewEncoder(w).Encode(games)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func (lc *LangueController) HandlePostGame(w http.ResponseWriter, r *http.Request) {
	var answer models.Answer
	err := decodeJSONBody(w, r, &answer)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	q, err := lc.Lm.GetUserCurrentQueux(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := lc.Lm.CheckTheAnswer(r.Context(), answer.Answer, (*q)[0].AssociatedTranslation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if b == true {

	}
}
func (lc *LangueController) HandlePostChooseLanguage(w http.ResponseWriter, r *http.Request) {
	l := chi.URLParam(r, "language")
	err := lc.Lm.UserLangueSelection(r.Context(), l)
	if err != nil {
		log.Printf("%s", err)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
