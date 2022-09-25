package usecase

import (
	"context"
	"time"

	"github.com/khuchuz/go-clean-arch/domain"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

// NewUserUsecase will create new an userUsecase object representation of domain.UserUsecase interface
func NewUserUsecase(a domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       a,
		contextTimeout: timeout,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (a *userUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.User, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.userRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}
	nextCursor = ""
	return
}

func (a *userUsecase) GetByID(c context.Context, id int64) (res domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.userRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	return
}

func (a *userUsecase) Update(c context.Context, ar *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.userRepo.Update(ctx, ar)
}

func (a *userUsecase) GetByEmail(c context.Context, email string) (res domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return
	}

	return
}

func (a *userUsecase) Signup(c context.Context, m *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedUser, _ := a.GetByEmail(ctx, m.Email)
	if existedUser != (domain.User{}) {
		return domain.ErrConflict
	}

	err = a.userRepo.Signup(ctx, m)
	return
}

func (a *userUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedUser, err := a.userRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedUser == (domain.User{}) {
		return domain.ErrNotFound
	}
	return a.userRepo.Delete(ctx, id)
}
