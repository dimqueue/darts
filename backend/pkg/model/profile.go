package model

import "time"

type UserProfile struct {
	UserId      int64      `json:"user_id" db:"user_id"`
	AvatarURL   *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	Bio         *string    `json:"bio,omitempty" db:"bio"`
	CountryCode *string    `json:"country_code,omitempty" db:"country_code"`
	Timezone    string     `json:"timezone" db:"timezone"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" db:"date_of_birth"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type UserSettings struct {
	UserId             int64     `json:"user_id" db:"user_id"`
	PreferredLanguage  string    `json:"preferred_language" db:"preferred_language"`
	Theme              string    `json:"theme" db:"theme"` // "light", "dark", "auto"
	SoundEnabled       bool      `json:"sound_enabled" db:"sound_enabled"`
	EmailNotifications bool      `json:"email_notifications" db:"email_notifications"`
	ShowProfilePublic  bool      `json:"show_profile_public" db:"show_profile_public"`
	ShowStatsPublic    bool      `json:"show_stats_public" db:"show_stats_public"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type UpdateProfileInput struct {
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	CountryCode *string `json:"country_code,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
}

type UpdateSettingsInput struct {
	PreferredLanguage  *string `json:"preferred_language,omitempty"`
	Theme              *string `json:"theme,omitempty"`
	SoundEnabled       *bool   `json:"sound_enabled,omitempty"`
	EmailNotifications *bool   `json:"email_notifications,omitempty"`
	ShowProfilePublic  *bool   `json:"show_profile_public,omitempty"`
	ShowStatsPublic    *bool   `json:"show_stats_public,omitempty"`
}

type UserProfileSummary struct {
	Id          int64      `json:"id" db:"id"`
	Username    string     `json:"username" db:"username"`
	Name        string     `json:"name" db:"name"`
	Email       *string    `json:"email,omitempty" db:"email"`
	MemberSince time.Time  `json:"member_since" db:"member_since"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`

	AvatarURL   *string `json:"avatar_url,omitempty" db:"avatar_url"`
	Bio         *string `json:"bio,omitempty" db:"bio"`
	CountryCode *string `json:"country_code,omitempty" db:"country_code"`

	PreferredLanguage string `json:"preferred_language" db:"preferred_language"`
	Theme             string `json:"theme" db:"theme"`
	ShowProfilePublic bool   `json:"show_profile_public" db:"show_profile_public"`
	ShowStatsPublic   bool   `json:"show_stats_public" db:"show_stats_public"`

	TotalGames       int        `json:"total_games" db:"total_games"`
	TotalWins        int        `json:"total_wins" db:"total_wins"`
	TotalLosses      int        `json:"total_losses" db:"total_losses"`
	CurrentWinStreak int        `json:"current_win_streak" db:"current_win_streak"`
	BestWinStreak    int        `json:"best_win_streak" db:"best_win_streak"`
	AverageGuesses   float64    `json:"average_guesses" db:"average_guesses"`
	TotalScore       int        `json:"total_score" db:"total_score"`
	LastGameAt       *time.Time `json:"last_game_at,omitempty" db:"last_game_at"`
}
