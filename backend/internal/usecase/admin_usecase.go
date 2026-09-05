package usecase

import (
	"context"

	"github.com/jaberpatwary/startech/internal/domain"
)

type AdminUsecase interface {
	GetDashboardStats(ctx context.Context) (map[string]interface{}, error)
	GetUsers(ctx context.Context, page, limit int) ([]domain.User, int64, error)
	ToggleBlockUser(ctx context.Context, id string) (*domain.User, error)
	GetSettings(ctx context.Context) (map[string]interface{}, error)
	UpdateSetting(ctx context.Context, key string, val interface{}) error
	GetReviews(ctx context.Context, status string, page, limit int) ([]domain.Review, int64, error)
	UpdateReviewStatus(ctx context.Context, id, status string) error
}

type adminUsecase struct {
	adminRepo domain.AdminRepository
	userRepo  domain.UserRepository
}

func NewAdminUsecase(adminRepo domain.AdminRepository, userRepo domain.UserRepository) AdminUsecase {
	return &adminUsecase{
		adminRepo: adminRepo,
		userRepo:  userRepo,
	}
}

func (u *adminUsecase) GetDashboardStats(ctx context.Context) (map[string]interface{}, error) {
	return u.adminRepo.GetDashboardStats(ctx)
}

func (u *adminUsecase) GetUsers(ctx context.Context, page, limit int) ([]domain.User, int64, error) {
	return u.userRepo.GetAll(ctx, page, limit)
}

func (u *adminUsecase) ToggleBlockUser(ctx context.Context, id string) (*domain.User, error) {
	return u.userRepo.ToggleBlock(ctx, id)
}

func (u *adminUsecase) GetSettings(ctx context.Context) (map[string]interface{}, error) {
	return u.adminRepo.GetSettings(ctx)
}

func (u *adminUsecase) UpdateSetting(ctx context.Context, key string, val interface{}) error {
	return u.adminRepo.UpdateSetting(ctx, key, val)
}

func (u *adminUsecase) GetReviews(ctx context.Context, status string, page, limit int) ([]domain.Review, int64, error) {
	return u.adminRepo.GetAllReviews(ctx, status, page, limit)
}

func (u *adminUsecase) UpdateReviewStatus(ctx context.Context, id, status string) error {
	return u.adminRepo.UpdateReviewStatus(ctx, id, status)
}
