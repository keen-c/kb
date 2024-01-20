package main

import (
	"database/sql"
	"kb/kb/controllers"
	"kb/kb/models"
	"kb/kb/views"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/markbates/goth/gothic"
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
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
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
	lc := controllers.LangueController{
		Lm: &models.LanguageManager{DB: db},
	}
	uc := controllers.UserControllers{
		Um: &models.UserManager{DB: db},
		SM: &models.SessionManager{DB: db},
		Lm: &models.LanguageManager{DB: db},
	}
	ad := controllers.AdminControllers{
		Am: &models.AdminManager{DB: db},
		Sm: &models.SessionManager{DB: db},
		Lm: &models.LanguageManager{DB: db},
	}

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(views.Static))))
	r.Route("/auth", func(r chi.Router) {
		r.Get("/", gothic.BeginAuthHandler)
		r.Get("/callback", uc.CreateWithGoogleCallback)
	})
	r.Route("/user", func(r chi.Router) {
		r.Post("/create", uc.Create)
		r.Post("/connexion", uc.Connexion)
		r.Get("/userbytoken/{token}", uc.UserByToken)
		r.Get("/languages", uc.HandleGetAvailableLanguages)
		r.With(mw.GetCurrentUser).Post("/choose/{language}", uc.HandlePostChooseLanguage)
	})
	r.Route("/language", func(r chi.Router) {
		r.With(mw.GetCurrentUser).Get("/", lc.HandleGetThemeByLanguage)
	})

	r.Route("/admin", func(r chi.Router) {
		r.Get("/connexion", ad.HandleGetConnexion)
		r.Post("/connexion", ad.HandlePostConnexion)
		r.Route("/panel", func(r chi.Router) {
			r.With(mw.GetCurrentAdmin).Get("/", ad.HandleGetPanel)
			r.Get("/languages", ad.HandleGetLanguages)
			r.Post("/languages", ad.HandlePostLanguage)
			r.Post("/theme", ad.HandlePostTheme)
			r.Get("/languages-themes", ad.HandleSearchThemeByLangue)
			r.Post("/word-theme", ad.HandlePostWordByTheme)
		})
	})
	log.Println("http://localhost:8080/admin/panel")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Panic(err)
	}
}
