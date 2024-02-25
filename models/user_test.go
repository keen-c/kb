package models

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword"

	// Test pour un succès de hashage
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword() error = %v, wantErr %v", err, false)
	}
	if len(hash) == 0 {
		t.Errorf("HashPassword() got = %v, want non-empty string", hash)
	}

	// Vous pouvez ajouter d'autres cas de test ici, si nécessaire
}
func TestCompareHashedPassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, _ := HashPassword(password)

	// Test de correspondance
	err := CompareHashedPassword(password, hashedPassword)
	if err != nil {
		t.Errorf("CompareHashedPassword() should have successfully compared the password, got error %v", err)
	}

	// Test d'échec
	err = CompareHashedPassword("wrongpassword", hashedPassword)
	if err == nil {
		t.Errorf("CompareHashedPassword() should have failed to compare the wrong password")
	}
}
func TestUserManager_FindUserByEmail(t *testing.T) {
	// Configuration du mock de la base de données ici
	// ...
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
	um := UserManager{DB: db}
	exists, err := um.FindUserByEmail(context.Background(), "email@test.com")
	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}
	if !exists {
		t.Errorf("FindUserByEmail() got = %v, want %v", exists, true)
	}

	// Test pour un email qui n'existe pas
	// ...
}
func TestUserManager_Connexion(t *testing.T) {
	// Configuration du mock de la base de données ici
	// ...

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
	um := UserManager{DB: db}
	id, err := um.Connexion(context.Background(), "email@test.com", "password")
	if err != nil {
		t.Fatalf("Connexion() error = %v", err)
	}
	if id == "" {
		t.Errorf("Connexion() got empty id, want non-empty string")
	}

	// Test avec un mot de passe incorrect
	// ...
}
