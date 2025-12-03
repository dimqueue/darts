package repository

import (
	"database/sql"
	"fmt"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/jmoiron/sqlx"
)

type ProfilePostgres struct {
	db *sqlx.DB
}

func NewProfilePostgres(db *sqlx.DB) *ProfilePostgres {
	return &ProfilePostgres{db: db}
}

func (r *ProfilePostgres) GetProfile(userId int64) (*model.UserProfile, error) {
	var profile model.UserProfile

	query := fmt.Sprintf(`
		SELECT user_id, avatar_url, bio, country_code, timezone, date_of_birth, created_at, updated_at
		FROM %s
		WHERE user_id = $1
	`, profilesTable)

	err := r.db.Get(&profile, query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, err
	}

	return &profile, nil
}

func (r *ProfilePostgres) CreateProfile(tx *sqlx.Tx, profile *model.UserProfile) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, avatar_url, bio, country_code, timezone, date_of_birth)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, profilesTable)

	_, err := tx.Exec(query,
		profile.UserId,
		profile.AvatarURL,
		profile.Bio,
		profile.CountryCode,
		profile.Timezone,
		profile.DateOfBirth,
	)

	return err
}

func (r *ProfilePostgres) UpdateProfile(userId int64, input model.UpdateProfileInput) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET avatar_url = COALESCE($2, avatar_url),
		    bio = COALESCE($3, bio),
		    country_code = COALESCE($4, country_code),
		    timezone = COALESCE($5, timezone)
		WHERE user_id = $1
	`, profilesTable)

	_, err := r.db.Exec(query,
		userId,
		input.AvatarURL,
		input.Bio,
		input.CountryCode,
		input.Timezone,
	)

	return err
}

func (r *ProfilePostgres) GetSettings(userId int64) (*model.UserSettings, error) {
	var settings model.UserSettings

	query := fmt.Sprintf(`
		SELECT user_id, preferred_language, theme, sound_enabled, email_notifications,
		       show_profile_public, show_stats_public, created_at, updated_at
		FROM %s
		WHERE user_id = $1
	`, settingsTable)

	err := r.db.Get(&settings, query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("settings not found")
		}
		return nil, err
	}

	return &settings, nil
}

func (r *ProfilePostgres) CreateSettings(tx *sqlx.Tx, userId int64) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id)
		VALUES ($1)
	`, settingsTable)

	_, err := tx.Exec(query, userId)
	return err
}

func (r *ProfilePostgres) UpdateSettings(userId int64, input model.UpdateSettingsInput) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET preferred_language = COALESCE($2, preferred_language),
		    theme = COALESCE($3, theme),
		    sound_enabled = COALESCE($4, sound_enabled),
		    email_notifications = COALESCE($5, email_notifications),
		    show_profile_public = COALESCE($6, show_profile_public),
		    show_stats_public = COALESCE($7, show_stats_public)
		WHERE user_id = $1
	`, settingsTable)

	_, err := r.db.Exec(query,
		userId,
		input.PreferredLanguage,
		input.Theme,
		input.SoundEnabled,
		input.EmailNotifications,
		input.ShowProfilePublic,
		input.ShowStatsPublic,
	)

	return err
}

func (r *ProfilePostgres) GetProfileSummary(userId int64) (*model.UserProfileSummary, error) {
	var summary model.UserProfileSummary

	query := `
		SELECT id, username, name, email, member_since, last_login_at,
		       avatar_url, bio, country_code,
		       preferred_language, theme, show_profile_public, show_stats_public,
		       total_games, total_wins, total_losses, current_win_streak,
		       best_win_streak, average_guesses, total_score, last_game_at
		FROM user_profile_summary
		WHERE id = $1
	`

	err := r.db.Get(&summary, query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &summary, nil
}

func (r *ProfilePostgres) GetProfileByUsername(username string) (*model.UserProfileSummary, error) {
	var summary model.UserProfileSummary

	query := `
		SELECT id, username, name, email, member_since, last_login_at,
		       avatar_url, bio, country_code,
		       preferred_language, theme, show_profile_public, show_stats_public,
		       total_games, total_wins, total_losses, current_win_streak,
		       best_win_streak, average_guesses, total_score, last_game_at
		FROM user_profile_summary
		WHERE username = $1
	`

	err := r.db.Get(&summary, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &summary, nil
}
