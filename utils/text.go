package utils

import (
        "fmt"
        "strings"
        "time"
)

// RemoveDuplicates removes duplicate strings from a slice
func RemoveDuplicates(slice []string) []string {
        keys := make(map[string]bool)
        var result []string

        for _, item := range slice {
                item = strings.TrimSpace(item)
                if item != "" && !keys[item] {
                        keys[item] = true
                        result = append(result, item)
                }
        }

        return result
}

// CleanText removes extra whitespace and normalizes text
func CleanText(text string) string {
        // Replace multiple whitespace with single space
        text = strings.Join(strings.Fields(text), " ")
        return strings.TrimSpace(text)
}

// ExtractNumbers extracts numbers from text string
func ExtractNumbers(text string) []string {
        var numbers []string
        words := strings.Fields(text)

        for _, word := range words {
                // Simple number detection
                if len(word) > 0 && (word[0] >= '0' && word[0] <= '9') {
                        numbers = append(numbers, word)
                }
        }

        return numbers
}

// GenerateTimestamp generates a timestamp for file naming
func GenerateTimestamp() int64 {
        return time.Now().Unix()
}

// SplitIntoSentences splits text into sentences
func SplitIntoSentences(text string) []string {
        // Simple sentence splitting by periods, exclamation marks, and question marks
        text = strings.ReplaceAll(text, "!", ".")
        text = strings.ReplaceAll(text, "?", ".")
        
        sentences := strings.Split(text, ".")
        var result []string
        
        for _, sentence := range sentences {
                sentence = strings.TrimSpace(sentence)
                if len(sentence) > 5 { // Filter out very short fragments
                        result = append(result, sentence)
                }
        }
        
        return result
}

// CalculateWordCount counts words in text
func CalculateWordCount(text string) int {
        return len(strings.Fields(text))
}

// ExtractEmailDomain extracts domain from email address
func ExtractEmailDomain(email string) string {
        parts := strings.Split(email, "@")
        if len(parts) == 2 {
                return parts[1]
        }
        return ""
}

// TruncateText truncates text to specified length with ellipsis
func TruncateText(text string, maxLength int) string {
        if len(text) <= maxLength {
                return text
        }
        
        if maxLength <= 3 {
                return text[:maxLength]
        }
        
        return text[:maxLength-3] + "..."
}

// IsValidEmailFormat checks if string is a valid email format
func IsValidEmailFormat(email string) bool {
        // Simple email validation
        if !strings.Contains(email, "@") {
                return false
        }
        
        parts := strings.Split(email, "@")
        if len(parts) != 2 {
                return false
        }
        
        if len(parts[0]) == 0 || len(parts[1]) == 0 {
                return false
        }
        
        if !strings.Contains(parts[1], ".") {
                return false
        }
        
        return true
}

// ContainsAny checks if text contains any of the provided substrings
func ContainsAny(text string, substrings []string) bool {
        textLower := strings.ToLower(text)
        for _, substring := range substrings {
                if strings.Contains(textLower, strings.ToLower(substring)) {
                        return true
                }
        }
        return false
}

// ExtractYearsFromText extracts 4-digit years from text
func ExtractYearsFromText(text string) []int {
        var years []int
        words := strings.Fields(text)
        
        for _, word := range words {
                // Remove common punctuation
                word = strings.Trim(word, "().,;:")
                
                // Check if it's a 4-digit year (1900-2030 range)
                if len(word) == 4 {
                        var year int
                        if _, err := fmt.Sscanf(word, "%d", &year); err == nil {
                                if year >= 1900 && year <= 2030 {
                                        years = append(years, year)
                                }
                        }
                }
        }
        
        return years
}
