package config

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	if err := loadEnvForTesting(); err != nil {
		os.Setenv("APP_NAME", "Senkou Catalyst")
		os.Setenv("APP_ENV", "testing")
		os.Setenv("TEST_INT_VAR", "42")
		os.Setenv("TEST_BOOL_VAR", "true")
		os.Setenv("TEST_DURATION_VAR", "5s")
	}

	code := m.Run()

	os.Exit(code)

}

func TestEnvironmentVariable(t *testing.T) {

	t.Run("GetEnv should return the value of an existing env variable", func(t *testing.T) {
		key := "APP_NAME"
		expectedValue := "Senkou Catalyst"

		value := GetEnv(key, "DefaultApp")
		if value != expectedValue {
			t.Errorf("Expected %s but got %s", expectedValue, value)
		} else {
			t.Logf("GetEnv returned expected value: %s", value)
		}
	})

	t.Run("GetEnv should return fallback for non-existing env variable", func(t *testing.T) {
		key := "NON_EXISTING_ENV_VAR"
		fallback := "FallbackValue"

		value := GetEnv(key, fallback)
		if value != fallback {
			t.Errorf("Expected fallback %s but got %s", fallback, value)
		} else {
			t.Logf("GetEnv returned fallback value as expected: %s", value)
		}
	})

	t.Run("GetEnvAsInt should return integer value", func(t *testing.T) {
		key := "TEST_INT_VAR"
		expectedValue := 42

		value := GetEnvAsInt(key, 0)
		if value != expectedValue {
			t.Errorf("Expected %d but got %d", expectedValue, value)
		} else {
			t.Logf("GetEnvAsInt returned expected value: %d", value)
		}
	})

	t.Run("GetEnvAsInt should return fallback for non-integer value", func(t *testing.T) {
		key := "NON_INTEGER_VAR"
		fallback := 99

		value := GetEnvAsInt(key, fallback)
		if value != fallback {
			t.Errorf("Expected fallback %d but got %d", fallback, value)
		} else {
			t.Logf("GetEnvAsInt returned fallback value as expected: %d", value)
		}
	})

	t.Run("GetEnvAsBool should return boolean value", func(t *testing.T) {
		key := "TEST_BOOL_VAR"
		expectedValue := true

		value := GetEnvAsBool(key, false)
		if value != expectedValue {
			t.Errorf("Expected %t but got %t", expectedValue, value)
		} else {
			t.Logf("GetEnvAsBool returned expected value: %t", value)
		}
	})

	t.Run("GetEnvAsDuration should return duration value", func(t *testing.T) {
		key := "TEST_DURATION_VAR"
		expectedValue := "5s"

		value := GetEnvAsDuration(key, 0)
		if value.String() != expectedValue {
			t.Errorf("Expected %s but got %s", expectedValue, value.String())
		} else {
			t.Logf("GetEnvAsDuration returned expected value: %s", value.String())
		}
	})

}
