package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Host     string
	DBname   string
	UserName string
	Password string
	Port     string
	SSLMode  string
}

func (cfg Config) ReturnDSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.DBname, cfg.UserName, cfg.Password, cfg.SSLMode)
}

func main() {
	cfg := Config{
		Host:     "localhost",
		DBname:   "lenslocked",
		UserName: "baloo",
		Password: "junglebook",
		Port:     "5432",
		SSLMode:  "disable",
	}

	db, err := sql.Open("pgx", cfg.ReturnDSN())
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to DB!")

	query := `CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT UNIQUE NOT NULL,
		contact INT
	);
	
	CREATE TABLE IF NOT EXISTS tweets(
		tweet_id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		tweet TEXT NOT NULL,
		parent_tweet_id INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (parent_tweet_id) REFERENCES tweets(tweet_id)
	);
	
	CREATE TABLE IF NOT EXISTS likes(
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		tweet_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (tweet_id) REFERENCES tweets(tweet_id),
		UNIQUE(user_id, tweet_id)
	);`

	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}
