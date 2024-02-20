package logging

import "log/slog"

func Init() *slog.Logger {
	return slog.Default()
}
