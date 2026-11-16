package codec

import "testing"

func TestIsPrintable(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{name: "empty", data: nil, want: true},
		{name: "ascii text", data: []byte("secret-value"), want: true},
		{name: "tab and newline are retained", data: []byte("a\tb\nc\r"), want: true},
		{name: "accented utf8", data: []byte("café"), want: true},
		{name: "cjk utf8", data: []byte("日本語"), want: true},
		{name: "nul is rejected", data: []byte{'a', 0, 'b'}, want: false},
		{name: "bell is rejected", data: []byte{'a', 7, 'b'}, want: false},
		{name: "invalid utf8 is rejected", data: []byte{0xff, 0xfe}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPrintable(tt.data); got != tt.want {
				t.Fatalf("isPrintable(%q) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}
