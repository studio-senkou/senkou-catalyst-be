package converter

import (
	"testing"
)

func TestStrToUint(t *testing.T) {
	t.Run("Should convert valid positive number string", func(t *testing.T) {
		result, err := StrToUint("123")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != 123 {
			t.Errorf("Expected 123, got %d", result)
		}
	})

	t.Run("Should convert zero string", func(t *testing.T) {
		result, err := StrToUint("0")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != 0 {
			t.Errorf("Expected 0, got %d", result)
		}
	})

	t.Run("Should convert large valid number", func(t *testing.T) {
		result, err := StrToUint("4294967295")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != 4294967295 {
			t.Errorf("Expected 4294967295, got %d", result)
		}
	})

	t.Run("Should return error for negative number", func(t *testing.T) {
		result, err := StrToUint("-123")

		if err == nil {
			t.Error("Expected error for negative number")
		}

		if result != 0 {
			t.Errorf("Expected 0 for error case, got %d", result)
		}
	})

	t.Run("Should return error for non-numeric string", func(t *testing.T) {
		result, err := StrToUint("abc")

		if err == nil {
			t.Error("Expected error for non-numeric string")
		}

		if result != 0 {
			t.Errorf("Expected 0 for error case, got %d", result)
		}
	})

	t.Run("Should return error for empty string", func(t *testing.T) {
		result, err := StrToUint("")

		if err == nil {
			t.Error("Expected error for empty string")
		}

		if result != 0 {
			t.Errorf("Expected 0 for error case, got %d", result)
		}
	})

	t.Run("Should return error for string with decimal", func(t *testing.T) {
		result, err := StrToUint("123.45")

		if err == nil {
			t.Error("Expected error for decimal string")
		}

		if result != 0 {
			t.Errorf("Expected 0 for error case, got %d", result)
		}
	})

	t.Run("Should return error for string with spaces", func(t *testing.T) {
		result, err := StrToUint(" 123 ")

		if err == nil {
			t.Error("Expected error for string with spaces")
		}

		if result != 0 {
			t.Errorf("Expected 0 for error case, got %d", result)
		}
	})

	t.Run("Should return error for mixed alphanumeric", func(t *testing.T) {
		result, err := StrToUint("123abc")

		if err == nil {
			t.Error("Expected error for mixed alphanumeric string")
		}

		if result != 0 {
			t.Errorf("Expected 0 for error case, got %d", result)
		}
	})

	t.Run("Should handle edge case with leading zeros", func(t *testing.T) {
		result, err := StrToUint("00123")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != 123 {
			t.Errorf("Expected 123, got %d", result)
		}
	})
}
