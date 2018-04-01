package demux

import (
	"bytes"
	"io"
	"testing"
)

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

func TestTopic_WriteToDevice(t *testing.T) {
	userID := "TestUser"
	goodDeviceID := "GoodDeviceID"
	badDeviceID := "BadDeviceID"

	var goodBuffer, badBuffer bytes.Buffer

	topic, _ := NewTopic("TestTopic_WriteToDevice")
	topic.RegisterDevice(&Device{
		UserID:   userID,
		DeviceID: goodDeviceID,
		writer:   &goodBuffer,
	})

	topic.RegisterDevice(&Device{
		UserID:   "non-existent",
		DeviceID: "non-existent",
		writer:   &badBuffer,
	})

	type args struct {
		userID   string
		deviceID string
		message  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Working", args{userID, goodDeviceID, []byte("Hello sir")}, false},
		{"Non-existent User", args{"Derp", goodDeviceID, []byte("Goodbye")}, true},
		{"Non-existent Device", args{userID, badDeviceID, []byte("Goodbye")}, true},
		{"Empty User", args{"", goodDeviceID, []byte("Goodbye")}, true},
		{"Empty Device", args{userID, "", []byte("Goodbye")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goodBuffer.Reset()
			badBuffer.Reset()

			if err := topic.WriteToDevice(tt.args.userID, tt.args.deviceID, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Topic.WriteToDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if goodBuffer.String() != "" {
					t.Errorf("Wanted WriteToDevice to fail, but data was written: %q", goodBuffer.String())
				}

				return
			}

			if goodBuffer.String() != string(tt.args.message) {
				t.Errorf("Written message %q not what expected: %q", goodBuffer.String(), tt.args.message)
			}

			if badBuffer.String() != "" {
				t.Errorf("Another buffer was written into.")
				return
			}
		})
	}
}

func TestTopic_WriteToUser(t *testing.T) {
	userID := "TestUser"

	var goodBufferOne, goodBufferTwo, badBuffer bytes.Buffer

	topic, _ := NewTopic("TestTopic_WriteToUser")
	topic.RegisterDevice(&Device{
		UserID:   userID,
		DeviceID: "ok",
		writer:   &goodBufferOne,
	})

	topic.RegisterDevice(&Device{
		UserID:   userID,
		DeviceID: "okagain",
		writer:   &goodBufferTwo,
	})

	topic.RegisterDevice(&Device{
		UserID:   "non-existent",
		DeviceID: "non-existent",
		writer:   &badBuffer,
	})

	type args struct {
		userID  string
		message []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Working", args{userID, []byte("Hello sir")}, false},
		{"Non-existent User", args{"Derp", []byte("Goodbye")}, true},
		{"Empty User", args{"", []byte("Goodbye")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goodBufferOne.Reset()
			goodBufferTwo.Reset()
			badBuffer.Reset()

			if err := topic.WriteToUser(tt.args.userID, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Topic.WriteToUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if goodBufferOne.String() != "" {
					t.Errorf("Wanted WriteToDevice to fail, but data was written: %q", goodBufferOne.String())
				}

				if goodBufferTwo.String() != "" {
					t.Errorf("Wanted WriteToDevice to fail, but data was written: %q", goodBufferTwo.String())
				}

				return
			}

			if goodBufferOne.String() != string(tt.args.message) {
				t.Errorf("Written message in buffer one: %q not what expected: %q", goodBufferOne.String(), tt.args.message)
			}

			if goodBufferTwo.String() != string(tt.args.message) {
				t.Errorf("Written message in buffer two %q not what expected: %q", goodBufferTwo.String(), tt.args.message)
			}

			if badBuffer.String() != "" {
				t.Errorf("Another buffer was written into.")
				return
			}
		})
	}
}
