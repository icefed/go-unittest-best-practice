package utils

import (
	"testing"
)

func TestFormatSize(t *testing.T) {
	cases := []struct {
		name     string
		size     int64
		expected string
	}{
		{
			name:     "zero",
			size:     0,
			expected: "0 B",
		},
		{
			name:     "1 byte",
			size:     1,
			expected: "1 B",
		},
		{
			name:     "64 bytes",
			size:     64,
			expected: "64 B",
		},
		{
			name:     "1024 bytes",
			size:     1024,
			expected: "1.0 KB",
		},
		{
			name:     "1034 bytes",
			size:     1034,
			expected: "1.0 KB",
		},
		{
			name:     "1934 bytes",
			size:     1934,
			expected: "1.9 KB",
		},
		{
			name:     "1000000 bytes",
			size:     1000000,
			expected: "976.6 KB",
		},
		{
			name:     "2000000 bytes",
			size:     2000000,
			expected: "1.9 MB",
		},
		{
			name:     "1000000000 bytes",
			size:     1000000000,
			expected: "953.7 MB",
		},
		{
			name:     "1000000000000 bytes",
			size:     1000000000000,
			expected: "931.3 GB",
		},
		{
			name:     "1000000000000000 bytes",
			size:     1000000000000000,
			expected: "909.5 TB",
		},
		{
			name:     "1000000000000000000 bytes",
			size:     1000000000000000000,
			expected: "888.2 PB",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			formatted := FormatSize(tc.size)
			if formatted != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, formatted)
			}
		})
	}
}

func TestFormatSize_subtests(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		formatted := FormatSize(0)
		if formatted != "0 B" {
			t.Errorf("expected \"0 B\", got \"%s\"", formatted)
		}
	})
	t.Run("100000", func(t *testing.T) {
		formatted := FormatSize(100000)
		if formatted != "97.7 KB" {
			t.Errorf("expected \"97.7 KB\", got \"%s\"", formatted)
		}
	})
	t.Run("-100", func(t *testing.T) {
		t.Skip()
	})
}

func BenchmarkFormatSize(b *testing.B) {
	// for i := 0; i < b.N; i++ {
	// 	FormatSize(100000000000)
	// }
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			FormatSize(100000000000)
		}
	})
}