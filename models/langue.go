package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Theme struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

func (lm *LanguageManager) QueryAvailableLanguage(ctx context.Context) (*[]Language, error) {
	query := `select id, name from languages where languages.status = 'available'`
	rows, err := lm.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("can not query language")
	}
	defer func(r *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("can not close language")
		}
	}(rows)
	var languages []Language
	for rows.Next() {
		var lang Language
		err = rows.Scan(&lang.ID, &lang.Name)
		if err != nil {
			return nil, fmt.Errorf("can not scan language")
		}
		languages = append(languages, lang)
	}
	return &languages, nil
}

func (lm *LanguageManager) UserLangueSelection(ctx context.Context, langueId string) error {
	u, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	query := `INSERT INTO user_languages (user_id, language_id) VALUES ($1, $2)`
	_, err = lm.DB.ExecContext(ctx, query, u.ID, langueId)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (lm *LanguageManager) QueryTheme(ctx context.Context, langueId string) (*Theme, error) {
	query := `SELECT id, name FROM themes WHERE language_id = $1 AND published = TRUE ORDER BY created_at LIMIT 1;
`
	var t Theme
	err := lm.DB.QueryRowContext(ctx, query, langueId).Scan(&t.ID, &t.Name)
	if err != nil {
		return nil, fmt.Errorf("can not fetch theme")
	}
	return &t, nil
}

func (lm *LanguageManager) QueryUserLangue(ctx context.Context) (string, error) {
	u, err := GetUserByContext(ctx)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}
	query := `select name from languages left join user_languages on user_languages.language_id = languages.id and user_languages.user_id =$1 limit 1;`
	var lang string
	err = lm.DB.QueryRowContext(ctx, query, u.ID).Scan(&lang)
	if err != nil {
		return "", fmt.Errorf("can not fetch user langue")
	}
	return lang, nil
}

func (lm *LanguageManager) QueryIDLanguage(ctx context.Context, languageID string) (string, error) {
	var id string
	query := `select id from languages where name = $1;`
	err := lm.DB.QueryRowContext(ctx, query, languageID).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("can not fetch id language")
	}
	return id, nil
}

func (lm *LanguageManager) GetWordByTheme(ctx context.Context, themeID string) (*Word, error) {
	query := `select w.id, w.word, w.language_id, w.theme_id, t.translation from words w left join translations t on t.word_id = w.id where w.theme_id = $1 limit 1;;`
	var w Word
	err := lm.DB.QueryRowContext(ctx, query, themeID).
		Scan(&w.ID, &w.Word, &w.LangId, &w.ThemeId, &w.Translation)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &w, nil
}

func (lm *LanguageManager) GetRandomWord(
	ctx context.Context,
	wordID, themeID string,
) (*Word, error) {
	query := `select id, word from words where id != $1 and theme_id = $2 order by random() limit 1;`
	var w Word
	err := lm.DB.QueryRowContext(ctx, query, wordID, themeID).Scan(&w.ID, &w.Word)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &w, nil
}

func (lm *LanguageManager) GetTheResponse(
	ctx context.Context,
	wordID, translation string,
) (bool, error) {
	var b bool
	query := `SELECT EXISTS (SELECT 1 FROM translations WHERE word_id = $1 AND translation = $2);`
	err := lm.DB.QueryRowContext(ctx, query, wordID, translation).Scan(&b)
	if err != nil {
		return false, fmt.Errorf("%w", err)
	}
	return b, nil
}

func (Lm *LanguageManager) QueryThemeByLangue(
	ctx context.Context,
	language_id string,
) (*[]Theme, error) {
	var t []Theme
	query := "select t.id, t.name from themes t where t.language_id = $1"
	rows, err := Lm.DB.QueryContext(ctx, query, language_id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer func(s *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("%s", err)
		}
	}(rows)
	for rows.Next() {
		var th Theme
		err := rows.Scan(&th.ID, &th.Name)
		if err != nil {
			log.Printf("%s", err)
			break
		}
		t = append(t, th)
	}
	return &t, nil
}

func (lm *LanguageManager) CreateThemeByLangue(
	ctx context.Context,
	langue_id, theme string,
) error {
	query := `insert into themes(name,language_id) values ($1,$2)`
	_, err := lm.DB.ExecContext(ctx, query, theme, langue_id)
	if err != nil {
		return err
	}
	return nil
}

func (lm *LanguageManager) CreateWordByThemeByLangue(
	ctx context.Context,
	langue_id, theme_id, word string,
) error {
	query := `insert into words(language_id, theme_id,word) values($1, $2, $3)`
	_, err := lm.DB.ExecContext(ctx, query, langue_id, theme_id, word)
	if err != nil {
		return err
	}
	return nil
}
