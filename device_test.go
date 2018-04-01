package demux

import (
	"bufio"
	"bytes"
	"io"
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
	user, _ := topic.addUser("TestUser")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			writer := bufio.NewWriter(&b)
			if err := user.addDevice(tt.args.name, writer); (err != nil) != tt.wantErr {
				t.Errorf("User.AddDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

		})

		// Test nil writer error
		if err := user.addDevice("nil writer", nil); err == nil {
			t.Errorf("AddDevice should have failed with nil writer argument.")
		}
	}
}

func TestDevice_validate(t *testing.T) {
	var testIOWriter bytes.Buffer
	type fields struct {
		UserID   string
		DeviceID string
		writer   io.Writer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"Valid", fields{UserID: "valid", DeviceID: "valid", writer: &testIOWriter}, false},
		{"Missing UserID", fields{UserID: "", DeviceID: "valid", writer: &testIOWriter}, true},
		{"Missing DeviceID", fields{UserID: "valid", DeviceID: "", writer: &testIOWriter}, true},
		{"Nil Writer", fields{UserID: "valid", DeviceID: "valid", writer: nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Device{
				UserID:   tt.fields.UserID,
				DeviceID: tt.fields.DeviceID,
				writer:   tt.fields.writer,
			}
			if err := d.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Device.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
