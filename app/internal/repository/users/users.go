package users

import (
	"context"

	"gorm.io/gorm"

	"github.com/gunzgo2mars/test-poke-service/app/internal/core/model"
)

type IUserRepository interface {
	GetUser(ctx context.Context, scope ...func(db *gorm.DB) *gorm.DB) (*model.UserSchema, error)
	CreateNewUser(ctx context.Context, schema *model.UserSchema) error
}

type userRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUser(
	ctx context.Context,
	scope ...func(db *gorm.DB) *gorm.DB,
) (*model.UserSchema, error) {
	var userSchema *model.UserSchema

	if err := r.db.WithContext(ctx).Scopes(scope...).First(&userSchema).Error; err != nil {
		return nil, err
	}

	return userSchema, nil
}

func (r *userRepository) CreateNewUser(ctx context.Context, schema *model.UserSchema) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(&schema).Error; err != nil {
			return err
		}
		return nil
	})
}

func ByUsername(username string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("username = ?", username)
	}
}

func ByUUID(uuid string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("uuid = ?", uuid)
	}
}
