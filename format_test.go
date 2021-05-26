package logger_test

import (
	"testing"
	"time"

	"github.com/hamba/logger/v2"
	"github.com/hamba/logger/v2/internal/bytes"
	"github.com/stretchr/testify/assert"
)

func TestJsonFormat(t *testing.T) {
	fmtr := logger.JSONFormat()

	buf := bytes.NewBuffer(512)
	fmtr.WriteMessage(buf, 0, logger.Error, "some message")
	fmtr.AppendKey(buf, "error")
	fmtr.AppendString(buf, "some error")
	fmtr.AppendEndMarker(buf)

	want := `{"lvl":"eror","msg":"some message","error":"some error"}` + "\n"
	assert.Equal(t, want, string(buf.Bytes()))
}

func TestJsonFormat_Array(t *testing.T) {
	fmtr := logger.JSONFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendArrayStart(buf)
	fmtr.AppendArraySep(buf)
	fmtr.AppendArrayEnd(buf)

	assert.Equal(t, "[,]", string(buf.Bytes()))
}

func TestJsonFormat_Strings(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "equals",
			in:   "=",
			want: `"="`,
		},
		{
			name: "quote",
			in:   "\"",
			want: `"\""`,
		},
		{
			name: "carriage_return",
			in:   "bang\rfoo",
			want: `"bang\rfoo"`,
		},
		{
			name: "tab",
			in: "bar	baz",
			want: `"bar\tbaz"`,
		},
		{
			name: "newline",
			in:   "foo\nbar",
			want: `"foo\nbar"`,
		},
		{
			name: "escape",
			in:   string('\\'),
			want: `"\\"`,
		},
		{
			name: "special chars",
			in:   "some string with \"special ❤️ chars\" and somewhat realistic length",
			want: `"some string with \"special ❤️ chars\" and somewhat realistic length"`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			fmtr := logger.JSONFormat()

			buf := bytes.NewBuffer(512)
			fmtr.AppendString(buf, test.in)

			assert.Equal(t, test.want, string(buf.Bytes()))
		})
	}
}

func TestJsonFormat_Bool(t *testing.T) {
	fmtr := logger.JSONFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendBool(buf, true)
	fmtr.AppendBool(buf, false)

	assert.Equal(t, "truefalse", string(buf.Bytes()))
}

func TestJsonFormat_Int(t *testing.T) {
	fmtr := logger.JSONFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendInt(buf, 5)

	assert.Equal(t, "5", string(buf.Bytes()))
}

func TestJsonFormat_Uint(t *testing.T) {
	fmtr := logger.JSONFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendUint(buf, 5)

	assert.Equal(t, "5", string(buf.Bytes()))
}

func TestJsonFormat_Float(t *testing.T) {
	fmtr := logger.JSONFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendFloat(buf, 4.56)

	assert.Equal(t, "4.56", string(buf.Bytes()))
}

func TestJsonFormat_Time(t *testing.T) {
	fmtr := logger.JSONFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendTime(buf, time.Unix(1541573670, 0).UTC())

	assert.Equal(t, `"2018-11-07T06:54:30+0000"`, string(buf.Bytes()))
}

func TestJsonFormat_Duration(t *testing.T) {
	fmtr := logger.JSONFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendDuration(buf, time.Second)

	assert.Equal(t, `"1s"`, string(buf.Bytes()))
}

func TestJsonFormat_Interface(t *testing.T) {
	obj := struct {
		Name string
	}{Name: "test"}

	fmtr := logger.JSONFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendInterface(buf, obj)
	fmtr.AppendInterface(buf, nil)

	assert.Equal(t, `"{Name:test}"null`, string(buf.Bytes()))
}

func TestLogfmtFormat(t *testing.T) {
	fmtr := logger.LogfmtFormat()

	buf := bytes.NewBuffer(512)
	fmtr.WriteMessage(buf, 0, logger.Error, "some message")
	fmtr.AppendKey(buf, "error")
	fmtr.AppendString(buf, "some error")
	fmtr.AppendEndMarker(buf)

	want := `lvl=eror msg="some message" error="some error"` + "\n"
	assert.Equal(t, want, string(buf.Bytes()))
}

func TestLogfmtFormat_Array(t *testing.T) {
	fmtr := logger.LogfmtFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendArrayStart(buf)
	fmtr.AppendArraySep(buf)
	fmtr.AppendArrayEnd(buf)

	assert.Equal(t, ",", string(buf.Bytes()))
}

func TestLogfmtFormat_Strings(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "equals",
			in:   "=",
			want: `"="`,
		},
		{
			name: "quote",
			in:   "\"",
			want: `"\""`,
		},
		{
			name: "carriage_return",
			in:   "bang\rfoo",
			want: `"bang\rfoo"`,
		},
		{
			name: "tab",
			in: "bar	baz",
			want: `"bar\tbaz"`,
		},
		{
			name: "newline",
			in:   "foo\nbar",
			want: `"foo\nbar"`,
		},
		{
			name: "escape",
			in:   string('\\'),
			want: `\\`,
		},
		{
			name: "special chars",
			in:   "some string with \"special ❤️ chars\" and somewhat realistic length",
			want: `"some string with \"special ❤️ chars\" and somewhat realistic length"`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			fmtr := logger.LogfmtFormat()

			buf := bytes.NewBuffer(512)
			fmtr.AppendString(buf, test.in)

			assert.Equal(t, test.want, string(buf.Bytes()))
		})
	}
}

func TestLogfmtFormat_Bool(t *testing.T) {
	fmtr := logger.LogfmtFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendBool(buf, true)
	fmtr.AppendBool(buf, false)

	assert.Equal(t, "truefalse", string(buf.Bytes()))
}

func TestLogfmtFormat_Int(t *testing.T) {
	fmtr := logger.LogfmtFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendInt(buf, 5)

	assert.Equal(t, "5", string(buf.Bytes()))
}

func TestLogfmtFormat_Uint(t *testing.T) {
	fmtr := logger.LogfmtFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendUint(buf, 5)

	assert.Equal(t, "5", string(buf.Bytes()))
}

func TestLogfmtFormat_Float(t *testing.T) {
	fmtr := logger.LogfmtFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendFloat(buf, 4.56)

	assert.Equal(t, "4.560", string(buf.Bytes()))
}

func TestLogfmtFormat_Time(t *testing.T) {
	fmtr := logger.LogfmtFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendTime(buf, time.Unix(1541573670, 0).UTC())

	assert.Equal(t, `2018-11-07T06:54:30+0000`, string(buf.Bytes()))
}

func TestLogfmtFormat_Duration(t *testing.T) {
	fmtr := logger.LogfmtFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendDuration(buf, time.Second)

	assert.Equal(t, `1s`, string(buf.Bytes()))
}

func TestLogfmtFormat_Interface(t *testing.T) {
	obj := struct {
		Name string
	}{Name: "test"}

	fmtr := logger.LogfmtFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendInterface(buf, obj)
	fmtr.AppendInterface(buf, nil)

	assert.Equal(t, `{Name:test}`, string(buf.Bytes()))
}

func TestConsoleFormat(t *testing.T) {
	fmtr := logger.ConsoleFormat()

	buf := bytes.NewBuffer(512)
	fmtr.WriteMessage(buf, 0, logger.Error, "some message")
	fmtr.AppendKey(buf, "error")
	fmtr.AppendString(buf, "some error")
	fmtr.AppendEndMarker(buf)

	want := "\x1b[31mEROR\x1b[0m some message \x1b[31merror=\x1b[0msome error\n"
	assert.Equal(t, want, string(buf.Bytes()))
}

func TestConsoleFormat_Array(t *testing.T) {
	fmtr := logger.ConsoleFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendArrayStart(buf)
	fmtr.AppendArraySep(buf)
	fmtr.AppendArrayEnd(buf)

	assert.Equal(t, ",", string(buf.Bytes()))
}

func TestConsoleFormat_Strings(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "equals",
			in:   "=",
			want: `=`,
		},
		{
			name: "quote",
			in:   "\"",
			want: `\"`,
		},
		{
			name: "carriage_return",
			in:   "bang\rfoo",
			want: `bang\rfoo`,
		},
		{
			name: "tab",
			in: "bar	baz",
			want: `bar\tbaz`,
		},
		{
			name: "newline",
			in:   "foo\nbar",
			want: `foo\nbar`,
		},
		{
			name: "escape",
			in:   string('\\'),
			want: `\\`,
		},
		{
			name: "special chars",
			in:   "some string with \"special ❤️ chars\" and somewhat realistic length",
			want: `some string with \"special ❤️ chars\" and somewhat realistic length`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			fmtr := logger.ConsoleFormat()

			buf := bytes.NewBuffer(512)
			fmtr.AppendString(buf, test.in)

			assert.Equal(t, test.want, string(buf.Bytes()))
		})
	}
}

func TestConsoleFormat_Bool(t *testing.T) {
	fmtr := logger.ConsoleFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendBool(buf, true)
	fmtr.AppendBool(buf, false)

	assert.Equal(t, "truefalse", string(buf.Bytes()))
}

func TestConsoleFormat_Int(t *testing.T) {
	fmtr := logger.ConsoleFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendInt(buf, 5)

	assert.Equal(t, "5", string(buf.Bytes()))
}

func TestConsoleFormat_Uint(t *testing.T) {
	fmtr := logger.ConsoleFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendUint(buf, 5)

	assert.Equal(t, "5", string(buf.Bytes()))
}

func TestConsoleFormat_Float(t *testing.T) {
	fmtr := logger.ConsoleFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendFloat(buf, 4.56)

	assert.Equal(t, "4.560", string(buf.Bytes()))
}

func TestConsoleFormat_Time(t *testing.T) {
	fmtr := logger.ConsoleFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendTime(buf, time.Unix(1541573670, 0).UTC())

	assert.Equal(t, `2018-11-07T06:54:30+0000`, string(buf.Bytes()))
}

func TestConsoleFormat_Duration(t *testing.T) {
	fmtr := logger.ConsoleFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendDuration(buf, time.Second)

	assert.Equal(t, `1s`, string(buf.Bytes()))
}

func TestConsoleFormat_Interface(t *testing.T) {
	obj := struct {
		Name string
	}{Name: "test"}

	fmtr := logger.ConsoleFormat()

	buf := bytes.NewBuffer(512)
	fmtr.AppendInterface(buf, obj)
	fmtr.AppendInterface(buf, nil)

	assert.Equal(t, `{Name:test}`, string(buf.Bytes()))
}
