package bytes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	buf := NewBuffer(512)

	tests := []struct {
		name string
		fn   func()
		want string
	}{
		{
			name: "WriteByte",
			fn:   func() { buf.WriteByte('v') },
			want: "v",
		},
		{
			name: "WriteString",
			fn:   func() { buf.WriteString("foo") },
			want: "foo",
		},
		{
			name: "WriteRune",
			fn:   func() { buf.WriteRune('f') },
			want: "f",
		},
		{
			name: "Write",
			fn:   func() { buf.Write([]byte("foo")) },
			want: "foo",
		},
		{
			name: "AppendIntPositive",
			fn:   func() { buf.AppendInt(42) },
			want: "42",
		},
		{
			name: "AppendIntNegative",
			fn:   func() { buf.AppendInt(-42) },
			want: "-42",
		},
		{
			name: "AppendUint",
			fn:   func() { buf.AppendUint(42) },
			want: "42",
		},
		{
			name: "AppendBool",
			fn:   func() { buf.AppendBool(true) },
			want: "true",
		},
		{
			name: "AppendFloat64",
			fn:   func() { buf.AppendFloat(3.14, 'f', 3, 64) },
			want: "3.140",
		},
		{
			name: "AppendTime",
			fn:   func() { buf.AppendTime(time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC), time.DateTime) },
			want: "2026-01-02 15:04:05",
		},
		{
			name: "AppendDuration",
			fn:   func() { buf.AppendDuration(3*time.Hour + 2*time.Minute + time.Second) },
			want: "3h2m1s",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf.Reset()

			test.fn()

			assert.Equal(t, len(test.want), buf.Len())
			assert.Equal(t, test.want, string(buf.Bytes()))
		})
	}
}
