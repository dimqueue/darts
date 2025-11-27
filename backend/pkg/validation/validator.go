package validation

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

var (
	ErrBadWord   = errors.New("input contains prohibited words")
	ErrTooShort  = errors.New("input is too short")
	ErrEmptyWord = errors.New("input cannot be empty")
)

type Validator struct {
	badWords map[string]struct{}
}

func New() *Validator {
	return &Validator{
		badWords: loadBadWords(),
	}
}

func (v *Validator) ValidateWord(word string) error {
	word = strings.ToLower(strings.TrimSpace(word))

	if word == "" {
		return ErrEmptyWord
	}

	if utf8.RuneCountInString(word) < 2 {
		return ErrTooShort
	}

	if v.containsBadWord(word) {
		return ErrBadWord
	}

	return nil
}

func (v *Validator) containsBadWord(word string) bool {
	_, exists := v.badWords[word]
	return exists
}

func (v *Validator) RegisterCustomValidators(validate *validator.Validate) {
	validate.RegisterValidation("nobadwords", func(fl validator.FieldLevel) bool {
		return !v.containsBadWord(strings.ToLower(fl.Field().String()))
	})
}
