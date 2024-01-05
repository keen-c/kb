package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type AdminCtxKey string

const (
	AdminKey AdminCtxKey = "admin"
)

func SetAdminContext(ctx context.Context, user *Admin) context.Context {
	return context.WithValue(ctx, AdminKey, user)
}

func GetAdminByContext(ctx context.Context) (*Admin, error) {
	v := ctx.Value(AdminKey)
	a, ok := v.(*Admin)
	if !ok {
		return nil, errors.New("casting admin went wrong")
	}
	return a, nil
}

type Admin struct {
	ID             string
	Email          string
	PasswordHashed string
}
type AdminManager struct {
	DB *sql.DB
}

func (am *AdminManager) CreateAdmin(email, password string) error {
	p, err := HashPassword(password)
	if err != nil {
		log.Printf("error creating admin %s", err)
		return err
	}
	_, err = am.DB.Exec(`INSERT INTO admin (email, passwordhashed) values ($1, $2) `, email, p)
	if err != nil {
		return err
	}
	return nil
}

func (am *AdminManager) Connexion(ctx context.Context, email, password string) (*Admin, error) {
	var admin Admin
	query := "SELECT id, email, passwordhashed FROM admin WHERE email = $1"
	err := am.DB.QueryRowContext(ctx, query, email).
		Scan(&admin.ID, &admin.Email, &admin.PasswordHashed)
	if err != nil {
		return nil, err
	}
	err = CompareHashedPassword(password, admin.PasswordHashed)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (am *AdminManager) FindAdmin(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `select exists(select email from admin where email = $1)`
	err := am.DB.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w", err)
	}
	return exists, nil
}

func (am *AdminManager) InsertLangue(
	ctx context.Context,
	langue, status string,
) (*Language, error) {
	var L Language
	query := `insert into languages (name, status) values ($1, $2) returning id, name, status`
	err := am.DB.QueryRowContext(ctx, query, langue, status).Scan(&L.ID, &L.Name, &L.Status)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &L, nil
}
