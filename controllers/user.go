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
		http.Error(w, ServerError, http.StatusInternalServerError)
		return
	}
	if b {
		m := ErrorMessage{
			Message: "Cet email a dèjà utilisé",
		}
		err := json.NewEncoder(w).Encode(m)
		if err != nil {
			http.Error(w, ServerError, http.StatusInternalServerError)
			return
		}
		return
	}
	if u.Password != u.RepeatPassword {
		http.Error(w, BR, http.StatusBadRequest)
		return
	}
	if _, err = mail.ParseAddress(u.Email); err != nil {
		http.Error(w, BR, http.StatusBadRequest)
		return
	}
	user, err := uc.Um.Create(r.Context(), u.Pseudo, u.Email, u.Password)
	if err != nil {
		log.Printf("Error %s", err)
		http.Error(w, BR, http.StatusBadRequest)
		return
	}
	s, err := uc.SM.CreateSession(r.Context(), user.ID)
	if err != nil {
		log.Printf("error on creating session %s", err)
		http.Error(w, ServerError, http.StatusInternalServerError)
		return
	}
	SetCookie(w, session, s.Token)
	w.WriteHeader(http.StatusCreated)
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
		http.Error(w, ServerError, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusFound)
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		log.Printf("%s", err)
		http.Error(w, ServerError, http.StatusInternalServerError)
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
	http.Redirect(w, r, "http://localhost:5173/", http.StatusMovedPermanently)
}
