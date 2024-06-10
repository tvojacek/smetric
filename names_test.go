package smetric

import (
	"errors"
	"strings"
	"testing"
)

func TestSnake(t *testing.T) {

	tests := []struct {
		name      string
		want      string
		input     string
		wantError bool
	}{
		{
			name:  "empty",
			input: "",
			want:  "",
		},
		{
			name:  "lower",
			input: "lower",
			want:  "lower",
		},
		{
			name:  "UPPER",
			input: "UPPER",
			want:  "upper",
		},
		{
			name:  "UPPER_SNAKE",
			input: "UPPER_SNAKE",
			want:  "upper_snake",
		},
		{
			name:  "Foo2Bar",
			input: "Foo2Bar",
			want:  "foo2_bar",
		},
		{
			name:  "camelCase",
			input: "camelCase",
			want:  "camel_case",
		},
		{
			name:  "camel_Snake",
			input: "camel_Snake",
			want:  "camel_snake",
		},
		{
			name:  "double__snake99",
			input: "double__Snake09",
			want:  "double__snake09",
		},
		{
			name:  "__prefix",
			input: "__prefix",
			want:  "__prefix",
		},
		{
			name:  "__prefix",
			input: "__prefix",
			want:  "__prefix",
		},
		{
			name:      "českýJazyk",
			input:     "českýJazyk",
			wantError: true,
			want:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSnakeMetricName(tt.input)
			if tt.wantError {
				if err == nil {
					t.Fatal("want error got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
			if tt.want != got {
				t.Fatalf("want: %s got: %s", tt.want, got)
			}
		})
	}
}

func BenchmarkSnake(b *testing.B) {
	input := "strcase_SnakeCase"

	b.Run("GetSnakeMetricName", func(b *testing.B) {
		b.ReportAllocs()
		for _ = range b.N {
			GetSnakeMetricName(input)
		}
	})
	b.Run("toLower", func(b *testing.B) {
		b.ReportAllocs()
		for _ = range b.N {
			strings.ToLower(input)
		}
	})
}

func BenchmarkSnakeLonger(b *testing.B) {
	input := "prometheusNotificationsTotalSeconds"
	b.Run("GetSnakeMetricName", func(b *testing.B) {
		b.ReportAllocs()
		for _ = range b.N {
			GetSnakeMetricName(input)
		}
	})
	b.Run("toLower", func(b *testing.B) {
		b.ReportAllocs()
		for _ = range b.N {
			strings.ToLower(input)
		}
	})
}

func TestError(t *testing.T) {
	err := InvalidCharacterError
	if !errors.Is(err, InvalidCharacterError) {
		t.Fatalf("err should be InvalidCharacterError")
	}

}
