package controllers

import (
	"database/sql"
	"kb/kb/models"
	"log"
	"net/http"
)

type Middleware struct {
	DB *sql.DB
	Sm *models.SessionManager
}

func (m *Middleware) GetCurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		c, err := ReadCookie(r, session)
		if err != nil {
			log.Println("Cookie error :", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		u, err := m.Sm.FindUserByCookie(ctx, c.Value)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx = models.SetUserContext(ctx, u)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) GetCurrentAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		c, err := ReadCookie(r, session)
		if err != nil {
			log.Printf("Cookie error %s :", err)
			http.Redirect(w, r, "/admin/connexion", http.StatusSeeOther)
			return
		}
		a, err := m.Sm.FindAdminByCookie(ctx, c.Value)
		if err != nil {
			log.Printf("find admin by cookie %s", err)
			http.Redirect(w, r, "/admin/connexion", http.StatusSeeOther)
			return
		}
		ctx = models.SetAdminContext(ctx, a)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// func (m *Middleware)AuthorizedLang(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 	})
// }
