package validation

func loadBadWords() map[string]struct{} {
	words := []string{
		// "word1",
		// "word2",
	}

	result := make(map[string]struct{}, len(words))
	for _, w := range words {
		result[w] = struct{}{}
	}
	return result
}
