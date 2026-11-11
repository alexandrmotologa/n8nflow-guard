package secrets

import (
	"math"
	"regexp"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// CalculateShannonEntropy computes the Shannon entropy of a given string in bits per symbol.
func CalculateShannonEntropy(input string) float64 {
	if len(input) == 0 {
		return 0.0
	}

	charCounts := make(map[rune]int)
	for _, r := range input {
		charCounts[r]++
	}

	total := float64(len(input))
	entropy := 0.0

	for _, count := range charCounts {
		p := float64(count) / total
		entropy -= p * math.Log2(p)
	}

	return entropy
}

// IsHighEntropySecret determines whether a candidate string exhibits characteristics of a generated secret.
// It excludes UUIDs, simple URLs, or strings shorter than 16 characters.
func IsHighEntropySecret(input string) bool {
	if len(input) < 16 || len(input) > 256 {
		return false
	}

	// Exclude UUIDs
	if uuidRegex.MatchString(input) {
		return false
	}

	entropy := CalculateShannonEntropy(input)
	// Base64 or mixed alphanumeric keys typically exceed 4.3 bits of entropy
	return entropy >= 4.3
}
