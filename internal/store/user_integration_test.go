package store

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	// _ "github.com/glebarez/go-sqlite"
	// "gorm.io/driver/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserTestSuite struct {
	suite.Suite

	dbname   string
	db       *gorm.DB
	userRepo UserRepository
}

func (s *UserTestSuite) SetupSuite() {
	dbhost := os.Getenv("TEST_DBHOST")
	if dbhost == "" {
		s.T().Skip("skip test: env TEST_DBHOST not set")
	}
	dbportStr := os.Getenv("TEST_DBPORT")
	dbport := 3306
	if dbportStr != "" {
		var err error
		dbport, err = strconv.Atoi(dbportStr)
		if err != nil {
			s.T().Skipf("skip test: parse dbport from ENV failed: %s", dbportStr)
		}
	}
	dbuser := os.Getenv("TEST_DBUSER")
	dbpassword := os.Getenv("TEST_DBPASSWORD")
	s.dbname = fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	// create test database
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql?charset=utf8mb4&parseTime=true&loc=Local",
		dbuser, dbpassword, dbhost, dbport)), &gorm.Config{})
	if err != nil {
		s.T().Skipf("skip test: open database failed: %v", err)
	}
	err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", s.dbname)).Error
	if err != nil {
		s.T().Skipf("skip test: create database failed: %v", err)
	}

	db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		dbuser, dbpassword, dbhost, dbport, s.dbname)), &gorm.Config{})
	if err != nil {
		s.T().Skipf("skip test: open database failed: %v", err)
	}

	// db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	// if err != nil {
	// 	s.T().Skipf("skip test: open database failed: %v", err)
	// }
	s.db = db
	s.userRepo = NewUserRepository(db)

	s.db.AutoMigrate(&User{})
}

func (s *UserTestSuite) TearDownSuite() {
	s.db.Exec(fmt.Sprintf("DROP DATABASE %s", s.dbname))
}

func (s *UserTestSuite) TestUser() {
	u := &User{
		ID:    uuid.NewString(),
		Name:  "liuhong",
		Email: "aaa@bb.com",
		Age:   22,
	}
	// Create
	err := s.userRepo.Create(u)
	s.Require().NoError(err)

	// GetByID
	user, err := s.userRepo.GetByID(u.ID)
	s.Require().NoError(err)
	if user.Name != u.Name || user.Email != u.Email {
		s.T().Fatalf("assert get user failed, expected: %v, got: %v", u, user)
	}

	// GetByEmail
	user, err = s.userRepo.GetByEmail(u.Email)
	s.Require().NoError(err)
	if user.Name != u.Name || user.ID != u.ID {
		s.T().Fatalf("assert get user failed, expected: %v, got: %v", u, user)
	}

	// Update
	u.Age = 23
	err = s.userRepo.Update(u)
	s.Require().NoError(err)
	user, err = s.userRepo.GetByID(u.ID)
	s.Require().NoError(err)
	s.Require().EqualValues(23, user.Age, "assert update user age failed")

	// List
	users, total, err := s.userRepo.List(1, 10)
	s.Require().NoError(err)
	s.Require().EqualValues(1, total, "asset list users total failed")
	s.Require().EqualValues(1, len(users), "asset list users total failed")
	s.Require().EqualValues(*user, users[0], "asset list users failed")

	// DeleteByID
	err = s.userRepo.DeleteByID(u.ID)
	s.Require().NoError(err)
	_, err = s.userRepo.GetByID(u.ID)
	s.Require().NotNil(err)
	s.Require().ErrorContains(err, "not found")
}

func TestUserIntegration(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}