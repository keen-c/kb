package models

import "database/sql"

type AdminCtxKey string

type UserCtxKey string
type Theme struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Answer struct {
	WordID                  string `json:"w_id,omitempty"`
	AssociatedTranslationID string `json:"atid,omitempty"`
	RandomTranslationID     string `json:"rid,omitempty"`
	Answer                  string `json:"answer"`
}
type Question struct {
	WordID                  string `json:"w_id,omitempty"`
	Word                    string `json:"word,omitempty"`
	AssociatedTranslationID string `json:"atid,omitempty"`
	AssociatedTranslation   string `json:"translation,omitempty"`
	RandomTranslationID     string `json:"rid,omitempty"`
	RandomTranslation       string `json:"random,omitempty"`
	CurrentTheme            string `json:"theme,omitempty"`
}
type Word struct {
	ID          string `json:"id"`
	Word        string `json:"word"`
	LangId      string `json:"lang_id"`
	ThemeId     string `json:"theme_id"`
	Translation string `json:"translation"`
}
type Language struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
type LanguageManager struct {
	DB *sql.DB
}
