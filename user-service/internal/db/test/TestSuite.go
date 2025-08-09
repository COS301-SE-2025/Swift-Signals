package test

import (
	"database/sql"
	"log"
	"time"

	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/db"
	"github.com/COS301-SE-2025/Swift-Signals/user-service/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo db.UserRepository
}

func (suite *TestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	suite.Require().NoError(err)

	suite.repo = db.NewPostgresUserRepo(suite.db)
}

func (suite *TestSuite) TearDownTest() {
	suite.mock.ExpectClose()
	if err := suite.db.Close(); err != nil {
		log.Printf("Failed to close suite: %v", err)
	}
}

var testUser = &model.User{
	ID:        uuid.NewString(),
	Name:      "Test Name",
	Email:     "test@gmail.com",
	Password:  "8characters",
	IsAdmin:   false,
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

var testUsers = []*model.User{
	{
		ID:        "test_user-1",
		Name:      "Test Name 1",
		Email:     "test1@gmail.com",
		Password:  "8characters",
		IsAdmin:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        "test_user-2",
		Name:      "Test Name 2",
		Email:     "test2@gmail.com",
		Password:  "8characters",
		IsAdmin:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

var testIntersectionIDs = []string{"int-1", "int-2", "int-3"}

var (
	getUserByEmailQuery = `SELECT uuid, name, email, password, is_admin, created_at, updated_at
FROM users
WHERE email = \$1`
	getUserByIDQuery = `SELECT uuid, name, email, password, is_admin, created_at, updated_at
	          FROM users
	          WHERE uuid = \$1`

	insertUserQuery = `INSERT INTO users \(uuid, name, email, password, is_admin, created_at, updated_at\)
	VALUES \(\$1, \$2, \$3, \$4, \$5, NOW\(\), NOW\(\)\)`
	getIntersectionIDQuery = `SELECT intersection_id FROM user_intersections WHERE user_id = \$1`
	updateUserQuery        = `UPDATE users
	          SET name = \$1, email = \$2, password = \$3, is_admin = \$4, updated_at = NOW\(\)
	          WHERE uuid = \$5`
	deleteUserQuery = `DELETE FROM users
	          WHERE uuid = \$1`
	listUsersQuery = `SELECT uuid, name, email, password, is_admin, created_at, updated_at
	          FROM users
	          ORDER BY uuid LIMIT \$1 OFFSET \$2`
)
