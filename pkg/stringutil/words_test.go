package stringutil

import (
	"math/rand"
	"testing"
	"time"
)

func TestRandomWord(t *testing.T) {
	word := RandomWord()
	if word == "" {
		t.Errorf("expected a non-empty word, got an empty string")
	}

	found := false
	for _, w := range words {
		if w == word {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected word to be from the predefined list, got %s", word)
	}
}

func TestRandomWords(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	tests := []struct {
		name     string
		count    int
		expected int
	}{
		{"Zero words", 0, 0},
		{"One word", 1, 1},
		{"Multiple words", 5, 5},
		{"Negative count", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words := RandomWords(tt.count)
			if len(words) != tt.expected {
				t.Errorf("expected %d words, got %d", tt.expected, len(words))
			}

			for _, word := range words {
				if word == "" {
					t.Errorf("expected a non-empty word, got an empty string")
				}

				found := false
				for _, w := range words {
					if w == word {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("expected word to be from the predefined list, got %s", word)
				}
			}
		})
	}
}
