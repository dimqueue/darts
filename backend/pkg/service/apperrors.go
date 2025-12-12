package service

import "errors"

var (
	// Game errors
	ErrGameNotFound  = errors.New("game not found")
	ErrGameNotActive = errors.New("game is not active")
	ErrGameExpired   = errors.New("game has expired")

	// Guess errors
	ErrWordNotFound    = errors.New("word not found in vocabulary")
	ErrWordAlreadyUsed = errors.New("word already guessed")
	ErrWordTooFar      = errors.New("word too far from target")

	// Auth errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrUserExists   = errors.New("user already exists")

	// Access errors
	ErrForbidden       = errors.New("access denied")
	ErrProfilePrivate  = errors.New("profile is private")
	ErrProfileNotFound = errors.New("profile not found")

	// Generic errors
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")

	// External service errors
	ErrComputeService = errors.New("compute service unavailable")
)

type LoggedError struct {
	Err error
}

func (e *LoggedError) Error() string {
	return e.Err.Error()
}

func (e *LoggedError) Unwrap() error {
	return e.Err
}

func Logged(err error) error {
	if err == nil {
		return nil
	}
	// Don't double-wrap
	if IsLogged(err) {
		return err
	}
	return &LoggedError{Err: err}
}

func IsLogged(err error) bool {
	var le *LoggedError
	return errors.As(err, &le)
}

func IsDomainError(err error) bool {
	switch {
	case errors.Is(err, ErrGameNotFound),
		errors.Is(err, ErrGameNotActive),
		errors.Is(err, ErrGameExpired),
		errors.Is(err, ErrWordNotFound),
		errors.Is(err, ErrWordAlreadyUsed),
		errors.Is(err, ErrWordTooFar),
		errors.Is(err, ErrUnauthorized),
		errors.Is(err, ErrUserExists),
		errors.Is(err, ErrForbidden),
		errors.Is(err, ErrProfilePrivate),
		errors.Is(err, ErrProfileNotFound),
		errors.Is(err, ErrNotFound),
		errors.Is(err, ErrInvalidInput),
		errors.Is(err, ErrComputeService):
		return true
	}
	return false
}
