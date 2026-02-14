package adapter

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	_ "github.com/lib/pq"
)

type storage struct {
	db *sql.DB
}

type Storage interface {
	Save(*domain.Link) error
	FindByShortCode(string) (*domain.Link, error)
	FindByURL(string) (*domain.Link, error)
}

func ConnToStorage() (*storage, error) {
	connStr := makeConnStr()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(1)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(2)
		return nil, err
	}

	query := `
		CREATE TABLE IF NOT EXISTS links (
			id SERIAL PRIMARY KEY,
			original_url TEXT NOT NULL,
			short_url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
	);`

	_, err = db.Exec(query)
	if err != nil {
		fmt.Println(3)
		return nil, err
	}

	return &storage{
		db: db,
	}, nil
}

func makeConnStr() string {
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
}
