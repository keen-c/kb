package controllers

import (
	"fmt"
	"html/template"
	"kb/kb/models"
	"kb/kb/views"
	"log"
	"net/http"
	"strings"
)

type AdminControllers struct {
	Am *models.AdminManager
	Sm *models.SessionManager
	Lm *models.LanguageManager
}

func (ac *AdminControllers) HandleGetConnexion(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFS(views.Static, "connexion.html"))
	if err := t.Execute(w, nil); err != nil {
		log.Printf("%s", err)
		http.Error(w, ServerError, http.StatusInternalServerError)
	}
}

func (ac *AdminControllers) HandlePostConnexion(w http.ResponseWriter, r *http.Request) {
	e := r.FormValue("email")
	p := r.FormValue("password")
	ad, err := ac.Am.Connexion(r.Context(), e, p)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	s, err := ac.Sm.CreateSessionAdmin(r.Context(), ad.ID)
	if err != nil {
		log.Printf("error creating session %s", err)
		return
	}
	SetCookie(w, session, s.Token)

	http.Redirect(w, r, "/admin/panel", http.StatusSeeOther)
}

func (ac *AdminControllers) HandleGetPanel(w http.ResponseWriter, r *http.Request) {
	t := template.Must(
		template.ParseFS(
			views.Static,
			"panel.html",
			"components/components.html",
		),
	)
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, ServerError, http.StatusInternalServerError)
	}
}

func (ac *AdminControllers) HandlePostLanguage(w http.ResponseWriter, r *http.Request) {
	l := strings.ToLower(r.FormValue("langue"))
	d := strings.ToLower(r.FormValue("disponible"))
	langue, err := ac.Am.InsertLangue(r.Context(), l, d)
	if err != nil {
		log.Printf("%s", err)
		http.Error(w, "Not Ok", http.StatusInternalServerError)
		return
	}
	log.Printf("%v", langue)
	t := template.Must(template.ParseFS(views.Static, "components/components.html"))
	if err := t.ExecuteTemplate(w, "li-langue", langue); err != nil {
		log.Printf("parsing template : %s", err)
		return
	}
}
func (ac *AdminControllers) HandleGetLanguages(w http.ResponseWriter, r *http.Request) {
	l, err := ac.Lm.QueryAvailableLanguage(r.Context())
	if err != nil {
		log.Printf("error querying available language %s", err)
		return
	}
	t := template.Must(template.ParseFS(views.Static, "components/components.html"))
	if err := t.ExecuteTemplate(w, "options-languages", l); err != nil {
		log.Printf("error executing %s", err)
		return
	}
}

func (ac *AdminControllers) HandlePostTheme(w http.ResponseWriter, r *http.Request) {
	t := strings.ToLower(r.FormValue("theme"))
	language_id := r.FormValue("languages")
	err := ac.Lm.CreateThemeByLangue(r.Context(), language_id, t)
	if err != nil {
		log.Printf("%s", err)
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}
}

func (ac *AdminControllers) HandleSearchThemeByLangue(w http.ResponseWriter, r *http.Request) {
	lan := r.FormValue("language")
	themes, err := ac.Lm.QueryThemeByLangue(r.Context(), lan)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	t := template.Must(template.ParseFS(views.Static, "components/components.html"))
	if err := t.ExecuteTemplate(w, "options-languages", themes); err != nil {
		log.Printf("error executing %s", err)
		return
	}
}

func (ac *AdminControllers) HandlePostWordByTheme(w http.ResponseWriter, r *http.Request) {
	lan := r.FormValue("language")
	th := r.FormValue("theme")
	wd := strings.TrimSpace(strings.ToLower(r.FormValue("word")))
	tra := strings.TrimSpace(strings.ToLower(r.FormValue("translation")))
	fmt.Println(lan, th, wd, tra)
	err := ac.Lm.InsertWordAndTraduction(r.Context(), lan, th, wd, tra)
	if err != nil {
		log.Printf("%s", err)
		return
	}
}
