package logger

type Config struct {
	Level          LogLevel
	ComponentLevel map[string]LogLevel
}
