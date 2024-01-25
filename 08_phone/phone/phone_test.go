package phone

import (
	"testing"
)

var (
	testCases = [...]struct {
		input string
		want  string
	}{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"123-456-7890", "1234567890"},
		{"1234567892", "1234567892"},
		{"(123)456-7892", "1234567892"},
	}
)

type normalizeFunc func(string) string

func TestNormalize(t *testing.T) {
	testNormalizeGeneral(t, func(s string) string { return Normalize(s) })
}

func TestNormalizeRegex(t *testing.T) {
	testNormalizeGeneral(t, func(s string) string { return NormalizeRegex(s) })
}

func testNormalizeGeneral(t *testing.T, fn normalizeFunc) {
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := fn(tc.input)
			if actual != tc.want {
				t.Errorf("got %s, want %s", actual, tc.want)
			}
		})
	}
}
