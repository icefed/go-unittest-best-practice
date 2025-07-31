package store

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	testUserRepo UserRepository
	sqlMock      sqlmock.Sqlmock
)

func TestMain(m *testing.M) {
	mockDb, mock, _ := sqlmock.New()
	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      mockDb,
		SkipInitializeWithVersion: true,
	})
	newLogger := logger.New(
		log.New(os.Stderr, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("test open mock db failed: %v", err)
	}
	defer mockDb.Close()
	testUserRepo = NewUserRepository(db)
	sqlMock = mock

	m.Run()
}

func TestUserWithMock(t *testing.T) {
	// Create
	t.Run("Create", func(t *testing.T) {
		u := &User{
			ID:    uuid.NewString(),
			Name:  "liuhong",
			Email: "aaa@bb.com",
			Age:   22,
		}
		sqlMock.ExpectBegin()
		sqlMock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(1, 1))
		sqlMock.ExpectCommit()
		err := testUserRepo.Create(u)
		assert.NoError(t, err)
	})

	// GetByID
	t.Run("GetByID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "age", "created_at", "updated_at", "deleted_at"}).
			AddRow("idddddddd", "liuliu", "aa@bb.com", "", 10, time.Now(), time.Now(), sql.NullTime{})
		sqlMock.ExpectQuery(`SELECT`).WillReturnRows(rows)
		user, err := testUserRepo.GetByID("idddddddd")
		require.NoError(t, err)
		if user.Name != "liuliu" || user.Email != "aa@bb.com" {
			t.Fatalf("assert get user failed, got: %v", user)
		}
	})

	// GetByEmail
	t.Run("GetByEmail", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "age", "created_at", "updated_at", "deleted_at"}).
			AddRow("idddddddd", "liuliu", "aa@bb.com", "", 10, time.Now(), time.Now(), sql.NullTime{})
		sqlMock.ExpectQuery(`SELECT`).WillReturnRows(rows)
		user, err := testUserRepo.GetByID("aa@bb.com")
		require.NoError(t, err)
		if user.Name != "liuliu" || user.ID != "idddddddd" {
			t.Fatalf("assert get user failed, got: %v", user)
		}
	})

	// Update
	t.Run("Update", func(t *testing.T) {
		u := &User{
			ID:    uuid.NewString(),
			Name:  "liuhong",
			Email: "aaa@bb.com",
			Age:   22,
		}
		sqlMock.ExpectBegin()
		sqlMock.ExpectExec("UPDATE `users`").WillReturnResult(sqlmock.NewResult(1, 1))
		sqlMock.ExpectCommit()
		err := testUserRepo.Update(u)
		assert.NoError(t, err)
	})

	// List
	t.Run("List", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "age", "created_at", "updated_at", "deleted_at"}).
			AddRow("idddddddd", "liuliu", "aa@bb.com", "", 10, time.Now(), time.Now(), sql.NullTime{})
		sqlMock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
		sqlMock.ExpectQuery(`SELECT`).WillReturnRows(rows)
		users, total, err := testUserRepo.List(1, 10)
		require.NoError(t, err)
		require.EqualValues(t, 1, total, "asset list users total failed")
		require.EqualValues(t, 1, len(users), "asset list users total failed")
	})

	// DeleteByID
	t.Run("DeleteByID", func(t *testing.T) {
		sqlMock.ExpectBegin()
		sqlMock.ExpectExec("UPDATE `users`").WillReturnResult(sqlmock.NewResult(1, 1))
		sqlMock.ExpectCommit()
		err := testUserRepo.DeleteByID("idddddddddd")
		assert.NoError(t, err)
	})
}