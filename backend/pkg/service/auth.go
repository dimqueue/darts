package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dimqueue/darts/pkg/config"
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/argon2"
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
	var err error
	user.Password, err = s.generatePasswordHash(user.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	var userId int64

	err = s.txManager.WithTransaction(func(tx *sqlx.Tx) error {
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

type Argon2Configuration struct {
	HashRaw    []byte
	Salt       []byte
	TimeCost   uint32
	MemoryCost uint32
	Threads    uint8
	KeyLength  uint32
}

func (s *AuthService) generatePasswordHash(password string) (string, error) {
	cfg := Argon2Configuration{
		TimeCost:   2,
		MemoryCost: 64 * 1024,
		Threads:    4,
		KeyLength:  32,
	}
	salt, err := s.generateCryptographicSalt(16)
	if err != nil {
		return "", fmt.Errorf("password hashing failed: %w", err)
	}
	cfg.Salt = salt
	cfg.HashRaw = argon2.IDKey(
		[]byte(password),
		cfg.Salt,
		cfg.TimeCost,
		cfg.MemoryCost,
		cfg.Threads,
		cfg.KeyLength)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		cfg.MemoryCost,
		cfg.TimeCost,
		cfg.Threads,
		base64.RawStdEncoding.EncodeToString(cfg.Salt),
		base64.RawStdEncoding.EncodeToString(cfg.HashRaw))

	return encodedHash, nil
}

func (s *AuthService) generateCryptographicSalt(saltSize uint32) ([]byte, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("salt generation failed: %w", err)
	}
	return salt, nil
}

func (s *AuthService) parseArgon2Hash(encodedHash string) (*Argon2Configuration, error) {
	components := strings.Split(encodedHash, "$")
	if len(components) != 6 {
		return nil, errors.New("invalid hash format structure")
	}

	if !strings.HasPrefix(components[1], "argon2id") {
		return nil, errors.New("unsupported algorithm variant")
	}

	var version int
	fmt.Sscanf(components[2], "v=%d", &version)

	cfg := &Argon2Configuration{}
	fmt.Sscanf(components[3], "m=%d,t=%d,p=%d",
		&cfg.MemoryCost, &cfg.TimeCost, &cfg.Threads)

	salt, err := base64.RawStdEncoding.DecodeString(components[4])
	if err != nil {
		return nil, fmt.Errorf("salt decoding failed: %w", err)
	}
	cfg.Salt = salt

	hash, err := base64.RawStdEncoding.DecodeString(components[5])
	if err != nil {
		return nil, fmt.Errorf("hash decoding failed: %w", err)
	}
	cfg.HashRaw = hash
	cfg.KeyLength = uint32(len(hash))

	return cfg, nil
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	res, err := s.verifyPasswordSecure(user.Password, password)
	if err != nil {
		return "", fmt.Errorf("failed to verify password: %w", err)
	}
	if !res {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.TokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(config.JWTSecret))
}

func (s *AuthService) verifyPasswordSecure(storedHash, providedPassword string) (bool, error) {
	cfg, err := s.parseArgon2Hash(storedHash)
	if err != nil {
		return false, fmt.Errorf("hash parsing failed: %w", err)
	}

	computedHash := argon2.IDKey(
		[]byte(providedPassword),
		cfg.Salt,
		cfg.TimeCost,
		cfg.MemoryCost,
		cfg.Threads,
		cfg.KeyLength,
	)

	match := subtle.ConstantTimeCompare(cfg.HashRaw, computedHash) == 1
	return match, nil
}

func (s *AuthService) ParseToken(accessToken string) (int64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.JWTSecret), nil
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
