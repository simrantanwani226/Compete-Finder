package matcher

import (
	"strings"
)

var stopWords = map[string]bool{
	"the": true,
	"a":   true, "an": true, "and": true, "or": true, "but": true, "in": true, "on": true, "at": true, "to": true, "for": true, "of": true, "with": true, "by": true, "is": true, "it": true, "that": true, "this": true, "as": true, "are": true, "was": true, "be": true, "has": true,
	"have": true, "from": true, "we": true, "our": true, "their": true, "they": true, "you": true,
}

func Tokenize(s string) []string {
	s = strings.ToLower(s)
	r := strings.NewReplacer(",", " ", ".", " ", "?", " ", "/", " ", "!", " ")
	s = r.Replace(s)

	words := strings.Fields(s)
	result := make([]string, 0, len(words))
	for _, w := range words {
		if !stopWords[w] {
			result = append(result, w)
		}

	}
	return result
}
