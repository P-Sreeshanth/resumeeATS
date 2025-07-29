package services

import (
	"math"
	"sort"
	"strings"
	"unicode"
)

// NLPService provides natural language processing capabilities
type NLPService struct {
	stopWords map[string]bool
}

// NewNLPService creates a new NLP service instance
func NewNLPService() *NLPService {
	stopWords := map[string]bool{
		"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
		"be": true, "by": true, "for": true, "from": true, "has": true, "he": true,
		"in": true, "is": true, "it": true, "its": true, "of": true, "on": true,
		"that": true, "the": true, "to": true, "was": true, "were": true, "will": true,
		"with": true, "you": true, "your": true, "have": true, "had": true, "this": true,
		"they": true, "we": true, "our": true, "us": true, "can": true, "could": true,
		"would": true, "should": true, "may": true, "might": true, "must": true,
	}

	return &NLPService{
		stopWords: stopWords,
	}
}

// Document represents a document for TF-IDF analysis
type Document struct {
	Text  string
	Terms map[string]int
}

// TFIDFResult represents a term with its TF-IDF score
type TFIDFResult struct {
	Term  string
	Score float64
}

// Tokenize splits text into tokens and removes stop words
func (nlp *NLPService) Tokenize(text string) []string {
	// Convert to lowercase and split by non-letter characters
	words := strings.FieldsFunc(strings.ToLower(text), func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})

	var tokens []string
	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) > 2 && !nlp.stopWords[word] {
			tokens = append(tokens, word)
		}
	}

	return tokens
}

// CalculateTFIDF calculates TF-IDF scores for terms in a document collection
func (nlp *NLPService) CalculateTFIDF(documents []string) map[string]float64 {
	docs := make([]Document, len(documents))
	allTerms := make(map[string]int)

	// Process each document
	for i, docText := range documents {
		tokens := nlp.Tokenize(docText)
		termFreq := make(map[string]int)

		for _, token := range tokens {
			termFreq[token]++
			allTerms[token]++
		}

		docs[i] = Document{
			Text:  docText,
			Terms: termFreq,
		}
	}

	// Calculate TF-IDF scores
	tfidfScores := make(map[string]float64)
	numDocs := float64(len(documents))

	for term := range allTerms {
		// Calculate document frequency (DF)
		df := 0
		for _, doc := range docs {
			if doc.Terms[term] > 0 {
				df++
			}
		}

		// Calculate IDF
		idf := math.Log(numDocs / float64(df))

		// Calculate TF-IDF for the term across all documents
		totalTFIDF := 0.0
		for _, doc := range docs {
			tf := float64(doc.Terms[term])
			if tf > 0 {
				// Normalize TF by document length
				docLength := 0
				for _, freq := range doc.Terms {
					docLength += freq
				}
				normalizedTF := tf / float64(docLength)
				totalTFIDF += normalizedTF * idf
			}
		}

		tfidfScores[term] = totalTFIDF
	}

	return tfidfScores
}

// ExtractKeywords extracts top keywords from text using TF-IDF
func (nlp *NLPService) ExtractKeywords(text string, topK int) []string {
	// For single document, we'll use term frequency as a simple approach
	tokens := nlp.Tokenize(text)
	termFreq := make(map[string]int)

	for _, token := range tokens {
		termFreq[token]++
	}

	// Sort by frequency
	type termScore struct {
		term  string
		score int
	}

	var scores []termScore
	for term, freq := range termFreq {
		scores = append(scores, termScore{term, freq})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Return top K terms
	var keywords []string
	for i := 0; i < topK && i < len(scores); i++ {
		keywords = append(keywords, scores[i].term)
	}

	return keywords
}

// CalculateCosineSimilarity calculates cosine similarity between two texts
func (nlp *NLPService) CalculateCosineSimilarity(text1, text2 string) float64 {
	tokens1 := nlp.Tokenize(text1)
	tokens2 := nlp.Tokenize(text2)

	// Create term frequency vectors
	allTerms := make(map[string]bool)
	freq1 := make(map[string]int)
	freq2 := make(map[string]int)

	for _, token := range tokens1 {
		freq1[token]++
		allTerms[token] = true
	}

	for _, token := range tokens2 {
		freq2[token]++
		allTerms[token] = true
	}

	// Calculate cosine similarity
	var dotProduct, norm1, norm2 float64

	for term := range allTerms {
		f1 := float64(freq1[term])
		f2 := float64(freq2[term])

		dotProduct += f1 * f2
		norm1 += f1 * f1
		norm2 += f2 * f2
	}

	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

// CalculateSkillMatch calculates skill matching percentage
func (nlp *NLPService) CalculateSkillMatch(resumeSkills, jobSkills []string) (float64, []string, []string) {
	resumeSet := make(map[string]bool)
	for _, skill := range resumeSkills {
		resumeSet[strings.ToLower(skill)] = true
	}

	var matched []string
	var missing []string

	for _, jobSkill := range jobSkills {
		jobSkillLower := strings.ToLower(jobSkill)
		if resumeSet[jobSkillLower] {
			matched = append(matched, jobSkill)
		} else {
			// Check for partial matches (fuzzy matching)
			found := false
			for resumeSkill := range resumeSet {
				if nlp.calculateStringSimilarity(resumeSkill, jobSkillLower) > 0.8 {
					matched = append(matched, jobSkill)
					found = true
					break
				}
			}
			if !found {
				missing = append(missing, jobSkill)
			}
		}
	}

	matchPercentage := 0.0
	if len(jobSkills) > 0 {
		matchPercentage = float64(len(matched)) / float64(len(jobSkills)) * 100
	}

	return matchPercentage, matched, missing
}

// calculateStringSimilarity calculates string similarity using Levenshtein distance
func (nlp *NLPService) calculateStringSimilarity(s1, s2 string) float64 {
	if len(s1) == 0 {
		return float64(len(s2))
	}
	if len(s2) == 0 {
		return float64(len(s1))
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}

	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	distance := matrix[len(s1)][len(s2)]
	maxLen := max(len(s1), len(s2))
	
	if maxLen == 0 {
		return 1.0
	}
	
	return 1.0 - float64(distance)/float64(maxLen)
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
