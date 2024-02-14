package logger

// test for LogLevel CanLog method

import (
	"testing"
)

func TestLogLevel_CanLog(t *testing.T) {
	testTable := []struct {
		name   string
		ll     LogLevel
		level  LogLevel
		result bool
	}{
		{
			name:   "LogOff",
			ll:     LogOff,
			level:  LogInfo,
			result: false,
		},
		{
			name:   "LogErrorInfo",
			ll:     LogError,
			level:  LogInfo,
			result: false,
		},
		{
			name:   "LogError",
			ll:     LogError,
			level:  LogError,
			result: true,
		},
		{
			name:   "LogErrorDebug",
			ll:     LogError,
			level:  LogDebug,
			result: false,
		},
		{
			name:   "LogInfo",
			ll:     LogInfo,
			level:  LogInfo,
			result: true,
		},
		{
			name:   "LogInfoError",
			ll:     LogInfo,
			level:  LogError,
			result: true,
		},
		{
			name:   "LogDebugInfo",
			ll:     LogDebug,
			level:  LogInfo,
			result: true,
		},
		{
			name:   "LogDebugError",
			ll:     LogDebug,
			level:  LogError,
			result: true,
		},
		{
			name:   "LogDebug",
			ll:     LogDebug,
			level:  LogDebug,
			result: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			if test.ll.CanLog(test.level) != test.result {
				t.Errorf("expected %v, got %v", test.result, test.ll.CanLog(test.level))
			}
		})
	}
}
