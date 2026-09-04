package service

import (
	"app/src/model"
	"app/src/response"
	"app/src/utils"
	"app/src/validation"
	"errors"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	Register(req *validation.RegisterRequest) (*model.User, string, error)
	Login(req *validation.LoginRequest) (*model.User, string, error)
	GetProfile(userID string) (*model.User, error)
	UpdateProfile(userID string, req *validation.UpdateProfileRequest) (*model.User, error)
	GetAll(query *validation.QueryUser) (*response.SuccessWithPaginate[model.User], error)
}

type userServiceImpl struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userServiceImpl{db: db}
}

func (s *userServiceImpl) Register(req *validation.RegisterRequest) (*model.User, string, error) {
	var count int64
	s.db.Model(&model.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return nil, "", errors.New("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	newUser := model.User{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		Role:         model.RoleUser,
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateJWTToken(newUser.ID, newUser.Email, newUser.Role)
	if err != nil {
		return nil, "", err
	}

	return &newUser, token, nil
}

func (s *userServiceImpl) Login(req *validation.LoginRequest) (*model.User, string, error) {
	var user model.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if user.IsBlocked {
		return nil, "", errors.New("user account is blocked")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWTToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *userServiceImpl) GetProfile(userID string) (*model.User, error) {
	var user model.User
	if err := s.db.Preload("Addresses").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (s *userServiceImpl) UpdateProfile(userID string, req *validation.UpdateProfileRequest) (*model.User, error) {
	var user model.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userServiceImpl) GetAll(query *validation.QueryUser) (*response.SuccessWithPaginate[model.User], error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}

	var users []model.User
	var totalResults int64

	dbQuery := s.db.Model(&model.User{})
	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		dbQuery = dbQuery.Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	dbQuery.Count(&totalResults)

	offset := (query.Page - 1) * query.Limit
	if err := dbQuery.Offset(offset).Limit(query.Limit).Order("created_at desc").Find(&users).Error; err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return &response.SuccessWithPaginate[model.User]{
		Code:         200,
		Status:       "success",
		Message:      "Users fetched successfully",
		Results:      users,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}, nil
}
