package demux

import (
	"bufio"
	"bytes"
	"testing"
)

func TestUser_AddDevice(t *testing.T) {
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
	topic, _ := NewTopic("TestUser_AddDevice")
	user, _ := topic.AddUser("TestUser")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			writer := bufio.NewWriter(&b)
			if err := user.AddDevice(tt.args.name, writer); (err != nil) != tt.wantErr {
				t.Errorf("User.AddDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

		})

		// Test nil writer error
		if err := user.AddDevice("nil writer", nil); err == nil {
			t.Errorf("AddDevice should have failed with nil writer argument.")
		}
	}
}
