package example

import (
	"reflect"
	"testing"
)

func TestNewResources(t *testing.T) {
	type args struct {
		locale string
	}
	tests := []struct {
		name string
		args args
		want Resources
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResources(tt.args.locale); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResources() = %v, want %v", got, tt.want)
			}
		})
	}
}