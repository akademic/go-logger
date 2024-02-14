package logger

type LogLevel string

const (
	LogOff   LogLevel = "off"
	LogError LogLevel = "error"
	LogInfo  LogLevel = "info"
	LogDebug LogLevel = "debug"
)

func (ll LogLevel) String() string {
	return string(ll)
}

func (ll LogLevel) CanLog(level LogLevel) bool {
	if ll == LogError {
		return level == LogError
	}

	if ll == LogInfo {
		return level == LogError || level == LogInfo
	}

	if ll == LogDebug {
		return true
	}

	return false
}
