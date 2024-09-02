package postgres_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/pkg/models"
	"github.com/oalexander6/passman/pkg/store/postgres"
	"github.com/testcontainers/testcontainers-go"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var pgOpts = config.PostgresConfig{
	URI: "",
}

func mustStartPostgresContainer() (func(context.Context) error, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := pg.Run(
		context.Background(),
		"postgres:latest",
		pg.WithDatabase(dbName),
		pg.WithUsername(dbUser),
		pg.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	pgOpts.URI = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", dbUser, dbPwd, dbHost, dbPort.Port(), dbName, "public")

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
}

func TestNew(t *testing.T) {
	srv := postgres.New(pgOpts)
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestCreateNote(t *testing.T) {
	srv := postgres.New(pgOpts)

	result, err := srv.NoteCreate(context.Background(), models.NoteCreateParams{Name: "Test Note 1", Value: "testval1"})
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if result.ID == 0 {
		t.Fatal("Expected non-zero id")
	}

	if result.CreatedAt != result.UpdatedAt {
		t.Fatal("Expected created_at and updated_at to match")
	}

	if result.Name != "Test Note 1" || result.Value != "testval1" {
		t.Fatal("Name or value did not match input")
	}
}