package converter

import (
	"testing"
	"time"
)

func TestParseStringToTime(t *testing.T) {
	t.Run("Should parse RFC3339 format", func(t *testing.T) {
		timeStr := "2023-12-25T15:30:45Z"
		result, err := ParseStringToTime(timeStr)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Should parse standard format", func(t *testing.T) {
		timeStr := "2023-12-25 15:30:45"
		result, err := ParseStringToTime(timeStr)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Should parse date only format", func(t *testing.T) {
		timeStr := "2023-12-25"
		result, err := ParseStringToTime(timeStr)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Should parse RFC3339Nano format", func(t *testing.T) {
		timeStr := "2023-12-25T15:30:45.123456789Z"
		result, err := ParseStringToTime(timeStr)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := time.Date(2023, 12, 25, 15, 30, 45, 123456789, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Should parse format with milliseconds", func(t *testing.T) {
		timeStr := "2023-12-25 15:30:45.123"
		result, err := ParseStringToTime(timeStr)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := time.Date(2023, 12, 25, 15, 30, 45, 123000000, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Should return error for empty string", func(t *testing.T) {
		result, err := ParseStringToTime("")

		if err == nil {
			t.Error("Expected error for empty string")
		}

		if !result.IsZero() {
			t.Errorf("Expected zero time, got %v", result)
		}
	})

	t.Run("Should return error for invalid format", func(t *testing.T) {
		timeStr := "invalid-time-format"
		result, err := ParseStringToTime(timeStr)

		if err == nil {
			t.Error("Expected error for invalid format")
		}

		if !result.IsZero() {
			t.Errorf("Expected zero time, got %v", result)
		}
	})

	t.Run("Should return error for partial invalid format", func(t *testing.T) {
		timeStr := "2023-13-45"
		result, err := ParseStringToTime(timeStr)

		if err == nil {
			t.Error("Expected error for invalid date")
		}

		if !result.IsZero() {
			t.Errorf("Expected zero time, got %v", result)
		}
	})
}

func TestParseStringToTimeWithLocation(t *testing.T) {
	t.Run("Should parse time with Jakarta location", func(t *testing.T) {
		jakarta, _ := time.LoadLocation("Asia/Jakarta")
		timeStr := "2023-12-25 15:30:45"
		result, err := ParseStringToTimeWithLocation(timeStr, jakarta)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := time.Date(2023, 12, 25, 15, 30, 45, 0, jakarta)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Should parse time with UTC when location is nil", func(t *testing.T) {
		timeStr := "2023-12-25 15:30:45"
		result, err := ParseStringToTimeWithLocation(timeStr, nil)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Should parse RFC3339 with specified location", func(t *testing.T) {
		tokyo, _ := time.LoadLocation("Asia/Tokyo")
		timeStr := "2023-12-25T15:30:45Z"
		result, err := ParseStringToTimeWithLocation(timeStr, tokyo)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Should parse date only with location", func(t *testing.T) {
		newYork, _ := time.LoadLocation("America/New_York")
		timeStr := "2023-12-25"
		result, err := ParseStringToTimeWithLocation(timeStr, newYork)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := time.Date(2023, 12, 25, 0, 0, 0, 0, newYork)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Should return error for empty string", func(t *testing.T) {
		jakarta, _ := time.LoadLocation("Asia/Jakarta")
		result, err := ParseStringToTimeWithLocation("", jakarta)

		if err == nil {
			t.Error("Expected error for empty string")
		}

		if !result.IsZero() {
			t.Errorf("Expected zero time, got %v", result)
		}
	})

	t.Run("Should return error for invalid format with location", func(t *testing.T) {
		jakarta, _ := time.LoadLocation("Asia/Jakarta")
		timeStr := "invalid-format"
		result, err := ParseStringToTimeWithLocation(timeStr, jakarta)

		if err == nil {
			t.Error("Expected error for invalid format")
		}

		if !result.IsZero() {
			t.Errorf("Expected zero time, got %v", result)
		}
	})
}

func TestParseMidtransTime(t *testing.T) {
	t.Run("Should parse Midtrans time format with Jakarta timezone", func(t *testing.T) {
		timeStr := "2023-12-25 15:30:45"
		result, err := ParseMidtransTime(timeStr)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		jakarta, _ := time.LoadLocation("Asia/Jakarta")
		expected := time.Date(2023, 12, 25, 15, 30, 45, 0, jakarta)
		if !result.Equal(expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Should parse standard format as fallback", func(t *testing.T) {
		timeStr := "2023-12-25 15:30:45"
		result, err := ParseMidtransTime(timeStr)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result.Year() != 2023 || result.Month() != 12 || result.Day() != 25 {
			t.Errorf("Expected correct date parsing, got %v", result)
		}
	})

	t.Run("Should return error for empty string", func(t *testing.T) {
		result, err := ParseMidtransTime("")

		if err == nil {
			t.Error("Expected error for empty string")
		}

		if !result.IsZero() {
			t.Errorf("Expected zero time, got %v", result)
		}
	})

	t.Run("Should return error for invalid format", func(t *testing.T) {
		timeStr := "invalid-midtrans-format"
		result, err := ParseMidtransTime(timeStr)

		if err == nil {
			t.Error("Expected error for invalid format")
		}

		if !result.IsZero() {
			t.Errorf("Expected zero time, got %v", result)
		}
	})

	t.Run("Should return error for partial invalid date", func(t *testing.T) {
		timeStr := "2023-13-45 25:70:99"
		result, err := ParseMidtransTime(timeStr)

		if err == nil {
			t.Error("Expected error for invalid date values")
		}

		if !result.IsZero() {
			t.Errorf("Expected zero time, got %v", result)
		}
	})

	t.Run("Should handle midnight time", func(t *testing.T) {
		timeStr := "2023-12-25 00:00:00"
		result, err := ParseMidtransTime(timeStr)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
			t.Errorf("Expected midnight time, got %02d:%02d:%02d", result.Hour(), result.Minute(), result.Second())
		}
	})

	t.Run("Should handle end of day time", func(t *testing.T) {
		timeStr := "2023-12-25 23:59:59"
		result, err := ParseMidtransTime(timeStr)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result.Hour() != 23 || result.Minute() != 59 || result.Second() != 59 {
			t.Errorf("Expected 23:59:59, got %02d:%02d:%02d", result.Hour(), result.Minute(), result.Second())
		}
	})
}
