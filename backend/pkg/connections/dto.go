package connections

type StartGameRequest struct {
	Language   string `json:"language"`
	SecretWord string `json:"secret_word"`
}

type StartGameResponse struct {
	CalculationTime float64 `json:"calculation_time"`
	HintWord        string  `json:"hint_word,omitempty"`
}

type GuessRequest struct {
	SecretWord string `json:"secret_word"`
	Guess      string `json:"guess"`
	Language   string `json:"language"`
}

type GuessResponse struct {
	Distance int `json:"distance"`
}

type HealthResponse struct {
	Status          string   `json:"status"`
	LoadedLanguages []string `json:"loaded_languages"`
}
