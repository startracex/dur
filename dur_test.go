package dur

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	var testCases = []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "only_year_month_day",
			input:   "2y 3mon 5d",
			want:    2*yearDuration + 3*monthDuration + 5*dayDuration,
			wantErr: false,
		},
		{
			name:    "only_ms_us_ns",
			input:   "100ms 200us 300ns",
			want:    100*time.Millisecond + 200*time.Microsecond + 300*time.Nanosecond,
			wantErr: false,
		},
		{
			name:    "full_unit_names",
			input:   "1year 2months 3days 4hours 5minutes 6seconds",
			want:    1*yearDuration + 2*monthDuration + 3*dayDuration + 4*time.Hour + 5*time.Minute + 6*time.Second,
			wantErr: false,
		},
		{
			name:    "mixed_whitespace",
			input:   "1 h\t2 m\n3 s  ",
			want:    time.Hour + 2*time.Minute + 3*time.Second,
			wantErr: false,
		},
		{
			name:    "duplicate_unit_hour",
			input:   "1h 2h",
			want:    0,
			wantErr: true,
		},
		{
			name:    "duplicate_unit_month",
			input:   "1mon 2months",
			want:    0,
			wantErr: true,
		},
		{
			name:    "zero_number",
			input:   "0h",
			want:    0,
			wantErr: false,
		},
		{
			name:    "too_long_unit",
			input:   "1millisecondssssss",
			want:    0,
			wantErr: true,
		},
		{
			name:    "mixed_case_unit",
			input:   "1H 2M 3S 4MS",
			want:    time.Hour + 2*time.Minute + 3*time.Second + 4*time.Millisecond,
			wantErr: false,
		},
		{
			name:    "unit_variants",
			input:   "1yrs 2hrs 3mins 4secs",
			want:    1*yearDuration + 2*time.Hour + 3*time.Minute + 4*time.Second,
			wantErr: false,
		},
		{
			name:    "no_unit_after_digit",
			input:   "123",
			want:    0,
			wantErr: true,
		},
		{
			name:    "illegal_character",
			input:   "1h$2m",
			want:    0,
			wantErr: true,
		},
		{
			name:    "min_vs_mon",
			input:   "5min 3mon",
			want:    5*time.Minute + 3*monthDuration,
			wantErr: false,
		},
		{
			name:    "only_whitespace",
			input:   "   \t\n",
			want:    0,
			wantErr: true,
		},
		{
			name:    "large_number",
			input:   "999999h 999999m",
			want:    999999*time.Hour + 999999*time.Minute,
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dur, err := Parse(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("Parse(%q) err=%v, wantErr=%v", tc.input, err, tc.wantErr)
				return
			}
			if !tc.wantErr && dur != tc.want {
				t.Errorf("Parse(%q) = %v, want %v", tc.input, dur, tc.want)
			}
		})
	}
}

func TestParse_WithSign(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "single_negative",
			input:   "-1h",
			want:    -time.Hour,
			wantErr: false,
		},
		{
			name:    "mixed_signs",
			input:   "1h -30m +5s",
			want:    time.Hour - 30*time.Minute + 5*time.Second,
			wantErr: false,
		},
		{
			name:    "positive_sign",
			input:   "+2h +10m",
			want:    2*time.Hour + 10*time.Minute,
			wantErr: false,
		},
		{
			name:    "sign_between_charters",
			input:   "1h-30m",
			want:    time.Hour - 30*time.Minute,
			wantErr: false,
		},
		{
			name:    "sign_without_digit",
			input:   "1h -",
			want:    0,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dur, err := Parse(tc.input)
			hasErr := err != nil
			if hasErr != tc.wantErr {
				t.Errorf("Parse(%q) err=%v, wantErr=%v", tc.input, err, tc.wantErr)
				return
			}
			if dur != tc.want {
				t.Errorf("Parse(%q) = %v, want %v", tc.input, dur, tc.want)
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	stdInput := "1h 30m 45s 100ms 500us 100ns"

	b.Run("time.ParseDuration", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			_, _ = time.ParseDuration(stdInput)
		}
	})

	b.Run("custom_Parse_std_input", func(b *testing.B) {
		b.ResetTimer()
		for range b.N {
			_, _ = Parse(stdInput)
		}
	})
}
