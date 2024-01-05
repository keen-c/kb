package models

import (
	"context"
	"database/sql"
	"fmt"
)

type Session struct {
	ID      string
	UserID  string
	Token   string
	AdminID string
}
type SessionManager struct {
	DB *sql.DB
}

func (sm *SessionManager) CreateSession(ctx context.Context, userID string) (*Session, error) {
	var s Session
	err := sm.DB.QueryRowContext(ctx, `insert into sessions (user_id) values ($1) returning id, user_id, token`, userID).
		Scan(&s.ID, &s.UserID, &s.Token)
	if err != nil {
		return nil, fmt.Errorf("error when creating session %w", err)
	}
	return &s, nil
}

func (sm *SessionManager) CreateSessionAdmin(ctx context.Context, userID string) (*Session, error) {
	var s Session
	err := sm.DB.QueryRowContext(ctx, `insert into sessions (admin_id) values ($1) returning id, admin_id, token`, userID).
		Scan(&s.ID, &s.AdminID, &s.Token)
	if err != nil {
		return nil, fmt.Errorf("error when creating session %w", err)
	}
	return &s, nil
}

func (sm *SessionManager) FindUserByCookie(ctx context.Context, Token string) (*User, error) {
	var u User
	err := sm.DB.QueryRowContext(ctx, `select users.id, users.username,
	 users.email from users inner join sessions
	on users.id = sessions.user_id where 
	sessions.token = $1`, Token).Scan(&u.ID, &u.Pseudo, &u.Email)
	if err != nil {
		return nil, fmt.Errorf("fetching user by session error : %w", err)
	}
	return &u, nil
}

func (sm *SessionManager) FindAdminByCookie(ctx context.Context, Token string) (*Admin, error) {
	var a Admin
	query := `select admin.id, admin.email from admin inner join sessions on admin.id = sessions.admin_id where sessions.token = $1`
	err := sm.DB.QueryRowContext(ctx, query, Token).Scan(&a.ID, &a.Email)
	if err != nil {
		return nil, fmt.Errorf("fetching user by session error : %w", err)
	}
	return &a, nil
}
