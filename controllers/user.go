package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"kb/kb/models"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

const (
	ServerError = "Internal Server Error"
	BR          = "Bad Request"
)

type UserControllers struct {
	Um *models.UserManager
	SM *models.SessionManager
	Lm *models.LanguageManager
}
type ErrorMessage struct {
	Message string `json:"message,omitempty"`
}

func (uc *UserControllers) Create(w http.ResponseWriter, r *http.Request) {
	var u models.UserCreate
	err := decodeJSONBody(w, r, &u)
	if err != nil {
		log.Printf("creating user %s", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	b, err := uc.Um.FindUserByEmail(r.Context(), u.Email)
	if err != nil {
		log.Printf("error searching user %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if b {
		m := ErrorMessage{
			Message: "Cet email a dèjà utilisé",
		}
		JsonError(w, m, http.StatusBadRequest)
	}
	if u.Password != u.RepeatPassword {
		m := ErrorMessage{
			Message: "Les mot de passe ne correspondent pas.",
		}
		JsonError(w, m, http.StatusBadRequest)
		return
	}
	if _, err = mail.ParseAddress(u.Email); err != nil {
		m := ErrorMessage{
			Message: "Address email invalide !",
		}
		JsonError(w, m, http.StatusBadRequest)
		return
	}
	user, err := uc.Um.Create(r.Context(), u.Pseudo, u.Email, u.Password)
	if err != nil {
		log.Printf("Error %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s, err := uc.SM.CreateSession(r.Context(), user.ID)
	if err != nil {
		log.Printf("error on creating session %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	SetCookie(w, session, s.Token)
}

func (uc *UserControllers) UserByToken(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "token")
	u, err := uc.SM.FindUserByCookie(r.Context(), t)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("error on getting user : %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (uc *UserControllers) CreateWithGoogleCallback(w http.ResponseWriter, r *http.Request) {
	u, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
	}
	b, err := uc.Um.FindUserByEmail(r.Context(), u.Email)
	if err != nil {
		log.Printf("error getting user : %s", err)
		http.Error(w, ServerError, http.StatusInternalServerError)
		return
	}
	if b {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Cet email est dèja utilisé")
		return
	}
	pseudo, err := ExtractUsernameFromEmail(u.Email)
	if err != nil {
		log.Printf("error on extract %s", err)
	}
	password := strings.Join([]string{u.UserID, pseudo}, "")
	user, err := uc.Um.Create(r.Context(), pseudo, u.Email, password)
	if err != nil {
		log.Printf("error on creating user provider : %s", err)
		http.Error(w, ServerError, http.StatusInternalServerError)
		return
	}
	s, err := uc.SM.CreateSession(r.Context(), user.ID)
	if err != nil {
		http.Error(w, ServerError, http.StatusInternalServerError)
		return
	}
	SetCookie(w, session, s.Token)
	http.Redirect(w, r, "http://localhost:5173/choisir-une-langue", http.StatusSeeOther)
}
func (uc *UserControllers) HandleGetAvailableLanguages(w http.ResponseWriter, r *http.Request) {
	ls, err := uc.Lm.QueryAvailableLanguage(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err)
		return
	}
	err = json.NewEncoder(w).Encode(ls)
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func (uc *UserControllers) Connexion(w http.ResponseWriter, r *http.Request) {
	var u models.UserCreate
	err := decodeJSONBody(w, r, &u)
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := uc.Um.Connexion(r.Context(), u.Email, u.Password)
	if err != nil {
		m := ErrorMessage{
			Message: "Email ou Mot de passe incorrect.",
		}
		JsonError(w, m, http.StatusBadRequest)
		return
	}
	s, err := uc.SM.CreateSession(r.Context(), id)
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	SetCookie(w, session, s.Token)
}
func (uc *UserControllers) HandlePostChooseLanguage(w http.ResponseWriter, r *http.Request) {
	l := chi.URLParam(r, "language")
	err := uc.Lm.UserLangueSelection(r.Context(), l)
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
func (uc *UserControllers) HandleGetGame(w http.ResponseWriter, r *http.Request) {
	id, err := uc.Lm.GetCurrentTheme(r.Context())
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("%s", id)

}
