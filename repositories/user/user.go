package repositories

import (
	"context"
	"errors"
	"user-service/domain/dto"
	"user-service/domain/models"

	errWrap "user-service/common/error"
	errConstant "user-service/constants/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	Register(context.Context, *dto.RegisterRequest) (*models.User, error)
	Update(context.Context, *dto.UpdateRequest, string) (*models.User, error)
	FindByUsername(context.Context, string) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
	FindByUUID(context.Context, string) (*models.User, error)
	FindByIDWithRole(ctx context.Context, id uint) (*models.User, error)
	// Preload(column string) *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	user := models.User{
		UUID:        uuid.New(),
		Name:        req.Name,
		UserName:    req.UserName,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		RoleID:      req.RoleID,
	}

	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*models.User, error) {
	user := models.User{
		Name:        req.Name,
		UserName:    req.Username,
		Password:    *req.Password,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
	}

	err := r.db.WithContext(ctx).Model(&user).
		Where("uuid = ?", uuid).
		Updates(user).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("user_name = ?", username).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &user, nil
}

func (r *UserRepository) FindByUUID(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("uuid = ?", uuid).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errConstant.ErrUserNotFound
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &user, nil
}

func (r *UserRepository) FindByIDWithRole(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// func (r *UserRepository) Preload(column string) *gorm.DB {
// 	return r.db.Preload(column)
// }
