//go:build integration
// +build integration

package categories

import (
	"context"
	"database/sql"
	"errors"
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

func setupCategoriesTestDB(t *testing.T) *sqlx.DB {
	_, err := testDB.Exec(`
		TRUNCATE TABLE 
			categories,
			users
		RESTART IDENTITY CASCADE
	`)
	require.NoError(t, err)

	return testDB
}

func insertTestUser(t *testing.T, userID string) {
	_, err := testDB.Exec(`
		INSERT INTO users (user_id, username, fullname, email, password_hash)
		VALUES ($1, $2, $3, $4, $5)
	`,
		userID,
		"testuser",
		"Test User",
		"test@test.com",
		"secret",
	)
	require.NoError(t, err)
}

func TestCategoriesRepository_Create_And_CountByDate_Success(t *testing.T) {
	db := setupCategoriesTestDB(t)
	repo := NewCategoriesRepository(db)
	ctx := context.Background()

	userID := "USER-001"
	insertTestUser(t, userID)

	date := time.Now().Format("20060102")

	category1 := &Categories{
		CATEGORY_ID: "CAT-" + date + "-000001",
		USER_ID:     userID,
		Name:        "Food",
		Description: "Food expenses",
	}

	category2 := &Categories{
		CATEGORY_ID: "CAT-" + date + "-000002",
		USER_ID:     userID,
		Name:        "Transport",
		Description: "Transport expenses",
	}

	require.NoError(t, repo.Create(ctx, category1))
	require.NoError(t, repo.Create(ctx, category2))

	count, err := repo.CountByDate(ctx, date)
	require.NoError(t, err)
	require.Equal(t, 2, count)
}

func TestCategoriesRepository_Create_Duplicate(t *testing.T) {
	db := setupCategoriesTestDB(t)
	repo := NewCategoriesRepository(db)
	ctx := context.Background()

	userID := "USER-001"
	insertTestUser(t, userID)

	date := time.Now().Format("20060102")

	category := &Categories{
		CATEGORY_ID: "CAT-" + date + "-000001",
		USER_ID:     userID,
		Name:        "Food",
		Description: "Food expenses",
	}

	require.NoError(t, repo.Create(ctx, category))

	// insert lagi dengan ID yang sama
	err := repo.Create(ctx, category)

	require.Error(t, err)
}

func TestCreate_DBError(t *testing.T) {
	mockRepo := &mockCategoriesRepo{
		createFunc: func(ctx context.Context, category *Categories) error {
			return errors.New("db error")
		},
	}

	service := &categoriesService{
		categoryRepo: mockRepo,
	}

	err := service.Create(context.Background(), "Food", "desc")

	require.Error(t, err)
	require.Equal(t, ErrInternal, err)
}

func TestCategoriesRepository_GetAllByUserID_Success(t *testing.T) {
	db := setupCategoriesTestDB(t)
	repo := NewCategoriesRepository(db)
	ctx := context.Background()

	userID := "USER-001"
	insertTestUser(t, userID)

	date := time.Now().Format("20060102")

	category := &Categories{
		CATEGORY_ID: "CAT-" + date + "-000001",
		USER_ID:     userID,
		Name:        "Food",
		Description: "Food expenses",
	}

	require.NoError(t, repo.Create(ctx, category))

	result, err := repo.GetAllByUserID(ctx, userID)
	require.NoError(t, err)

	// cek 1 array
	require.Len(t, result, 1)

	// assert first data
	require.Equal(t, "Food", result[0].Name)
	require.Equal(t, "Food expenses", result[0].Description)
}

func TestCategoriesRepository_GetAllByUserID_WrongUser(t *testing.T) {
	db := setupCategoriesTestDB(t)
	repo := NewCategoriesRepository(db)
	ctx := context.Background()

	ownerID := "USER-001"
	otherUserID := "USER-002"

	insertTestUser(t, ownerID)
	insertTestUser(t, otherUserID)

	date := time.Now().Format("20060102")

	category := &Categories{
		CATEGORY_ID: "CAT-" + date + "-000001",
		USER_ID:     ownerID,
		Name:        "Food",
		Description: "Food expenses",
	}

	require.NoError(t, repo.Create(ctx, category))

	// attempt Get All by different user
	result, err := repo.GetAllByUserID(ctx, otherUserID)

	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestCategoriesRepository_GetByCategoryID_Success(t *testing.T) {
	db := setupCategoriesTestDB(t)
	repo := NewCategoriesRepository(db)
	ctx := context.Background()

	userID := "USER-001"
	insertTestUser(t, userID)

	date := time.Now().Format("20060102")

	category := &Categories{
		CATEGORY_ID: "CAT-" + date + "-000001",
		USER_ID:     userID,
		Name:        "Food",
		Description: "Food expenses",
	}

	require.NoError(t, repo.Create(ctx, category))

	result, err := repo.GetByCategoryID(ctx, userID, categories.CATEGORY_ID)
	require.NoError(t, err)

	require.NoError(t, err)
	require.NotNil(t, result)

	// assert object fields
	require.Equal(t, "Food", result.Name)
	require.Equal(t, "Food expenses", result.Description)
}

func TestCategoriesRepository_GetByCategoryID_WrongUserID(t *testing.T) {
	db := setupCategoriesTestDB(t)
	repo := NewCategoriesRepository(db)
	ctx := context.Background()

	ownerID := "USER-001"
	otherUserID := "USER-002"

	insertTestUser(t, ownerID)
	insertTestUser(t, otherUserID)

	date := time.Now().Format("20060102")

	category := &Categories{
		CATEGORY_ID: "CAT-" + date + "-000001",
		USER_ID:     ownerID,
		Name:        "Food",
		Description: "Food expenses",
	}

	require.NoError(t, repo.Create(ctx, category))

	result, err := repo.GetByCategoryID(ctx, otherUserID, category.CATEGORY_ID)
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrCategoryNotFound)
}

func TestCategoriesRepository_GetByCategoryID_CategoryNotFound(t *testing.T) {
	db := setupCategoriesTestDB(t)
	repo := NewCategoriesRepository(db)
	ctx := context.Background()

	ownerID := "USER-001"

	insertTestUser(t, ownerID)

	date := time.Now().Format("20060102")

	category := &Categories{
		CATEGORY_ID: "CAT-" + date + "-000001",
		USER_ID:     ownerID,
		Name:        "Food",
		Description: "Food expenses",
	}

	require.NoError(t, repo.Create(ctx, category))

	result, err := repo.GetByCategoryID(ctx, ownerID, "CAT-20250101-999999")
	require.Nil(t, result)
	require.ErrorIs(t, err, ErrCategoryNotFound)
}

func TestCategoriesRepository_Update_Success(t *testing.T) {
	db := setupCategoriesTestDB(t)
	repo := NewCategoriesRepository(db)
	ctx := context.Background()

	userID := "USER-001"
	insertTestUser(t, userID)

	date := time.Now().Format("20060102")

	category := &Categories{
		CATEGORY_ID: "CAT-" + date + "-000001",
		USER_ID:     userID,
		Name:        "Food",
		Description: "Food expenses",
	}

	require.NoError(t, repo.Create(ctx, category))

	// 2️⃣ Update field
	category.Name = "Food Updated"
	category.Description = "Updated desc"

	require.NoError(t, repo.Update(ctx, category))

	result, err := repo.GetByCategoryID(ctx, userID, category.CATEGORY_ID)
	require.NoError(t, err)

	// 4️⃣ Assert perubahan
	require.Equal(t, "Food Updated", result.Name)
	require.Equal(t, "Updated desc", result.Description)
}

func TestCategoriesRepository_Delete_Success(t *testing.T) {
	db := setupCategoriesTestDB(t)
	repo := NewCategoriesRepository(db)
	ctx := context.Background()

	userID := "USER-001"
	insertTestUser(t, userID)

	date := time.Now().Format("20060102")

	category := &Categories{
		CATEGORY_ID: "CAT-" + date + "-000001",
		USER_ID:     userID,
		Name:        "Food",
		Description: "Food expenses",
	}

	require.NoError(t, repo.Create(ctx, category))

	// delete
	err := repo.Delete(ctx, category.CATEGORY_ID, userID)
	require.NoError(t, err)

	// pastikan sudah hilang
	count, err := repo.CountByDate(ctx, date)
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func TestCategoriesRepository_Delete_NotFound(t *testing.T) {
	db := setupCategoriesTestDB(t)
	repo := NewCategoriesRepository(db)
	ctx := context.Background()

	userID := "USER-001"
	insertTestUser(t, userID)

	err := repo.Delete(ctx, "CAT-20250101-999999", userID)

	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestCategoriesRepository_Delete_WrongUser(t *testing.T) {
	db := setupCategoriesTestDB(t)
	repo := NewCategoriesRepository(db)
	ctx := context.Background()

	ownerID := "USER-001"
	otherUserID := "USER-002"

	insertTestUser(t, ownerID)
	insertTestUser(t, otherUserID)

	date := time.Now().Format("20060102")

	category := &Categories{
		CATEGORY_ID: "CAT-" + date + "-000001",
		USER_ID:     ownerID,
		Name:        "Food",
		Description: "Food expenses",
	}

	require.NoError(t, repo.Create(ctx, category))

	// attempt delete by different user
	err := repo.Delete(ctx, category.CATEGORY_ID, otherUserID)

	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestCategoriesRepository_Delete_Twice(t *testing.T) {
	db := setupCategoriesTestDB(t)
	repo := NewCategoriesRepository(db)
	ctx := context.Background()

	userID := "USER-001"
	insertTestUser(t, userID)

	date := time.Now().Format("20060102")

	category := &Categories{
		CATEGORY_ID: "CAT-" + date + "-000001",
		USER_ID:     userID,
		Name:        "Food",
		Description: "Food expenses",
	}

	require.NoError(t, repo.Create(ctx, category))

	// first delete
	require.NoError(t, repo.Delete(ctx, category.CATEGORY_ID, userID))

	// second delete should fail
	err := repo.Delete(ctx, category.CATEGORY_ID, userID)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)
}
