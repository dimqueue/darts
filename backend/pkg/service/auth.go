package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
)

const (
	salt      = "llll"
	signInKey = "klkkk"
	tokenTTL  = 12 * time.Hour
)

type AuthService struct {
	repo         repository.Authorization
	profileRepo  repository.Profile
	statsService *StatsService
	txManager    *repository.TransactionManager
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int64 `json:"user_id"`
}

func NewAuthService(repo repository.Authorization, profileRepo repository.Profile, statsService *StatsService, txManager *repository.TransactionManager) *AuthService {
	return &AuthService{
		repo:         repo,
		profileRepo:  profileRepo,
		statsService: statsService,
		txManager:    txManager,
	}
}

func (s *AuthService) CreateUser(user model.User) (int64, error) {
	user.Password = s.generatePasswordHash(user.Password)

	var userId int64

	err := s.txManager.WithTransaction(func(tx *sqlx.Tx) error {
		id, err := s.repo.CreateUser(tx, user)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		userId = id

		profile := &model.UserProfile{
			UserId:   userId,
			Timezone: "UTC",
		}
		if err := s.profileRepo.CreateProfile(tx, profile); err != nil {
			return fmt.Errorf("failed to create profile: %w", err)
		}

		if err := s.profileRepo.CreateSettings(tx, userId); err != nil {
			return fmt.Errorf("failed to create settings: %w", err)
		}

		if err := s.statsService.InitializeStats(tx, userId); err != nil {
			return fmt.Errorf("failed to create statistics: %w", err)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (s *AuthService) ParseToken(accessToken string) (int64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signInKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserId, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	password = password + salt
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "failed to find user with this username and password", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(signInKey))
}
