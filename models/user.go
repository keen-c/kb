package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserCtxKey string

const (
	Userkey UserCtxKey = "user"
)

func SetUserContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, Userkey, user)
}

func GetUserByContext(ctx context.Context) (*User, error) {
	v := ctx.Value(Userkey)
	u, ok := v.(*User)
	if !ok {
		return nil, errors.New("casting user went wrong")
	}
	return u, nil
}

type UserCreate struct {
	Pseudo         string `json:"pseudo"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}
type User struct {
	ID     string `json:"id"`
	Pseudo string `json:"pseudo"`
	Email  string `json:"email"`
}
type UserManager struct {
	DB *sql.DB
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error generating password hash: %w", err)
	}
	return string(hashedPassword), nil
}

func CompareHashedPassword(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (um *UserManager) Create(
	ctx context.Context,
	username, email, password string,
) (*User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password during user creation: %w", err)
	}

	newUser := UserCreate{
		Pseudo:   username,
		Email:    email,
		Password: hashedPassword,
	}

	var user User
	err = um.DB.QueryRowContext(
		ctx,
		`insert into users (username, email, password_hashed) values ($1, $2, $3) returning id, username, email;`,
		newUser.Pseudo,
		newUser.Email,
		newUser.Password,
	).Scan(&user.ID, &user.Pseudo, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("error inserting new user: %w", err)
	}

	return &user, nil
}

func (um *UserManager) FindUserByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `select exists(select email from users where email = $1)`
	err := um.DB.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w", err)
	}
	return exists, nil
}