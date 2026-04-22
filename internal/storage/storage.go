package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// SanitizeTitle converts a title to a valid filename
func SanitizeTitle(title string) string {
	// Convert to lowercase
	title = strings.ToLower(title)

	// Replace spaces with hyphens
	title = strings.ReplaceAll(title, " ", "-")

	// Remove special characters except hyphens and alphanumeric
	re := regexp.MustCompile(`[^a-z0-9-]`)
	title = re.ReplaceAllString(title, "")

	// Remove multiple consecutive hyphens
	re = regexp.MustCompile(`-+`)
	title = re.ReplaceAllString(title, "-")

	// Trim hyphens from start and end
	title = strings.Trim(title, "-")

	return title
}

// GenerateAutoTitle generates an auto-title from the note content
// Returns first 32 characters + "..." if longer than 32 chars
func GenerateAutoTitle(content string) string {
	// Get the first line
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return "untitled"
	}

	firstLine := strings.TrimSpace(lines[0])
	if firstLine == "" {
		return "untitled"
	}

	// Truncate to 32 characters
	if len(firstLine) > 32 {
		return firstLine[:32] + "..."
	}

	return firstLine
}

// GenerateFilename creates a filename with date prefix and sanitized title
func GenerateFilename(title string) string {
	date := time.Now().Format("2006-01-02")
	sanitizedTitle := SanitizeTitle(title)

	if sanitizedTitle == "" {
		sanitizedTitle = "untitled"
	}

	return fmt.Sprintf("%s-%s.md", date, sanitizedTitle)
}

// ResolveDuplicate checks if a file exists and appends a number if it does
// Returns the final filepath that doesn't conflict
func ResolveDuplicate(targetDir, filename string) (string, error) {
	filePath := filepath.Join(targetDir, filename)

	// If file doesn't exist, we're good
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return filePath, nil
	}

	// File exists, need to append number
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	// Extract date prefix and title
	parts := strings.SplitN(base, "-", 4) // YYYY-MM-DD-title
	if len(parts) < 3 {
		// Unexpected format, just append number to full base
		counter := 1
		for {
			newFilename := fmt.Sprintf("%s-%d%s", base, counter, ext)
			newFilepath := filepath.Join(targetDir, newFilename)
			if _, err := os.Stat(newFilepath); os.IsNotExist(err) {
				return newFilepath, nil
			}
			counter++
		}
	}

	// Build base without counter
	datePrefix := strings.Join(parts[:3], "-")
	titlePart := ""
	if len(parts) > 3 {
		titlePart = parts[3]
	}

	counter := 1
	for {
		var newFilename string
		if titlePart == "" {
			newFilename = fmt.Sprintf("%s-%d%s", datePrefix, counter, ext)
		} else {
			newFilename = fmt.Sprintf("%s-%s-%d%s", datePrefix, titlePart, counter, ext)
		}

		newFilepath := filepath.Join(targetDir, newFilename)
		if _, err := os.Stat(newFilepath); os.IsNotExist(err) {
			return newFilepath, nil
		}
		counter++
	}
}

// SaveNote saves a note to the target directory
func SaveNote(targetDir, content, title string) (string, error) {
	// Check if target directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return "", fmt.Errorf("target directory does not exist: %s", targetDir)
	}

	// Generate filename
	filename := GenerateFilename(title)

	// Resolve any duplicates
	filepath, err := ResolveDuplicate(targetDir, filename)
	if err != nil {
		return "", fmt.Errorf("failed to resolve duplicate filename: %w", err)
	}

	// Write the file
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filepath, nil
}
