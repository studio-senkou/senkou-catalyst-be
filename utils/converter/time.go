package converter

import (
	"fmt"
	"time"
)

func ParseStringToTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, fmt.Errorf("input time string is empty")
	}

	timeFormats := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		time.RFC3339Nano,
		"2006-01-02 15:04:05.000",
		"2006-01-02",
	}

	var lastErr error
	for _, format := range timeFormats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		} else {
			lastErr = err
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time string '%s' with any known format. Last error: %v", timeStr, lastErr)
}

func ParseStringToTimeWithLocation(timeStr string, location *time.Location) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, fmt.Errorf("input time string is empty")
	}

	if location == nil {
		location = time.UTC
	}

	timeFormats := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		time.RFC3339Nano,
		"2006-01-02 15:04:05.000",
		"2006-01-02",
	}

	var lastErr error
	for _, format := range timeFormats {
		if t, err := time.ParseInLocation(format, timeStr, location); err == nil {
			return t, nil
		} else {
			lastErr = err
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time string '%s' with any known format in location %s. Last error: %v", timeStr, location.String(), lastErr)
}

func ParseMidtransTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, fmt.Errorf("input time string is empty")
	}

	jakartaLocation, _ := time.LoadLocation("Asia/Jakarta")

	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, jakartaLocation)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("unable to parse Midtrans time '%s': %v", timeStr, err)
		}
	}

	return t, nil
}
