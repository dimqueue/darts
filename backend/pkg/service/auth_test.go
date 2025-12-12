package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dimqueue/darts/pkg/config"
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/dimqueue/darts/pkg/repository/mocks"
)

func init() {
	config.JWTSecret = "test-secret-key-for-testing-purposes"
	config.TokenTTL = 12 * time.Hour
}

func newTestAuthService(
	authRepo *mocks.MockAuthRepository,
	profileRepo *mocks.MockProfileRepository,
	statsRepo *mocks.MockStatisticsRepository,
) *AuthService {
	txManager := &repository.TransactionManager{}

	statsService := &StatsService{
		statsRepo: statsRepo,
	}

	return &AuthService{
		repo:         authRepo,
		profileRepo:  profileRepo,
		statsService: statsService,
		txManager:    txManager,
	}
}

func TestGeneratePasswordHash(t *testing.T) {
	s := &AuthService{}

	hash, err := s.generatePasswordHash("testpassword123")
	if err != nil {
		t.Fatalf("generatePasswordHash failed: %v", err)
	}

	if len(hash) == 0 {
		t.Error("hash should not be empty")
	}

	if hash[:9] != "$argon2id" {
		t.Errorf("hash should start with $argon2id, got: %s", hash[:20])
	}
}

func TestGeneratePasswordHash_DifferentSalts(t *testing.T) {
	s := &AuthService{}

	hash1, _ := s.generatePasswordHash("samepassword")
	hash2, _ := s.generatePasswordHash("samepassword")

	if hash1 == hash2 {
		t.Error("same password should produce different hashes due to random salt")
	}
}

func TestVerifyPasswordSecure_CorrectPassword(t *testing.T) {
	s := &AuthService{}

	password := "correctpassword123"
	hash, err := s.generatePasswordHash(password)
	if err != nil {
		t.Fatalf("failed to generate hash: %v", err)
	}

	match, err := s.verifyPasswordSecure(hash, password)
	if err != nil {
		t.Fatalf("verifyPasswordSecure failed: %v", err)
	}
	if !match {
		t.Error("correct password should match")
	}
}

func TestVerifyPasswordSecure_WrongPassword(t *testing.T) {
	s := &AuthService{}

	password := "correctpassword123"
	hash, err := s.generatePasswordHash(password)
	if err != nil {
		t.Fatalf("failed to generate hash: %v", err)
	}

	match, err := s.verifyPasswordSecure(hash, "wrongpassword")
	if err != nil {
		t.Fatalf("verifyPasswordSecure failed: %v", err)
	}
	if match {
		t.Error("wrong password should not match")
	}
}

func TestVerifyPasswordSecure_InvalidHashFormat(t *testing.T) {
	s := &AuthService{}

	_, err := s.verifyPasswordSecure("invalid-hash-format", "anypassword")
	if err == nil {
		t.Error("should return error for invalid hash format")
	}
}

func TestParseArgon2Hash_ValidHash(t *testing.T) {
	s := &AuthService{}

	hash, _ := s.generatePasswordHash("testpassword")

	cfg, err := s.parseArgon2Hash(hash)
	if err != nil {
		t.Fatalf("parseArgon2Hash failed: %v", err)
	}

	if cfg.TimeCost != 2 {
		t.Errorf("expected TimeCost 2, got %d", cfg.TimeCost)
	}
	if cfg.MemoryCost != 64*1024 {
		t.Errorf("expected MemoryCost %d, got %d", 64*1024, cfg.MemoryCost)
	}
	if cfg.Threads != 4 {
		t.Errorf("expected Threads 4, got %d", cfg.Threads)
	}
	if len(cfg.Salt) == 0 {
		t.Error("salt should not be empty")
	}
	if len(cfg.HashRaw) == 0 {
		t.Error("hash should not be empty")
	}
}

func TestParseArgon2Hash_InvalidFormat(t *testing.T) {
	s := &AuthService{}

	testCases := []struct {
		name string
		hash string
	}{
		{"empty", ""},
		{"wrong parts", "$argon2id$v=19"},
		{"wrong algorithm", "$bcrypt$v=19$m=65536,t=2,p=4$salt$hash"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.parseArgon2Hash(tc.hash)
			if err == nil {
				t.Error("should return error for invalid hash format")
			}
		})
	}
}

func TestParseToken_ValidToken(t *testing.T) {
	s := &AuthService{}

	mockAuth := &mocks.MockAuthRepository{
		GetUserByUsernameFn: func(ctx context.Context, username string) (model.User, error) {
			hash, _ := s.generatePasswordHash("testpassword")
			return model.User{
				Id:       42,
				Username: username,
				Password: hash,
			}, nil
		},
	}

	authService := &AuthService{repo: mockAuth}

	token, err := authService.GenerateToken(context.Background(), "testuser", "testpassword")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	userId, err := s.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	if userId != 42 {
		t.Errorf("expected userId 42, got %d", userId)
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	s := &AuthService{}

	testCases := []struct {
		name  string
		token string
	}{
		{"empty", ""},
		{"garbage", "not.a.valid.token"},
		{"wrong signature", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0Mn0.invalid_signature"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.ParseToken(tc.token)
			if !errors.Is(err, ErrUnauthorized) {
				t.Errorf("expected ErrUnauthorized, got: %v", err)
			}
		})
	}
}

func TestGenerateToken_Success(t *testing.T) {
	s := &AuthService{}
	password := "testpassword123"
	hash, _ := s.generatePasswordHash(password)

	mockAuth := &mocks.MockAuthRepository{
		GetUserByUsernameFn: func(ctx context.Context, username string) (model.User, error) {
			return model.User{
				Id:       1,
				Username: username,
				Password: hash,
			}, nil
		},
	}

	authService := &AuthService{repo: mockAuth}

	token, err := authService.GenerateToken(context.Background(), "testuser", password)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("token should not be empty")
	}
}

func TestGenerateToken_UserNotFound(t *testing.T) {
	mockAuth := &mocks.MockAuthRepository{
		GetUserByUsernameFn: func(ctx context.Context, username string) (model.User, error) {
			return model.User{}, errors.New("user not found")
		},
	}

	authService := &AuthService{repo: mockAuth}

	_, err := authService.GenerateToken(context.Background(), "nonexistent", "anypassword")
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got: %v", err)
	}
}

func TestGenerateToken_WrongPassword(t *testing.T) {
	s := &AuthService{}
	hash, _ := s.generatePasswordHash("correctpassword")

	mockAuth := &mocks.MockAuthRepository{
		GetUserByUsernameFn: func(ctx context.Context, username string) (model.User, error) {
			return model.User{
				Id:       1,
				Username: username,
				Password: hash,
			}, nil
		},
	}

	authService := &AuthService{repo: mockAuth}

	_, err := authService.GenerateToken(context.Background(), "testuser", "wrongpassword")
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got: %v", err)
	}
}

func TestCreateUser_Success(t *testing.T) {
	var createdUserId int64 = 123
	var profileCreated, settingsCreated, statsCreated bool

	mockAuth := &mocks.MockAuthRepository{
		CreateUserFn: func(ctx context.Context, q repository.Querier, user model.User) (int64, error) {
			return createdUserId, nil
		},
	}

	mockProfile := &mocks.MockProfileRepository{
		CreateProfileFn: func(ctx context.Context, q repository.Querier, profile *model.UserProfile) error {
			profileCreated = true
			if profile.UserId != createdUserId {
				t.Errorf("expected profile.UserId %d, got %d", createdUserId, profile.UserId)
			}
			return nil
		},
		CreateSettingsFn: func(ctx context.Context, q repository.Querier, userId int64) error {
			settingsCreated = true
			if userId != createdUserId {
				t.Errorf("expected userId %d, got %d", createdUserId, userId)
			}
			return nil
		},
	}

	mockStats := &mocks.MockStatisticsRepository{
		CreateGlobalStreaksFn: func(ctx context.Context, q repository.Querier, userId int64) error {
			statsCreated = true
			if userId != createdUserId {
				t.Errorf("expected userId %d, got %d", createdUserId, userId)
			}
			return nil
		},
	}

	statsService := &StatsService{statsRepo: mockStats}

	authService := &AuthService{
		repo:         mockAuth,
		profileRepo:  mockProfile,
		statsService: statsService,
	}

	user := model.User{Username: "testuser", Password: "password123"}
	_, err := authService.generatePasswordHash(user.Password)
	if err != nil {
		t.Fatalf("password hashing failed: %v", err)
	}

	ctx := context.Background()
	userId, err := mockAuth.CreateUser(ctx, nil, user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if userId != createdUserId {
		t.Errorf("expected userId %d, got %d", createdUserId, userId)
	}

	err = mockProfile.CreateProfile(ctx, nil, &model.UserProfile{UserId: userId})
	if err != nil {
		t.Fatalf("CreateProfile failed: %v", err)
	}

	err = mockProfile.CreateSettings(ctx, nil, userId)
	if err != nil {
		t.Fatalf("CreateSettings failed: %v", err)
	}

	err = mockStats.CreateGlobalStreaks(ctx, nil, userId)
	if err != nil {
		t.Fatalf("CreateGlobalStreaks failed: %v", err)
	}

	if !profileCreated {
		t.Error("profile should have been created")
	}
	if !settingsCreated {
		t.Error("settings should have been created")
	}
	if !statsCreated {
		t.Error("stats should have been created")
	}
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	mockAuth := &mocks.MockAuthRepository{
		CreateUserFn: func(ctx context.Context, q repository.Querier, user model.User) (int64, error) {
			return 0, repository.ErrDuplicateKey
		},
	}

	mockProfile := &mocks.MockProfileRepository{}
	mockStats := &mocks.MockStatisticsRepository{}

	ctx := context.Background()
	_, err := mockAuth.CreateUser(ctx, nil, model.User{Username: "existing"})
	if !errors.Is(err, repository.ErrDuplicateKey) {
		t.Errorf("expected ErrDuplicateKey, got: %v", err)
	}

	_ = mockProfile
	_ = mockStats
}

func TestCreateUser_FullFlow(t *testing.T) {
	var createdUserId int64 = 456
	callOrder := []string{}

	mockAuth := &mocks.MockAuthRepository{
		CreateUserFn: func(ctx context.Context, q repository.Querier, user model.User) (int64, error) {
			callOrder = append(callOrder, "CreateUser")
			if len(user.Password) < 50 || user.Password[:9] != "$argon2id" {
				t.Error("password should be hashed with argon2id")
			}
			return createdUserId, nil
		},
	}

	mockProfile := &mocks.MockProfileRepository{
		CreateProfileFn: func(ctx context.Context, q repository.Querier, profile *model.UserProfile) error {
			callOrder = append(callOrder, "CreateProfile")
			return nil
		},
		CreateSettingsFn: func(ctx context.Context, q repository.Querier, userId int64) error {
			callOrder = append(callOrder, "CreateSettings")
			return nil
		},
	}

	mockStats := &mocks.MockStatisticsRepository{
		CreateGlobalStreaksFn: func(ctx context.Context, q repository.Querier, userId int64) error {
			callOrder = append(callOrder, "CreateGlobalStreaks")
			return nil
		},
	}

	statsService := &StatsService{statsRepo: mockStats}

	ctx := context.Background()

	s := &AuthService{}
	hashedPassword, err := s.generatePasswordHash("rawpassword")
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}

	user := model.User{Username: "newuser", Password: hashedPassword}
	userId, err := mockAuth.CreateUser(ctx, nil, user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if userId != createdUserId {
		t.Errorf("expected userId %d, got %d", createdUserId, userId)
	}

	err = mockProfile.CreateProfile(ctx, nil, &model.UserProfile{UserId: userId, Timezone: "UTC"})
	if err != nil {
		t.Fatalf("CreateProfile failed: %v", err)
	}

	err = mockProfile.CreateSettings(ctx, nil, userId)
	if err != nil {
		t.Fatalf("CreateSettings failed: %v", err)
	}

	err = statsService.InitializeStats(ctx, nil, userId)
	if err != nil {
		t.Fatalf("InitializeStats failed: %v", err)
	}

	expectedOrder := []string{"CreateUser", "CreateProfile", "CreateSettings", "CreateGlobalStreaks"}
	if len(callOrder) != len(expectedOrder) {
		t.Errorf("expected %d calls, got %d", len(expectedOrder), len(callOrder))
	}
	for i, expected := range expectedOrder {
		if i >= len(callOrder) || callOrder[i] != expected {
			t.Errorf("call %d: expected %s, got %s", i, expected, callOrder[i])
		}
	}
}
