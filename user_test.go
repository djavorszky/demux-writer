package demux

import (
	"bytes"
	"testing"
)

func TestTopic_AddUser(t *testing.T) {
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
	topic, _ := NewTopic("TestTopic_AddUser")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				user *User
				err  error
			)
			if user, err = topic.AddUser(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Topic.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if user == nil {
				t.Errorf("User is nil")
				return
			}

			if user.Name != tt.args.name {
				t.Errorf("User added with wrong name. Expected: %q, got: %q", tt.args.name, user.Name)
				return
			}

			addedUser, ok := topic.users[tt.args.name]
			if !ok {
				t.Errorf("User should have been added, wasn't")
				return
			}

			if addedUser.Name != user.Name {
				t.Errorf("Expected user to be added under name %q, got %q", user.Name, addedUser.Name)
				return
			}

			if user.devices == nil {
				t.Errorf("Devices added with nil slice")
			}
		})
	}
}

func TestTopic_RegisterDevice(t *testing.T) {
	var testIOWriter bytes.Buffer
	tests := []struct {
		name    string
		d       *Device
		wantErr bool
	}{
		{"Valid", &Device{UserID: "valid", DeviceID: "valid", writer: &testIOWriter}, false},
		{"Missing UserID", &Device{UserID: "", DeviceID: "valid", writer: &testIOWriter}, true},
		{"Missing DeviceID", &Device{UserID: "valid", DeviceID: "", writer: &testIOWriter}, true},
		{"Nil Writer", &Device{UserID: "valid", DeviceID: "valid", writer: nil}, true},
		// Duplicate
		{"Duplicate DeviceID for UserID", &Device{UserID: "valid", DeviceID: "valid", writer: &testIOWriter}, true},
	}
	topic, _ := NewTopic("TestTopic_RegisterDevice")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := topic.RegisterDevice(tt.d); (err != nil) != tt.wantErr {
				t.Errorf("Topic.RegisterDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}
		})
	}
}

func TestTopic_getUser(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
	}{
		{"User1", args{"User1"}},
		{"User1 Again", args{"User1"}},
		{"User2", args{"User2"}},
		{"User2 again", args{"User2"}},
	}
	topic, _ := NewTopic("TestTopic_getUser")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := topic.getUser(tt.args.name)

			if user == nil {
				t.Errorf("Returned user is nil")
			}

			if user.Name != tt.args.name {
				t.Errorf("Name mismatch. Expected %q, got %q", tt.args.name, user.Name)
			}
		})
	}
}
