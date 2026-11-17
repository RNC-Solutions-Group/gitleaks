package codec

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMayContainEncoding(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{
			name: "ordinary source text",
			data: "const ordinary_value = 'synthetic content for scanner timing only'\n",
			want: false,
		},
		{
			name: "ordinary long identifier without a likely base64 character",
			data: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			want: false,
		},
		{
			name: "base64",
			data: "c2VjcmV0LWtleS0xMjM0NTY3ODkw",
			want: true,
		},
		{
			name: "base64url",
			data: base64.RawURLEncoding.EncodeToString([]byte("secret-key-1234567890")),
			want: true,
		},
		{
			name: "percent",
			data: "secret%2Dkey%2D1234567890",
			want: true,
		},
		{
			name: "unicode escape",
			data: `config = "\u0073\u0065\u0063\u0072\u0065\u0074"`,
			want: true,
		},
		{
			name: "unicode code point",
			data: "U+0073 U+0065 U+0063 U+0072 U+0065 U+0074",
			want: true,
		},
		{
			name: "hex",
			data: "7365637265742d6b65792d31323334353637383930",
			want: true,
		},
		{
			name: "empty",
			data: "",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, mayContainEncoding(tt.data))
		})
	}
}

func TestFindEncodingMatchesStillFindsEncodedValues(t *testing.T) {
	// This is the behavior the prefilter must never change. Each case is a
	// value that the real decoder turns into a different string.
	tests := []struct {
		name string
		data string
	}{
		{name: "base64", data: "c2VjcmV0LWtleS0xMjM0NTY3ODkw"},
		{name: "base64url", data: base64.RawURLEncoding.EncodeToString([]byte("secret-key-1234567890"))},
		{name: "percent", data: "secret%2Dkey%2D1234567890"},
		{name: "unicode", data: `\u0073\u0065\u0063\u0072\u0065\u0074`},
		{name: "hex", data: "7365637265742d6b65792d31323334353637383930"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := findEncodingMatches(tt.data)
			require.NotEmpty(t, matches)
			decoded := NewDecoder()
			got, segments := decoded.Decode(tt.data, nil)
			assert.NotEqual(t, tt.data, got)
			assert.NotEmpty(t, segments)
		})
	}
}

func TestFindEncodingMatchesRejectsOrdinaryText(t *testing.T) {
	data := "const ordinary_value = 'synthetic content for scanner timing only'\n"
	assert.Empty(t, findEncodingMatches(data))
}
