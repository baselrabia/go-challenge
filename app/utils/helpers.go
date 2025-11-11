package utils

import "fmt"

// parseIntParam parses a string parameter to int, returning defaultValue if parsing fails
func ParseIntParam(param string, defaultValue int) int {
	if param == "" {
		return defaultValue
	}
	var value int
	if _, err := fmt.Sscanf(param, "%d", &value); err != nil {
		return defaultValue
	}
	return value
}
