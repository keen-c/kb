package controllers

import (
	"database/sql"
	"errors"
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
	correct, err := lc.Lm.CheckTheAnswer(r.Context(), answer_id, q.First().AssociatedTranslation)
	if err != nil {
		log.Printf("%s", err)
		HttpInternalError(w, err)
		return
	}
	if correct {
		lc.CorrectProcessing(w, r, q, correct)
	} else {
		lc.IncorrectProcess(w, r, q, correct)
	}
}
func (lc *LangueController) ChangeTheme(w http.ResponseWriter, r *http.Request) {
	t, err := lc.Lm.GetCurrentTheme(r.Context())
	if err != nil {
		HttpInternalError(w, err)
		return
	}
	if err = lc.Lm.InsertThemeDone(r.Context(), t.ID); err != nil {
		log.Printf("%s", err)
		HttpInternalError(w, err)
		return
	}
	lc.HandleGetGame(w, r)
}

func (lc *LangueController) CorrectProcessing(w http.ResponseWriter, r *http.Request, q *models.Queux, answer bool) {
	if err := lc.Lm.InsertWordViews(r.Context(), q.First().WordID); err != nil {
		log.Printf("%s", err)
		HttpInternalError(w, err)
		return
	}
	q.Dequeue()
	if q.Len() == 0 {
		lc.EmptyQueue(w, r)
		return
	}
	lc.UpdateQueueAndSendGame(w, r, q, answer)
}

func (lc *LangueController) EmptyQueue(w http.ResponseWriter, r *http.Request) {
	next, err := lc.Lm.GetFiveWordFromTheCurrentTheme(r.Context())
	if err != nil || *next == nil {
		if *next == nil {
			lc.ChangeTheme(w, r)
		} else {
			HttpInternalError(w, err)
		}
		return
	}
	lc.UpdateQueueAndSendGame(w, r, next, true)
}
func (lc *LangueController) UpdateQueueAndSendGame(w http.ResponseWriter, r *http.Request, q *models.Queux, answer bool) {
	marshalledQueue, err := q.Marshall()
	if err != nil {
		HttpInternalError(w, err)
		return
	}

	if err := lc.Lm.UpdateQueue(r.Context(), marshalledQueue); err != nil {
		HttpInternalError(w, err)
		return
	}

	g := Game{
		Answer: answer,
		Game:   *q.First(),
	}
	JsonSendGame(w, g)
}
func (lc *LangueController) IncorrectProcess(w http.ResponseWriter, r *http.Request, q *models.Queux, answer bool) {
	q.SwapEnd()
	lc.UpdateQueueAndSendGame(w, r, q, answer)
}
