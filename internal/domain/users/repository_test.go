//go:build integration
// +build integration

package users

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var testDB *sqlx.DB

// 🔥 Run once before all tests
func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		panic("TEST_DATABASE_URL must be set")
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	// 🔥 Resolve absolute path ke folder migrations
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// dari internal/domain/categories naik 3 level ke root
	rootPath := filepath.Join(wd, "../../..")
	migrationsPath := "file://" + filepath.Join(rootPath, "migrations")

	migrator, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}

	// 🔥 Drop & recreate schema ONCE
	// err = migrator.Drop()
	// if err != nil && err != migrate.ErrNoChange {
	// 	panic(err)
	// }

	// err = migrator.Up()
	// if err != nil && err != migrate.ErrNoChange {
	// 	panic(err)
	// }

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	testDB = db

	code := m.Run()

	db.Close()
	os.Exit(code)
}

func setupTestDB(t *testing.T) *sqlx.DB {
	_, err := testDB.Exec(`
		TRUNCATE TABLE 
			users,
			categories,
			savings,
			category_budgets,
			transactions,
			salaries
		RESTART IDENTITY CASCADE
	`)
	require.NoError(t, err)

	return testDB
}

func TestUserRepository_Create_And_FindByEmail(t *testing.T) {
	db := setupTestDB(t)

	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &User{
		USER_ID:  "USER-20260216-000001",
		Fullname: "Bos Backend",
		Password: "secret",
		Email:    "bos@test.com",
		Username: "bos",
	}

	// Create
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// FindByEmail
	found, err := repo.FindByEmail(ctx, "bos@test.com")
	require.NoError(t, err)
	require.Equal(t, "bos", found.Username)
	require.Equal(t, "bos@test.com", found.Email)
}

func TestUserRepository_Find_NotFound(t *testing.T) {
	db := setupTestDB(t)

	repo := NewUserRepository(db)
	ctx := context.Background()

	_, err := repo.FindByEmail(ctx, "tidakada@test.com")
	require.Error(t, err)
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db := setupTestDB(t)

	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &User{
		USER_ID:  "USER-20260216-000002",
		Fullname: "Bos Username",
		Password: "secret",
		Email:    "username@test.com",
		Username: "bosuser",
	}

	require.NoError(t, repo.Create(ctx, user))

	found, err := repo.FindByUsername(ctx, "bosuser")
	require.NoError(t, err)
	require.Equal(t, "username@test.com", found.Email)
}

func TestUserRepository_FindByUserID(t *testing.T) {
	db := setupTestDB(t)

	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &User{
		USER_ID:  "USER-20260216-000003",
		Fullname: "Bos UserID",
		Password: "secret",
		Email:    "userid@test.com",
		Username: "bosid",
	}

	require.NoError(t, repo.Create(ctx, user))

	found, err := repo.FindByUserID(ctx, "USER-20260216-000003")
	require.NoError(t, err)
	require.Equal(t, "bosid", found.Username)
}

func TestUserRepository_CountByDate(t *testing.T) {
	db := setupTestDB(t)

	repo := NewUserRepository(db)
	ctx := context.Background()

	date := time.Now().Format("20060102")

	user1 := &User{
		USER_ID:  "USER-" + date + "-000001",
		Fullname: "User 1",
		Password: "secret",
		Email:    "count1@test.com",
		Username: "count1",
	}

	user2 := &User{
		USER_ID:  "USER-" + date + "-000002",
		Fullname: "User 2",
		Password: "secret",
		Email:    "count2@test.com",
		Username: "count2",
	}

	require.NoError(t, repo.Create(ctx, user1))
	require.NoError(t, repo.Create(ctx, user2))

	count, err := repo.CountByDate(ctx, date)
	require.NoError(t, err)
	require.Equal(t, 2, count)
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)

	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &User{
		USER_ID:  "USER-20260216-000010",
		Fullname: "Duplicate",
		Password: "secret",
		Email:    "duplicate@test.com",
		Username: "dup1",
	}

	require.NoError(t, repo.Create(ctx, user))

	duplicate := &User{
		USER_ID:  "USER-20260216-000011",
		Fullname: "Duplicate 2",
		Password: "secret",
		Email:    "duplicate@test.com", // sama
		Username: "dup2",
	}

	err := repo.Create(ctx, duplicate)
	require.Error(t, err)
}
