package main

import (
	"testing"
)

func Test_sanitizeLine(t *testing.T) {
	tests := []struct {
		name     string
		line     []byte
		expected string
	}{
		{"leading_white_space", []byte("  test"), "test"},
		{"trailing_white_space", []byte("test  "), "test"},
		{"leading_trailing_white_space", []byte("  test  "), "test"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			line := sanitizeLine(test.line)

			if line != test.expected {
				t.Errorf("sanitizeLine() line = %v, want %v", line, test.expected)
			}
		})
	}
}
