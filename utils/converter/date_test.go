package converter_test

import (
	"senkou-catalyst-be/utils/converter"
	"testing"
	"time"
)

func TestParseDateFunction(t *testing.T) {

	defaultDate := time.Date(2025, 9, 6, 0, 0, 0, 0, time.UTC)
	validDate := time.Date(2023, 12, 25, 10, 30, 45, 0, time.UTC)

	t.Run("Should return default date when input date is nil", func(t *testing.T) {
		var nilDate *time.Time = nil

		result := converter.ParseDate(nilDate, defaultDate)

		if !result.Equal(defaultDate) {
			t.Errorf("Expected default date %v, got %v", defaultDate, result)
		}

		t.Logf("Correctly returned default date for nil input: %v", result.Format("2006-01-02"))
	})

	t.Run("Should return default date when input date is zero time", func(t *testing.T) {
		zeroTime := time.Time{}

		result := converter.ParseDate(&zeroTime, defaultDate)

		if !result.Equal(defaultDate) {
			t.Errorf("Expected default date %v, got %v", defaultDate, result)
		}

		t.Logf("Correctly returned default date for zero time: %v", result.Format("2006-01-02"))
	})

	t.Run("Should parse valid date correctly", func(t *testing.T) {
		inputDate := validDate

		result := converter.ParseDate(&inputDate, defaultDate)

		expectedDate := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)

		if !result.Equal(expectedDate) {
			t.Errorf("Expected parsed date %v, got %v", expectedDate, result)
		}

		if result.Year() != 2023 {
			t.Errorf("Expected year 2023, got %d", result.Year())
		}
		if result.Month() != 12 {
			t.Errorf("Expected month 12, got %d", result.Month())
		}
		if result.Day() != 25 {
			t.Errorf("Expected day 25, got %d", result.Day())
		}

		if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
			t.Errorf("Expected time to be reset to 00:00:00, got %02d:%02d:%02d",
				result.Hour(), result.Minute(), result.Second())
		}

		t.Logf("Successfully parsed date: %v", result.Format("2006-01-02 15:04:05"))
	})

	t.Run("Should handle date with different time zones", func(t *testing.T) {
		jakarta, _ := time.LoadLocation("Asia/Jakarta")
		jakartaDate := time.Date(2023, 12, 25, 15, 30, 0, 0, jakarta)

		result := converter.ParseDate(&jakartaDate, defaultDate)

		expectedDate := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)

		if !result.Equal(expectedDate) {
			t.Errorf("Expected parsed date %v, got %v", expectedDate, result)
		}

		t.Logf("Correctly handled timezone conversion: %v -> %v",
			jakartaDate.Format("2006-01-02 15:04:05 MST"),
			result.Format("2006-01-02 15:04:05 MST"))
	})

	t.Run("Should handle leap year dates", func(t *testing.T) {
		leapDate := time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC)

		result := converter.ParseDate(&leapDate, defaultDate)

		expectedDate := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)

		if !result.Equal(expectedDate) {
			t.Errorf("Expected leap year date %v, got %v", expectedDate, result)
		}

		if result.Month() != 2 || result.Day() != 29 {
			t.Errorf("Expected February 29, got %s %d", result.Month(), result.Day())
		}

		t.Logf("Correctly handled leap year date: %v", result.Format("2006-01-02"))
	})

	t.Run("Should handle edge case dates", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    time.Time
			expected time.Time
		}{
			{
				name:     "New Year",
				input:    time.Date(2024, 1, 1, 23, 59, 59, 0, time.UTC),
				expected: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				name:     "Year End",
				input:    time.Date(2023, 12, 31, 12, 30, 45, 0, time.UTC),
				expected: time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			},
			{
				name:     "Mid Year",
				input:    time.Date(2023, 6, 15, 8, 45, 30, 0, time.UTC),
				expected: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := converter.ParseDate(&tc.input, defaultDate)

				if !result.Equal(tc.expected) {
					t.Errorf("%s: Expected %v, got %v", tc.name, tc.expected, result)
				}

				t.Logf("%s: %v -> %v", tc.name,
					tc.input.Format("2006-01-02 15:04:05"),
					result.Format("2006-01-02 15:04:05"))
			})
		}
	})

	t.Run("Should preserve UTC timezone in result", func(t *testing.T) {
		inputDate := time.Date(2023, 12, 25, 15, 30, 0, 0, time.UTC)

		result := converter.ParseDate(&inputDate, defaultDate)

		if result.Location() != time.UTC {
			t.Errorf("Expected UTC timezone, got %v", result.Location())
		}

		t.Logf("Result timezone is UTC: %v", result.Location())
	})

	t.Run("Should handle very old dates", func(t *testing.T) {
		oldDate := time.Date(1900, 1, 1, 12, 0, 0, 0, time.UTC)

		result := converter.ParseDate(&oldDate, defaultDate)

		expectedDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

		if !result.Equal(expectedDate) {
			t.Errorf("Expected old date %v, got %v", expectedDate, result)
		}

		t.Logf("Correctly handled old date: %v", result.Format("2006-01-02"))
	})

	t.Run("Should handle future dates", func(t *testing.T) {
		futureDate := time.Date(2050, 12, 31, 23, 59, 59, 0, time.UTC)

		result := converter.ParseDate(&futureDate, defaultDate)

		expectedDate := time.Date(2050, 12, 31, 0, 0, 0, 0, time.UTC)

		if !result.Equal(expectedDate) {
			t.Errorf("Expected future date %v, got %v", expectedDate, result)
		}

		t.Logf("Correctly handled future date: %v", result.Format("2006-01-02"))
	})
}

func TestParseDatePerformance(t *testing.T) {

	t.Run("Should perform well with multiple calls", func(t *testing.T) {
		defaultDate := time.Date(2025, 9, 6, 0, 0, 0, 0, time.UTC)
		testDate := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)

		start := time.Now()

		for range 1000 {
			converter.ParseDate(&testDate, defaultDate)
		}

		duration := time.Since(start)

		if duration > time.Millisecond*100 {
			t.Errorf("Performance test failed: took %v for 1000 calls", duration)
		}

		t.Logf("Performance test passed: 1000 calls completed in %v", duration)
	})
}

func TestParseDateEdgeCases(t *testing.T) {

	t.Run("Should handle same date as default", func(t *testing.T) {
		defaultDate := time.Date(2025, 9, 6, 0, 0, 0, 0, time.UTC)
		sameDate := time.Date(2025, 9, 6, 12, 30, 45, 0, time.UTC)

		result := converter.ParseDate(&sameDate, defaultDate)

		expectedDate := time.Date(2025, 9, 6, 0, 0, 0, 0, time.UTC)

		if !result.Equal(expectedDate) {
			t.Errorf("Expected %v, got %v", expectedDate, result)
		}

		t.Logf("Correctly handled same date with different time: %v", result.Format("2006-01-02 15:04:05"))
	})
}
