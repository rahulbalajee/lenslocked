package main

import (
	"context"
	"flag"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pressly/goose/v3"
	"github.com/rahulbalajee/lenslocked/cmd/goose/migrations"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

const postgresConfig = "host=localhost port=5432 user=baloo password=junglebook dbname=lenslocked sslmode=disable"

// ➜  goose git:(master) ✗ ./custom up
func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("goose: failed to parse flags: %v", err)
	}
	args := flags.Args()

	command := args[0]

	db, err := goose.OpenDBWithDriver("pgx", postgresConfig)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v", err)
		}
	}()

	// Set the embedded filesystem for migrations
	goose.SetBaseFS(migrations.FS)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	ctx := context.Background()
	if err := goose.RunContext(ctx, command, db, *dir); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
