package service

import (
	"fmt"

	"github.com/google/uuid"

	"voicechat/internal/domain"
	"voicechat/internal/repository"
	"voicechat/internal/util"
)

type AuthService struct {
	userRepo *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtSecret: jwtSecret,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *UserDTO    `json:"user"`
}

type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	if len(req.Password) < 6 {
		return nil, fmt.Errorf("password must be at least 6 characters")
	}

	// Hash password
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := util.GenerateToken(user.ID, user.Username, user.Email, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User: &UserDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !util.VerifyPassword(user.PasswordHash, req.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := util.GenerateToken(user.ID, user.Username, user.Email, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User: &UserDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *AuthService) GetUser(userID string) (*UserDTO, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
