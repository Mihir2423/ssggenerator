package markdown

import "testing"

func TestToHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "h1 heading",
			input:    "# Hello",
			expected: "<h1>Hello</h1>",
		},
		{
			name:     "h2 heading",
			input:    "## World",
			expected: "<h2>World</h2>",
		},
		{
			name:     "paragraph",
			input:    "Hello world",
			expected: "<p>Hello world</p>",
		},
		{
			name:     "multiple lines",
			input:    "# Title\n\nHello\nWorld",
			expected: "<h1>Title</h1><p>Hello</p><p>World</p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := ToHTML([]byte(tt.input))
			if output != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, output)
			}
		})
	}
}
