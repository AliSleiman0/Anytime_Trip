package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// ValidateRequired returns an error message if the value is empty/whitespace.
func ValidateRequired(field, value string) string {
	if strings.TrimSpace(value) == "" {
		return field + " is required"
	}
	return ""
}

// ValidateLengthBetween ensures a string length is within the provided bounds (inclusive).
func ValidateLengthBetween(field, value string, min, max int) string {
	l := len(strings.TrimSpace(value))
	if l < min || l > max {
		return field + " must be between " + intToString(min) + " and " + intToString(max) + " characters"
	}
	return ""
}

// ValidateEmail checks basic email format.
func ValidateEmail(field, value string) string {
	if strings.TrimSpace(value) == "" {
		return field + " is required"
	}
	// Simple RFC 5322-ish email regex
	re := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
	if !re.MatchString(value) {
		return field + " is not a valid email"
	}
	return ""
}

// ValidatePhone checks a simple international phone pattern (digits, optional leading +, 7-15 digits).
func ValidatePhone(field, value string) string {
	if strings.TrimSpace(value) == "" {
		return field + " is required"
	}
	re := regexp.MustCompile(`^\+?[0-9]{7,15}$`)
	if !re.MatchString(value) {
		return field + " must be a valid phone number"
	}
	return ""
}

// ValidateEnum ensures the value is one of the allowed options.
func ValidateEnum(field, value string, allowed []string) string {
	for _, v := range allowed {
		if strings.EqualFold(value, v) {
			return ""
		}
	}
	return field + " must be one of: " + strings.Join(allowed, ", ")
}

// ValidateFloatRange ensures a float value lies within min/max (inclusive).
func ValidateFloatRange(field string, value, min, max float64) string {
	if value < min || value > max {
		return field + " must be between " + floatToString(min) + " and " + floatToString(max)
	}
	return ""
}

// ValidatePositiveFloat ensures value is >= 0.
func ValidatePositiveFloat(field string, value float64) string {
	if value < 0 {
		return field + " must be zero or positive"
	}
	return ""
}

// intToString avoids fmt import for small helpers.
func intToString(v int) string {
	return strconv.Itoa(v)
}

// floatToString avoids fmt import for small helpers.
func floatToString(v float64) string {
	// Trim trailing zeros via strconv
	return strings.TrimRight(strings.TrimRight(strconv.FormatFloat(v, 'f', -1, 64), "0"), ".")
}
