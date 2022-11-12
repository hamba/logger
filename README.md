![Logo](http://svg.wiersma.co.za/hamba/project?title=logger&tag=fast%20Go%20logger)

[![Go Report Card](https://goreportcard.com/badge/github.com/hamba/logger)](https://goreportcard.com/report/github.com/hamba/logger)
[![Build Status](https://github.com/hamba/logger/actions/workflows/test.yml/badge.svg)](https://github.com/hamba/logger/actions)
[![Coverage Status](https://coveralls.io/repos/github/hamba/logger/badge.svg?branch=master)](https://coveralls.io/github/hamba/logger?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/hamba/logger/v2.svg)](https://pkg.go.dev/github.com/hamba/logger/v2)
[![GitHub release](https://img.shields.io/github/release/hamba/logger.svg)](https://github.com/hamba/logger/releases)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hamba/logger/master/LICENSE)

logger is a fast Go logging package made to be simple but effective.

## Overview

Install with:

```shell
go get github.com/hamba/logger/v2
```

#### Formatters

* **JSON**
* **Logfmt**
* **Console**

#### Writers

* **SyncWriter** Write synchronised to a Writer

**Note:** This project has renamed the default branch from `master` to `main`. You will need to update your local environment.

## Examples

```go
log := logger.New(os.Stdout, logger.LogfmtFormat(), logger.Info)

// Logger can have scoped context
log = log.With(ctx.Str("env", "prod"))

// All messages can have a context
log.Warn("connection error", ctx.Str("redis", "dsn_1"), ctx.Int("timeout", conn.Timeout()))
```

Will log the message

```
lvl=warn msg="connection error" env=prod redis=dsn_1 timeout=0.500
```

More examples can be found in the [godocs](https://godoc.org/github.com/hamba/logger).
