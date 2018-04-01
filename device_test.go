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
