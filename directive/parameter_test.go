package directive

import (
	"reflect"
	"testing"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
)

func Test_unescapeParameter(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			`123`,
			`123`,
		},
		{
			`abc`,
			`abc`,
		},
		{
			`aa \" bb \" cc`,
			`aa \" bb \" cc`,
		},
		{
			`""`,
			``,
		},
		{
			`"abc"`,
			`abc`,
		},
		{
			`"aa \" bb \" cc"`,
			`aa " bb " cc`,
		},
		{
			`"aa \\ bb \\ cc"`,
			`aa \ bb \ cc`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := unescapeParameter(bytes.Bytes(tt.name))
			if !reflect.DeepEqual(got, bytes.Bytes(tt.want)) {
				t.Errorf("unescapeParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsArrayOfTypes(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			`[@a]`,
			true,
		},
		{
			`[@世界]`,
			false,
		},
		{
			`[@@]`,
			false,
		},
		{
			`[@]`,
			false,
		},
		{
			`[]`,
			false,
		},
		{
			`[`,
			false,
		},
		{
			`]`,
			false,
		},
		{
			`[@aaa`,
			false,
		},
		{
			`@aaa]`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsArrayOfTypes(bytes.Bytes(tt.name)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsArrayOfTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}
