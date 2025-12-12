package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	v := New()
	assert.NotNil(t, v)
	assert.NotNil(t, v.badWords)
}

func TestValidateWord(t *testing.T) {
	v := New()

	tests := []struct {
		name    string
		word    string
		wantErr error
	}{
		{
			name:    "valid word",
			word:    "hello",
			wantErr: nil,
		},
		{
			name:    "valid word with spaces",
			word:    "  hello  ",
			wantErr: nil,
		},
		{
			name:    "empty word",
			word:    "",
			wantErr: ErrEmptyWord,
		},
		{
			name:    "whitespace only",
			word:    "   ",
			wantErr: ErrEmptyWord,
		},
		{
			name:    "single character",
			word:    "a",
			wantErr: ErrTooShort,
		},
		{
			name:    "two characters - valid",
			word:    "ab",
			wantErr: nil,
		},
		{
			name:    "unicode word",
			word:    "привіт",
			wantErr: nil,
		},
		{
			name:    "single unicode character",
			word:    "я",
			wantErr: ErrTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateWord(tt.word)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateLanguage(t *testing.T) {
	v := New()

	tests := []struct {
		name    string
		lang    string
		wantErr bool
	}{
		{
			name:    "english supported",
			lang:    "en",
			wantErr: false,
		},
		{
			name:    "ukrainian supported",
			lang:    "ua",
			wantErr: false,
		},
		{
			name:    "unsupported language",
			lang:    "fr",
			wantErr: true,
		},
		{
			name:    "empty language",
			lang:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateLanguage(tt.lang)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegisterCustomValidators(t *testing.T) {
	v := New()
	validate := validator.New()

	v.RegisterCustomValidators(validate)

	type TestStruct struct {
		Word string `validate:"nobadwords"`
	}

	tests := []struct {
		name    string
		word    string
		wantErr bool
	}{
		{
			name:    "valid word passes",
			word:    "hello",
			wantErr: false,
		},
		{
			name:    "empty word passes nobadwords",
			word:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := TestStruct{Word: tt.word}
			err := validate.Struct(ts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestContainsBadWord(t *testing.T) {
	v := New()

	assert.False(t, v.containsBadWord("hello"))
	assert.False(t, v.containsBadWord("world"))
}

func TestSupportedLanguages(t *testing.T) {
	assert.True(t, SupportedLanguages["en"])
	assert.True(t, SupportedLanguages["ua"])
	assert.False(t, SupportedLanguages["fr"])
	assert.False(t, SupportedLanguages[""])
}
