![Logo](http://svg.wiersma.co.za/github/project?lang=go&title=logger&tag=fast%20Go%20logger)

[![Go Report Card](https://goreportcard.com/badge/github.com/hamba/logger)](https://goreportcard.com/report/github.com/hamba/logger)
[![Build Status](https://travis-ci.org/hamba/logger.svg?branch=master)](https://travis-ci.org/hamba/logger)
[![Coverage Status](https://coveralls.io/repos/github/hamba/logger/badge.svg?branch=master)](https://coveralls.io/github/hamba/logger?branch=master)
[![GoDoc](https://godoc.org/github.com/hamba/logger?status.svg)](https://godoc.org/github.com/hamba/logger)
[![GitHub release](https://img.shields.io/github/release/hamba/logger.svg)](https://github.com/hamba/logger/releases)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hamba/logger/master/LICENSE)

logger is a fast Go logging package made to be simple but effective.

## Overview

Install with:

```shell
go get github.com/hamba/logger
```

#### Formatters

* **JSON**
* **Logfmt**

#### Handlers

* **StreamHandler** Write directly to a Writer, usually `os.Stdout`
* **BufferedStreamHandler** A buffered version of `StreamHandler`
* **FilterHandler** Filter log line using a function
* **LevelFilterHandler** Filter log line by level
* **DiscardHandler** Discard everything

## Examples

```go
// Composable handlers
h := logger.LevelFilterHandler(
    logger.Info,
    logger.StreamHandler(os.Stdout, logger.LogfmtFormat()),
)

// The logger can have an initial context
l := logger.New(h, "env", "prod")

// All messages can have a context
l.Warn("connection error", "redis", conn.Name(), "timeout", conn.Timeout())
```

Will log the message

```
lvl=warn msg="connection error" redis=dsn_1 timeout=0.500
```

More examples can be found in the [godocs](https://godoc.org/github.com/hamba/logger).