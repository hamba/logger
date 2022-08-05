/*
Package logger implements a logging package.

Example usage:

	log := logger.New(os.Stdout, logger.LogfmtFormat(), logger.Info)

	// Logger can have scoped context
	log = log.With(ctx.Str("env", "prod"))

	// All messages can have a context
	log.Error("connection error", ctx.Str("redis", conn.Name()), ctx.Int("timeout", conn.Timeout()))
*/
package logger
