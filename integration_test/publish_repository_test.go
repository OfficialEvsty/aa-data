package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	integration "github.com/OfficialEvsty/aa-data/integration_test"
	"github.com/OfficialEvsty/aa-data/repos"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"os"
	"testing"
)

var testDB *sql.DB
var pgContainer testcontainers.Container

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../migrations", // путь к миграциям
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up failed: %w", err)
	}

	return nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, dsn, err := integration.PostgresSetup(ctx)
	if err != nil {
		log.Fatal(err)
	}
	pgContainer = container
	defer pgContainer.Terminate(ctx)

	testDB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := testDB.PingContext(ctx); err != nil {
		log.Fatalf("cannot ping test database: %v", err)
	}
	err = runMigrations(testDB)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("migrations success")
	code := m.Run()
	os.Exit(code)
}
func TestInsertPublish(t *testing.T) {
	ctx := context.Background()
	log.Println("starts insert publishing transaction...")
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	repo := repos.NewPublishRepository(testDB)
	id := uuid.New()

	key := "somekey.png"
	bucket := "somecoolbucket"
	s3name := "selectel"
	publish := domain.PublishedScreenshot{
		ID: id,
		S3Data: serializable.S3Screenshot{
			Key:    key,
			Bucket: bucket,
			S3Name: s3name,
		},
	}
	err = repo.WithTx(tx).Add(ctx, publish)
	if err != nil {
		t.Fatal(err)
	}
	pub, err := repo.WithTx(tx).Get(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	isPubDataCorrect := pub.S3Data.Bucket == bucket &&
		pub.S3Data.S3Name == s3name &&
		pub.S3Data.Key == key
	if !isPubDataCorrect {
		t.Fatal("expected published screen shot to be present")
	}
}
