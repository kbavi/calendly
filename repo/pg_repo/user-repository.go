package pg_repo

import (
	"context"

	"github.com/kbavi/calendly/db/models"
	"github.com/kbavi/calendly/pkg"
	"github.com/kbavi/calendly/repo"
	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) repo.UserRepository {
	db.AutoMigrate(&models.UserModel{})
	return &userRepository{db: db}
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) Create(ctx context.Context, input *pkg.CreateUserInput) (*pkg.User, error) {
	id := repo.GenerateID()

	user := &models.UserModel{
		ID:    id,
		Email: string(input.Email),
		Name:  input.Name,
	}
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return &pkg.User{
		ID:    user.ID,
		Email: pkg.Email(user.Email),
		Name:  user.Name,
	}, nil
}

func (r *userRepository) Get(ctx context.Context, id string) (*pkg.User, error) {
	user := &models.UserModel{}
	if err := r.db.First(user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &pkg.User{
		ID:    user.ID,
		Email: pkg.Email(user.Email),
		Name:  user.Name,
	}, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.Unscoped().Delete(&models.UserModel{}, "id = ?", id).Error
}
