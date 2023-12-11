package main

import (
	"reflect"
	"testing"
)

func TestMyersDiff(t *testing.T) {
	tests := []struct {
		name     string
		aLines   []string
		bLines   []string
		expected []EditAction
	}{
		{
			name:     "No changes",
			aLines:   []string{"line1", "line2", "line3"},
			bLines:   []string{"line1", "line2", "line3"},
			expected: []EditAction{Keep{"line1"}, Keep{"line2"}, Keep{"line3"}},
		},
		{
			name:     "Insertion",
			aLines:   []string{"line1", "line2"},
			bLines:   []string{"line1", "line2", "line3"},
			expected: []EditAction{Keep{"line1"}, Keep{"line2"}, Insert{"line3"}},
		},
		{
			name:     "Removal",
			aLines:   []string{"line1", "line2", "line3"},
			bLines:   []string{"line1", "line2"},
			expected: []EditAction{Keep{"line1"}, Keep{"line2"}, Remove{"line3"}},
		},
		{
			name:     "Replacement",
			aLines:   []string{"line1", "line2", "line3"},
			bLines:   []string{"line1", "newLine2", "line3"},
			expected: []EditAction{Keep{"line1"}, Remove{"line2"}, Insert{"newLine2"}, Keep{"line3"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := myersDiff(tt.aLines, tt.bLines)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("myersDiff() = %v, want %v", result, tt.expected)
			}
		})
	}
}
