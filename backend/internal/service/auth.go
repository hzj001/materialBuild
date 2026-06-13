package service

import (
	"errors"

	"material-build/internal/model"
	"material-build/internal/repository"
	"material-build/pkg/config"
	"material-build/pkg/jwtutil"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewAuthService(repo *repository.Repository, cfg *config.Config) *AuthService {
	return &AuthService{repo: repo, cfg: cfg}
}

type LoginResult struct {
	Token    string      `json:"token"`
	User     *model.User `json:"user"`
}

func (s *AuthService) Login(phone, password, role string) (*LoginResult, error) {
	user, err := s.repo.FindUserByPhoneAndRole(phone, role)
	if err != nil {
		return nil, errors.New("账号或密码错误")
	}
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("账号或密码错误")
	}

	token, err := jwtutil.Generate(user.ID, user.Role, s.cfg.JWT.Secret, s.cfg.JWT.ExpireHours)
	if err != nil {
		return nil, err
	}
	return &LoginResult{Token: token, User: user}, nil
}

type RegisterInput struct {
	Phone    string
	Password string
	Nickname string
	Role     string
	CityID   *uint64
}

func (s *AuthService) Register(in RegisterInput) (*LoginResult, error) {
	if _, err := s.repo.FindUserByPhoneAndRole(in.Phone, in.Role); err == nil {
		return nil, errors.New("该手机号已注册")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Phone:    in.Phone,
		Password: string(hash),
		Nickname: in.Nickname,
		Role:     in.Role,
		CityID:   in.CityID,
		Status:   1,
	}
	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	token, err := jwtutil.Generate(user.ID, user.Role, s.cfg.JWT.Secret, s.cfg.JWT.ExpireHours)
	if err != nil {
		return nil, err
	}
	return &LoginResult{Token: token, User: user}, nil
}
