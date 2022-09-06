package mysql_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/khuchuz/go-clean-arch/user/repository"
	userMysqlRepo "github.com/khuchuz/go-clean-arch/user/repository/mysql"
	"github.com/khuchuz/go-clean-arch/domain"
)

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUsers := []domain.User{
		domain.User{
			ID: 1, Name: "title 1", Password: "Password 1", Email: "test1@gmail.com", UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
		domain.User{
			ID: 2, Name: "title 2", Password: "Password 2", Email: "test2@gmail.com", UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "password", "email", "updated_at", "created_at"}).
		AddRow(mockUsers[0].ID, mockUsers[0].Name, mockUsers[0].Password, mockUsers[0].Email,
			mockUsers[0].UpdatedAt, mockUsers[0].CreatedAt).
		AddRow(mockUsers[1].ID, mockUsers[1].Name, mockUsers[1].Password, mockUsers[1].Email,
			mockUsers[1].UpdatedAt, mockUsers[1].CreatedAt)

	query := "SELECT id,name,password,email, updated_at, created_at FROM user WHERE created_at > \\? ORDER BY created_at LIMIT \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := userMysqlRepo.NewMysqlUserRepository(db)
	cursor := repository.EncodeCursor(mockUsers[1].CreatedAt)
	num := int64(2)
	list, nextCursor, err := a.Fetch(context.TODO(), cursor, num)
	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "password", "email", "updated_at", "created_at"}).
		AddRow(1, "title 1", "Password 1", "test1@gmail.com", time.Now(), time.Now())

	query := "SELECT id,name,password,email, updated_at, created_at FROM user WHERE ID = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := userMysqlRepo.NewMysqlUserRepository(db)

	num := int64(5)
	anUser, err := a.GetByID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, anUser)
}

func TestStore(t *testing.T) {
	now := time.Now()
	ar := &domain.User{
		Name:      "Name",
		Email:     "Email",
		Password:  "Password",
		CreatedAt: now,
		UpdatedAt: now,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "INSERT  user SET name=\\? , password=\\? , email=\\? , updated_at=\\? , created_at=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Name, ar.Password, ar.Email, ar.CreatedAt, ar.UpdatedAt).WillReturnResult(sqlmock.NewResult(12, 1))

	a := userMysqlRepo.NewMysqlUserRepository(db)

	err = a.Store(context.TODO(), ar)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), ar.ID)
}

func TestGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "password", "email", "updated_at", "created_at"}).
		AddRow(1, "title 1", "Password 1", "test1@gmail.com", time.Now(), time.Now())

	query := "SELECT id,name,password,email, updated_at, created_at FROM user WHERE email = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := userMysqlRepo.NewMysqlUserRepository(db)

	email := "test1@gmail.com"
	anUser, err := a.GetByEmail(context.TODO(), email)
	assert.NoError(t, err)
	assert.NotNil(t, anUser)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "DELETE FROM user WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(12).WillReturnResult(sqlmock.NewResult(12, 1))

	a := userMysqlRepo.NewMysqlUserRepository(db)

	num := int64(12)
	err = a.Delete(context.TODO(), num)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	ar := &domain.User{
		ID:        12,
		Name:      "Name",
		Email:     "Email",
		Password:  "Password",
		CreatedAt: now,
		UpdatedAt: now,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "UPDATE user set name=\\?, password=\\?, email=\\?, updated_at=\\? WHERE ID = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Name, ar.Password, ar.Email, ar.UpdatedAt, ar.ID).WillReturnResult(sqlmock.NewResult(12, 1))

	a := userMysqlRepo.NewMysqlUserRepository(db)

	err = a.Update(context.TODO(), ar)
	assert.NoError(t, err)
}
