package demux

import (
	"testing"
)

func TestNewTopic(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Empty name", args{""}, true},
		{"Valid name", args{"name"}, false},
		// Duplicate name has to come after Valid name!
		{"Duplicate name", args{"name"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTopic(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTopic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got.users == nil {
				t.Errorf("NewTopic() returned with nil Users map")
				return
			}
		})
	}
}
