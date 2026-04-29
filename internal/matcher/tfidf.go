package matcher

import (
	"math"
	"sort"

	"github.com/simrantanwani226/compete-finder/internal/provider"
)

type Result struct {
	Startup provider.Startup
	Score   float64
}

func tf(tokens []string) map[string]float64 {
	counts := make(map[string]float64)
	for _, w := range tokens {
		counts[w] += 1
	}
	for word, count := range counts {
		counts[word] = count / float64(len(tokens))
	}
	return counts
}

func idf(docs [][]string) map[string]float64 {

	counts := make(map[string]float64)
	for _, w := range docs {
		seen := make(map[string]bool)
		for _, c := range w {
			if !seen[c] {
				seen[c] = true
				counts[c] += 1
			}
		}
	}
	for word, count := range counts {
		counts[word] = math.Log(float64(len(docs)) / count)
	}
	return counts
}

func tfidfVec(tokens []string, idfScores map[string]float64) map[string]float64 {
	tfscores := tf(tokens)
	for word, tfscore := range tfscores {
		tfscores[word] = tfscore * idfScores[word]

	}
	return tfscores
}
func cosineSim(a, b map[string]float64) float64 {
	var dot float64
	var magniA float64
	var magnib float64
	for word, scoreA := range a {
		dot += scoreA * b[word]
		magniA += scoreA * scoreA
	}
	for _, scoreB := range b {

		magnib += scoreB * scoreB
	}
	magniA = math.Sqrt(magniA)
	magnib = math.Sqrt(magnib)

	if magniA == 0 || magnib == 0 {

		return 0
	}
	return dot / (magniA * magnib)
}
func Match(query string, startups []provider.Startup, limit int) []Result {
	corpus := make([][]string, 0, len(startups))
	for _, s := range startups {
		tokens := Tokenize(s.Description)
		corpus = append(corpus, tokens)
	}
	querytoken := Tokenize(query)
	corpus = append(corpus, querytoken)
	idfscore := idf(corpus)
	vec := tfidfVec(querytoken, idfscore)
	results := make([]Result, 0, len(startups))
	for _, s := range startups {
		startuptoken := Tokenize(s.Description)
		vector := tfidfVec(startuptoken, idfscore)
		score := cosineSim(vec, vector)
		results = append(results, Result{Startup: s, Score: score})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	if limit < len(results) {
		results = results[:limit]
	}
	return results

}
