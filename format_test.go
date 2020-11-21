package logger_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/hamba/logger"
	"github.com/stretchr/testify/assert"
)

func TestJsonFormat(t *testing.T) {
	f := logger.JSONFormat()

	e := &logger.Event{
		Msg:     "some message",
		Lvl:     logger.Error,
		BaseCtx: []interface{}{"x", 1, "y", 3.2},
		Ctx: []interface{}{
			"bool", true,
			"carriage_return", "bang" + string('\r') + "foo",
			"tab", "bar	baz",
			"newline", "foo\nbar",
			"escape", string('\\'),
		},
	}
	b := f.Format(e)

	expect := []byte(`{"lvl":"eror","msg":"some message","x":1,"y":3.2,"bool":true,"carriage_return":"bang\rfoo","tab":"bar\tbaz","newline":"foo\nbar","escape":"\\"}` + "\n")
	assert.Equal(t, expect, b)

	m := map[string]interface{}{}
	err := json.Unmarshal(b, &m)
	assert.NoError(t, err)
}

func TestJsonFormat_KeyError(t *testing.T) {
	f := logger.JSONFormat()

	e := &logger.Event{
		Msg: "some message",
		Lvl: logger.Error,
		Ctx: []interface{}{1, "y"},
	}
	b := f.Format(e)

	expect := []byte(`{"lvl":"eror","msg":"some message","LOGGER_ERROR":1}` + "\n")
	assert.Equal(t, expect, b)

	m := map[string]interface{}{}
	err := json.Unmarshal(b, &m)
	assert.NoError(t, err)
}

func TestJsonFormat_Ints(t *testing.T) {
	f := logger.JSONFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"int", 1, "int8", int8(2), "int16", int16(3), "int32", int32(4), "int64", int64(5)},
	}
	b := f.Format(e)

	expect := []byte(`{"lvl":"eror","msg":"","int":1,"int8":2,"int16":3,"int32":4,"int64":5}` + "\n")
	assert.Equal(t, expect, b)

	m := map[string]interface{}{}
	err := json.Unmarshal(b, &m)
	assert.NoError(t, err)
}

func TestJsonFormat_Uints(t *testing.T) {
	f := logger.JSONFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"uint", uint(1), "uint8", uint8(2), "uint16", uint16(3), "uint32", uint32(4), "uint64", uint64(5)},
	}
	b := f.Format(e)

	expect := []byte(`{"lvl":"eror","msg":"","uint":1,"uint8":2,"uint16":3,"uint32":4,"uint64":5}` + "\n")
	assert.Equal(t, expect, b)

	m := map[string]interface{}{}
	err := json.Unmarshal(b, &m)
	assert.NoError(t, err)
}

func TestJsonFormat_Floats(t *testing.T) {
	f := logger.JSONFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"float32", float32(1), "float64", float64(4.56)},
	}
	b := f.Format(e)

	expect := []byte(`{"lvl":"eror","msg":"","float32":1,"float64":4.56}` + "\n")
	assert.Equal(t, expect, b)

	m := map[string]interface{}{}
	err := json.Unmarshal(b, &m)
	assert.NoError(t, err)
}

func TestJsonFormat_Time(t *testing.T) {
	f := logger.JSONFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"time", time.Unix(1541573670, 0).UTC()},
	}
	b := f.Format(e)

	expect := []byte(`{"lvl":"eror","msg":"","time":"2018-11-07T06:54:30+0000"}` + "\n")
	assert.Equal(t, expect, b)

	m := map[string]interface{}{}
	err := json.Unmarshal(b, &m)
	assert.NoError(t, err)
}

func TestJsonFormat_Unknown(t *testing.T) {
	obj := struct {
		Name string
	}{Name: "test"}

	f := logger.JSONFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"what", obj, "nil", nil},
	}
	b := f.Format(e)

	expect := []byte(`{"lvl":"eror","msg":"","what":"{Name:test}","nil":null}` + "\n")
	assert.Equal(t, expect, b)
}

func TestLogfmtFormat(t *testing.T) {
	f := logger.LogfmtFormat()

	e := &logger.Event{
		Msg:     "some message",
		Lvl:     logger.Error,
		BaseCtx: []interface{}{"x", 1, "y", 3.2},
		Ctx: []interface{}{
			"bool", true,
			"equals", "=",
			"quote", "\"",
			"carriage_return", "bang" + string('\r') + "foo",
			"tab", "bar	baz",
			"newline", "foo\nbar",
			"escape", string('\\'),
		},
	}
	b := f.Format(e)

	expect := []byte(`lvl=eror msg="some message" x=1 y=3.200 bool=true equals="=" quote="\"" carriage_return="bang\rfoo" tab="bar\tbaz" newline="foo\nbar" escape=\\` + "\n")
	assert.Equal(t, expect, b)
}

func TestLogfmtFormat_KeyError(t *testing.T) {
	f := logger.LogfmtFormat()

	e := &logger.Event{
		Msg: "some message",
		Lvl: logger.Error,
		Ctx: []interface{}{1, "y"},
	}
	b := f.Format(e)

	expect := []byte(`lvl=eror msg="some message" LOGGER_ERROR=1` + "\n")
	assert.Equal(t, expect, b)
}

func TestLogfmtFormat_Ints(t *testing.T) {
	f := logger.LogfmtFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"int", 1, "int8", int8(2), "int16", int16(3), "int32", int32(4), "int64", int64(5)},
	}
	b := f.Format(e)

	expect := []byte(`lvl=eror msg= int=1 int8=2 int16=3 int32=4 int64=5` + "\n")
	assert.Equal(t, expect, b)
}

func TestLogfmtFormat_Uints(t *testing.T) {
	f := logger.LogfmtFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"uint", uint(1), "uint8", uint8(2), "uint16", uint16(3), "uint32", uint32(4), "uint64", uint64(5)},
	}
	b := f.Format(e)

	expect := []byte(`lvl=eror msg= uint=1 uint8=2 uint16=3 uint32=4 uint64=5` + "\n")
	assert.Equal(t, expect, b)
}

func TestLogfmtFormat_Floats(t *testing.T) {
	f := logger.LogfmtFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"float32", float32(1.23), "float64", float64(4.56)},
	}
	b := f.Format(e)

	expect := []byte(`lvl=eror msg= float32=1.230 float64=4.560` + "\n")
	assert.Equal(t, expect, b)
}

func TestLogfmtFormat_Time(t *testing.T) {
	f := logger.LogfmtFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"time", time.Unix(1541573670, 0).UTC()},
	}
	b := f.Format(e)

	expect := []byte(`lvl=eror msg= time=2018-11-07T06:54:30+0000` + "\n")
	assert.Equal(t, expect, b)
}

func TestLogfmtFormat_Unknown(t *testing.T) {
	obj := struct {
		Name string
	}{Name: "test"}

	f := logger.LogfmtFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"what", obj, "nil", nil},
	}
	b := f.Format(e)

	expect := []byte(`lvl=eror msg= what={Name:test} nil=` + "\n")
	assert.Equal(t, expect, b)
}

func TestConsoleFormat(t *testing.T) {
	f := logger.ConsoleFormat()

	e := &logger.Event{
		Msg:     "some message",
		Lvl:     logger.Error,
		BaseCtx: []interface{}{"x", 1, "y", 3.2},
		Ctx: []interface{}{
			"bool", true,
			"equals", "=",
			"quote", "\"",
			"carriage_return", "bang" + string('\r') + "foo",
			"tab", "bar	baz",
			"newline", "foo\nbar",
			"escape", string('\\'),
			"error", "test",
		},
	}
	b := f.Format(e)

	expect := []byte("\x1b[31mEROR\x1b[0m some message \x1b[36mx=\x1b[0m1 \x1b[36my=\x1b[0m3.200 \x1b[36mbool=\x1b[0mtrue \x1b[36mequals=\x1b[0m= \x1b[36mquote=\x1b[0m\\\" \x1b[36mcarriage_return=\x1b[0mbang\\rfoo \x1b[36mtab=\x1b[0mbar\\tbaz \x1b[36mnewline=\x1b[0mfoo\\nbar \x1b[36mescape=\x1b[0m\\\\ \x1b[31merror=\x1b[0m\x1b[31mtest\x1b[0m\n")
	assert.Equal(t, expect, b)
}

func TestConsoleFormat_Levels(t *testing.T) {
	tests := []struct {
		name string
		lvl  logger.Level
		want string
	}{
		{
			name: "debug",
			lvl:  logger.Debug,
			want: "\x1b[34mDBUG\x1b[0m \n",
		},
		{
			name: "info",
			lvl:  logger.Info,
			want: "\x1b[32mINFO\x1b[0m \n",
		},
		{
			name: "warning",
			lvl:  logger.Warn,
			want: "\x1b[33mWARN\x1b[0m \n",
		},
		{
			name: "error",
			lvl:  logger.Error,
			want: "\x1b[31mEROR\x1b[0m \n",
		},
		{
			name: "crit",
			lvl:  logger.Crit,
			want: "\x1b[31;1mCRIT\x1b[0m \n",
		},
		{
			name: "unknown",
			lvl:  logger.Level(1234),
			want: "\x1b[37mUNKN\x1b[0m \n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := logger.ConsoleFormat()
			e := &logger.Event{
				Lvl: test.lvl,
			}

			b := f.Format(e)

			assert.Equal(t, test.want, string(b))
		})
	}
}

func TestConsoleFormat_KeyError(t *testing.T) {
	f := logger.ConsoleFormat()

	e := &logger.Event{
		Msg: "some message",
		Lvl: logger.Error,
		Ctx: []interface{}{1, "y"},
	}
	b := f.Format(e)

	expect := []byte("\x1b[31mEROR\x1b[0m some message \x1b[31mLOGGER_ERROR=1\x1b[0m\n")
	assert.Equal(t, expect, b)
}

func TestConsoleFormat_Ints(t *testing.T) {
	f := logger.ConsoleFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"int", 1, "int8", int8(2), "int16", int16(3), "int32", int32(4), "int64", int64(5)},
	}
	b := f.Format(e)

	expect := []byte("\x1b[31mEROR\x1b[0m  \x1b[36mint=\x1b[0m1 \x1b[36mint8=\x1b[0m2 \x1b[36mint16=\x1b[0m3 \x1b[36mint32=\x1b[0m4 \x1b[36mint64=\x1b[0m5\n")
	assert.Equal(t, expect, b)
}

func TestConsoleFormat_Uints(t *testing.T) {
	f := logger.ConsoleFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"uint", uint(1), "uint8", uint8(2), "uint16", uint16(3), "uint32", uint32(4), "uint64", uint64(5)},
	}
	b := f.Format(e)

	expect := []byte("\x1b[31mEROR\x1b[0m  \x1b[36muint=\x1b[0m1 \x1b[36muint8=\x1b[0m2 \x1b[36muint16=\x1b[0m3 \x1b[36muint32=\x1b[0m4 \x1b[36muint64=\x1b[0m5\n")
	assert.Equal(t, expect, b)
}

func TestConsoleFormat_Floats(t *testing.T) {
	f := logger.ConsoleFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"float32", float32(1.23), "float64", float64(4.56)},
	}
	b := f.Format(e)

	expect := []byte("\x1b[31mEROR\x1b[0m  \x1b[36mfloat32=\x1b[0m1.230 \x1b[36mfloat64=\x1b[0m4.560\n")
	assert.Equal(t, expect, b)
}

func TestConsoleFormat_Time(t *testing.T) {
	f := logger.ConsoleFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"time", time.Unix(1541573670, 0).UTC()},
	}
	b := f.Format(e)

	expect := []byte("\x1b[31mEROR\x1b[0m  \x1b[36mtime=\x1b[0m2018-11-07T06:54:30+0000\n")
	assert.Equal(t, expect, b)
}

func TestConsoleFormat_Unknown(t *testing.T) {
	obj := struct {
		Name string
	}{Name: "test"}

	f := logger.ConsoleFormat()

	e := &logger.Event{
		Lvl: logger.Error,
		Ctx: []interface{}{"what", obj, "nil", nil},
	}
	b := f.Format(e)

	expect := []byte("\x1b[31mEROR\x1b[0m  \x1b[36mwhat=\x1b[0m{Name:test} \x1b[36mnil=\x1b[0m\n")
	assert.Equal(t, expect, b)
}
