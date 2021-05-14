package bytes

import (
	"testing"

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()

			tt.fn()

			assert.Equal(t, len(tt.want), buf.Len())
			assert.Equal(t, tt.want, string(buf.Bytes()))
		})
	}
}
