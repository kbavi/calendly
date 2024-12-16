package user

import (
	"context"

	"github.com/kbavi/calendly/pkg"
	"github.com/kbavi/calendly/repo"
)

type Service interface {
	Create(ctx context.Context, input *pkg.CreateUserInput) (*pkg.User, error)
	Get(ctx context.Context, id string) (*pkg.User, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	userRepo repo.UserRepository
}

func NewService(userRepo repo.UserRepository) Service {
	return &service{userRepo: userRepo}
}

func (s *service) Create(ctx context.Context, input *pkg.CreateUserInput) (*pkg.User, error) {
	return s.userRepo.Create(ctx, input)
}

func (s *service) Get(ctx context.Context, id string) (*pkg.User, error) {
	return s.userRepo.Get(ctx, id)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
