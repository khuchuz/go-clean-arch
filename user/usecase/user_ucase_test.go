package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/khuchuz/go-clean-arch/domain"
	"github.com/khuchuz/go-clean-arch/domain/mocks"
	ucase "github.com/khuchuz/go-clean-arch/user/usecase"
)

func TestFetch(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := domain.User{
		Name:  "Name",
		Email: "Email",
	}

	mockListArtilce := make([]domain.User, 0)
	mockListArtilce = append(mockListArtilce, mockUser)

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(mockListArtilce, "next-cursor", nil).Once()
		u := ucase.NewUserUsecase(mockUserRepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)
		cursorExpected := "next-cursor"
		assert.Equal(t, cursorExpected, nextCursor)
		assert.NotEmpty(t, nextCursor)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListArtilce))

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockUserRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(nil, "", errors.New("Unexpexted Error")).Once()

		u := ucase.NewUserUsecase(mockUserRepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)

		assert.Empty(t, nextCursor)
		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockUserRepo.AssertExpectations(t)
	})

}

func TestGetByID(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := domain.User{
		Name:  "Name",
		Email: "Email",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockUser, nil).Once()
		u := ucase.NewUserUsecase(mockUserRepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockUser.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockUserRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.User{}, errors.New("Unexpected")).Once()
		u := ucase.NewUserUsecase(mockUserRepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockUser.ID)

		assert.Error(t, err)
		assert.Equal(t, domain.User{}, a)

		mockUserRepo.AssertExpectations(t)
	})

}

func TestStore(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := domain.User{
		Name:  "Name",
		Email: "Email",
	}

	t.Run("success", func(t *testing.T) {
		tempMockUser := mockUser
		tempMockUser.ID = 0
		mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{}, domain.ErrNotFound).Once()
		mockUserRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		u := ucase.NewUserUsecase(mockUserRepo, time.Second*2)

		err := u.Store(context.TODO(), &tempMockUser)

		assert.NoError(t, err)
		assert.Equal(t, mockUser.Email, tempMockUser.Email)
		mockUserRepo.AssertExpectations(t)
	})
	t.Run("existing-email", func(t *testing.T) {
		existingUser := mockUser
		mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(existingUser, nil).Once()
		u := ucase.NewUserUsecase(mockUserRepo, time.Second*2)

		err := u.Store(context.TODO(), &mockUser)

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})

}

func TestDelete(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := domain.User{
		Name:  "Name",
		Email: "Email",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockUser, nil).Once()

		mockUserRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		u := ucase.NewUserUsecase(mockUserRepo, time.Second*2)

		err := u.Delete(context.TODO(), mockUser.ID)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})
	t.Run("user-is-not-exist", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.User{}, nil).Once()

		u := ucase.NewUserUsecase(mockUserRepo, time.Second*2)

		err := u.Delete(context.TODO(), mockUser.ID)

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
	t.Run("error-happens-in-db", func(t *testing.T) {
		mockUserRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.User{}, errors.New("Unexpected Error")).Once()

		u := ucase.NewUserUsecase(mockUserRepo, time.Second*2)

		err := u.Delete(context.TODO(), mockUser.ID)

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})

}

func TestUpdate(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUser := domain.User{
		Name:  "Name",
		Email: "Email",
		ID:    23,
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Update", mock.Anything, &mockUser).Once().Return(nil)

		u := ucase.NewUserUsecase(mockUserRepo, time.Second*2)

		err := u.Update(context.TODO(), &mockUser)
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}
