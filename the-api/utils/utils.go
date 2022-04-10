package utils

import (
	"strings"
	"time"
)

// Predicate that checks if the given trimmed string is empty
//  string string to check
//  bool true if string is nil or empty string, false otherwise
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// returns a random boolean value based on the current time
func RandomBool() bool {
	return time.Now().UnixNano()%2 == 0
}
