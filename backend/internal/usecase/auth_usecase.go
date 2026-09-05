package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jaberpatwary/startech/internal/domain"
	"github.com/jaberpatwary/startech/internal/infrastructure/security"
)

type AuthUsecase interface {
	Register(ctx context.Context, name, email, password, phone string) (*domain.User, string, error)
	Login(ctx context.Context, email, password string) (*domain.User, string, error)
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID, name, phone, avatar string) (*domain.User, error)
	AddAddress(ctx context.Context, userID string, addr *domain.Address) error
	DeleteAddress(ctx context.Context, id, userID string) error
}

type authUsecase struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewAuthUsecase(userRepo domain.UserRepository, jwtSecret string) AuthUsecase {
	return &authUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *authUsecase) Register(ctx context.Context, name, email, password, phone string) (*domain.User, string, error) {
	existing, err := u.userRepo.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, "", domain.ErrAlreadyExists
	}

	hash, err := security.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user := &domain.User{
		ID:           uuid.NewString(),
		Name:         strings.TrimSpace(name),
		Email:        strings.ToLower(strings.TrimSpace(email)),
		Phone:        strings.TrimSpace(phone),
		PasswordHash: hash,
		Role:         domain.RoleUser,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := security.GenerateToken(user.ID, user.Email, user.Role, u.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", domain.ErrInvalidCredentials
	}

	if user.IsBlocked {
		return nil, "", domain.ErrAccountBlocked
	}

	if !security.CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", domain.ErrInvalidCredentials
	}

	token, err := security.GenerateToken(user.ID, user.Email, user.Role, u.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (u *authUsecase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	return u.userRepo.FindByID(ctx, userID)
}

func (u *authUsecase) UpdateProfile(ctx context.Context, userID, name, phone, avatar string) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if name != "" {
		user.Name = name
	}
	if phone != "" {
		user.Phone = phone
	}
	if avatar != "" {
		user.Avatar = avatar
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *authUsecase) AddAddress(ctx context.Context, userID string, addr *domain.Address) error {
	addr.ID = uuid.NewString()
	addr.UserID = userID
	return u.userRepo.AddAddress(ctx, addr)
}

func (u *authUsecase) DeleteAddress(ctx context.Context, id, userID string) error {
	return u.userRepo.DeleteAddress(ctx, id, userID)
}
