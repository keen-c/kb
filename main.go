package main

import (
	"database/sql"
	"fmt"
	"kb/kb/controllers"
	"kb/kb/models"
	"kb/kb/views"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/gorilla/csrf"
	_ "github.com/lib/pq"
	"github.com/markbates/goth/gothic"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panic("cant load env variables :", err)
	}
	db, err := sql.Open(
		"postgres",
		os.Getenv("DATABASE"),
	)
	if err != nil {
		log.Printf("error on db :%s", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("close db %s", err)
		}
	}(db)
	err = db.Ping()
	if err != nil {
		log.Printf("ping %s", err)
	}
	controllers.Oauth()
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},

		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Cookie"},
		ExposedHeaders:   []string{},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mw := controllers.Middleware{
		DB: db,
		Sm: &models.SessionManager{DB: db},
	}
	uc := controllers.UserControllers{
		Um: &models.UserManager{DB: db},
		SM: &models.SessionManager{DB: db},
	}
	ad := controllers.AdminControllers{
		Am: &models.AdminManager{DB: db},
		Sm: &models.SessionManager{DB: db},
		Lm: &models.LanguageManager{DB: db},
	}
	CSRF := csrf.Protect([]byte(os.Getenv("CSRF")))
	csrf.Secure(false)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(views.Static))))

	r.Route("/auth", func(r chi.Router) {
		r.Get("/", gothic.BeginAuthHandler)
		r.Get("/callback", uc.CreateWithGoogleCallback)
	})
	r.Route("/user", func(r chi.Router) {
		r.Post("/create", uc.Create)
		r.Get("/userbytoken/{token}", uc.UserByToken)
	})
	r.Route("/admin", func(r chi.Router) {
		r.Get("/connexion", ad.HandleGetConnexion)
		r.Post("/connexion", ad.HandlePostConnexion)
		r.Route("/panel", func(r chi.Router) {
			r.With(mw.GetCurrentAdmin).Get("/", ad.HandleGetPanel)
			r.Get("/languages", ad.HandleGetLanguage)
			r.Post("/languages", ad.HandlePostLanguage)
			r.Post("/theme", ad.HandlePostTheme)
			r.Get("/languages-themes", ad.HandleSearchThemeByLangue)
			r.Post("/word-theme", ad.HandlePostWordByTheme)
		})
	})
	fmt.Println("http://localhost:8080/admin/panel")
	err = http.ListenAndServe(":8080", CSRF(r))
	if err != nil {
		log.Panic(err)
	}
}
