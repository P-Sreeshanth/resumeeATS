package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IsValidResumeFile checks if the uploaded file is a valid resume format
func IsValidResumeFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExtensions := []string{".pdf", ".docx"}
	
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	
	return false
}

// ValidateFileSize checks if file size is within acceptable limits
func ValidateFileSize(fileSize int64) error {
	const maxSize = 10 * 1024 * 1024 // 10MB
	
	if fileSize > maxSize {
		return fmt.Errorf("file size too large: %d bytes (max: %d bytes)", fileSize, maxSize)
	}
	
	if fileSize == 0 {
		return fmt.Errorf("file is empty")
	}
	
	return nil
}

// CleanupFile removes a file from the filesystem
func CleanupFile(filename string) error {
	if filename == "" {
		return nil
	}
	
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to cleanup file %s: %v", filename, err)
	}
	
	return nil
}

// SanitizeFilename removes potentially dangerous characters from filename
func SanitizeFilename(filename string) string {
	// Replace dangerous characters with underscores
	dangerous := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", ".."}
	
	for _, char := range dangerous {
		filename = strings.ReplaceAll(filename, char, "_")
	}
	
	// Limit filename length
	if len(filename) > 255 {
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)
		filename = name[:255-len(ext)] + ext
	}
	
	return filename
}

// ValidateJobDescription validates job description content
func ValidateJobDescription(jobDesc string) error {
	jobDesc = strings.TrimSpace(jobDesc)
	
	if len(jobDesc) == 0 {
		return fmt.Errorf("job description cannot be empty")
	}
	
	if len(jobDesc) < 50 {
		return fmt.Errorf("job description too short (minimum 50 characters)")
	}
	
	if len(jobDesc) > 10000 {
		return fmt.Errorf("job description too long (maximum 10000 characters)")
	}
	
	return nil
}

// IsTextFile checks if a file contains readable text
func IsTextFile(content []byte) bool {
	// Simple heuristic: check if most bytes are printable ASCII or common UTF-8
	printableCount := 0
	totalCount := len(content)
	
	if totalCount == 0 {
		return false
	}
	
	for _, b := range content {
		if (b >= 32 && b <= 126) || b == 9 || b == 10 || b == 13 {
			printableCount++
		}
	}
	
	// If more than 80% of bytes are printable, consider it text
	ratio := float64(printableCount) / float64(totalCount)
	return ratio > 0.8
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	
	if !IsValidEmailFormat(email) {
		return fmt.Errorf("invalid email format")
	}
	
	return nil
}

// ValidatePhone validates phone number format
func ValidatePhone(phone string) error {
	phone = strings.TrimSpace(phone)
	
	if phone == "" {
		return fmt.Errorf("phone number cannot be empty")
	}
	
	// Remove common formatting characters
	cleanPhone := strings.ReplaceAll(phone, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "(", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, ")", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, ".", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "+", "")
	
	// Check if remaining characters are digits
	for _, char := range cleanPhone {
		if char < '0' || char > '9' {
			return fmt.Errorf("phone number contains invalid characters")
		}
	}
	
	// Check length (US phone numbers)
	if len(cleanPhone) < 10 || len(cleanPhone) > 15 {
		return fmt.Errorf("phone number length invalid (should be 10-15 digits)")
	}
	
	return nil
}

// EnsureDirectoryExists creates directory if it doesn't exist
func EnsureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
