package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"kb/kb/models"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type LangueController struct {
	Lm *models.LanguageManager
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
func (lc *LangueController) HandleGetGame(w http.ResponseWriter, r *http.Request) {
	q, err := lc.Lm.GetUserCurrentQueue(r.Context())
	if err != nil {
		switch err.Error() {
		case models.UnexpectedErrorJson:
			five, err := lc.Lm.GetFiveWordFromTheCurrentTheme(r.Context())
			if err != nil {
				log.Printf("%s", err)
				HttpInternalError(w, err)
				return
			}
			b, err := five.Marshall()
			if err != nil {
				log.Printf("%s", err)
				HttpInternalError(w, err)
				return
			}
			if err := lc.Lm.UpdateQueue(r.Context(), b); err != nil {
				log.Printf("%s", err)
				HttpInternalError(w, err)
				return
			}
			g := Game{
				Game: *five.First(),
			}
			JsonSendGame(w, g)
			return
		default:
			HttpInternalError(w, err)
			return
		}
	}
	g := Game{
		Game: *q.First(),
	}
	JsonSendGame(w, g)
}

func (lc *LangueController) HandlePostGame(w http.ResponseWriter, r *http.Request) {
	answer_id := chi.URLParam(r, "answer")
	q, err := lc.Lm.GetUserCurrentQueue(r.Context())
	if err != nil {
		log.Printf("%s", err)
		HttpInternalError(w, err)
		return
	}
	b, err := lc.Lm.CheckTheAnswer(r.Context(), answer_id, q.First().AssociatedTranslation)
	if err != nil {
		log.Printf("%s", err)
		HttpInternalError(w, err)
		return
	}
	switch b {
	case true:
		q.Dequeue()
		if q.Len() == 0 {
			next, err := lc.Lm.GetFiveWordFromTheCurrentTheme(r.Context())
			if err != nil {
				fmt.Println("there is a error")
				return
			}
			// if errors.Is(err, sql.ErrNoRows) {
			// 	t, err := lc.Lm.GetCurrentTheme(r.Context())
			// 	if err != nil {
			// 		HttpInternalError(w, err)
			// 		return
			// 	}
			// 	if err = lc.Lm.InsertThemeDone(r.Context(), t.ID); err != nil {
			// 		log.Printf("%s", err)
			// 		HttpInternalError(w, err)
			// 		return
			// 	}
			// 	lc.HandleGetGame(w, r)
			// 	return
			// }
			m, err := next.Marshall()
			if err != nil {
				log.Printf("%s", err)
				HttpInternalError(w, err)
				return
			}
			if err := lc.Lm.UpdateQueue(r.Context(), m); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			g := Game{
				Answer: true,
				Game: *next.First(),
			}
			JsonSendGame(w, g)
			return
		}
		m, err := q.Marshall()
		if err != nil {
			log.Printf("%s", err)
			HttpInternalError(w, err)
			return
		}
		if err := lc.Lm.UpdateQueue(r.Context(), m); err != nil {
			log.Printf("%s", err)
			HttpInternalError(w, err)
			return
		}
		g := Game{
			Answer: true,
			Game:   *q.First(),
		}
		JsonSendGame(w, g)
		return
	case false:
		q.SwapEnd()
		m, err := q.Marshall()
		if err != nil {
			HttpInternalError(w, err)
			return
		}

		if err := lc.Lm.UpdateQueue(r.Context(), m); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		g := Game{
			Answer: false,
		}

		JsonSendGame(w, g)
		return
	}
}
