/*
Package logger implements a logging package.

logger implements github.com/hamba/pkg Logger interface.

Example usage:
	// Composable handlers
	h := logger.LevelFilterHandler(
		logger.Info,
		logger.StreamHandler(os.Stdout, logger.LogfmtFormat()),
	)

	// The logger can have an initial context
	l := logger.New(h, "env", "prod")

	// All messages can have a context
	l.Error("connection error", "redis", conn.Name(), "timeout", conn.Timeout())
*/
package logger
