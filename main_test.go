package main

import (
	"testing"
)

func TestIsValidCEP(t *testing.T) {
	var tests = []struct {
		cep    string
		wanted bool
	}{
		{"12345-678", true},
		{"12345678", true},
		{"1234-567", false},
		{"12345-6789", false},
		{"abcd-efgh", false},
	}

	for _, tt := range tests {
		t.Run(tt.cep, func(t *testing.T) {
			got := isValidCEP(tt.cep)
			if got != tt.wanted {
				t.Errorf("isValidCEP(%s) = %v, want %v", tt.cep, got, tt.wanted)
			}
		})
	}
}
