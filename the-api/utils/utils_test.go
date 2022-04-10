package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetunTrueOnBlankStrings(t *testing.T) {
	actual := IsEmpty("") || IsEmpty("  ")
	assert.Equal(t, actual, true /*expected*/)
}

func TestFalseOnNonBlankStrings(t *testing.T) {
	actual := IsEmpty("a") || IsEmpty(" a b c ")
	assert.Equal(t, actual, false /*expected*/)
}
